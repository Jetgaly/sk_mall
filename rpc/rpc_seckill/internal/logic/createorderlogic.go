package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_seckill/internal/svc"
	"sk_mall/rpc/rpc_seckill/types"

	RMQUtils "sk_mall/utils/RabbitMQ"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var (
	seckillLuaStr string = `
	local stock = redis.call('GET', KEYS[1])
	if stock == false then
	    return 0
	end
	local stock_num = tonumber(stock)
	if stock_num <= 0 then
	    return 1
	end
	local ok = redis.call('sadd', KEYS[2], ARGV[1])
	if ok == 0 then
	    return 2
	end
	redis.call('decr', KEYS[1])
	return 3
	`
)

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type OrderMessage struct {
	OrderNo     int64 `json:"orderNo"`     // 雪花ID
	UserID      int64 `json:"userId"`      // 用户ID
	SKProductID int64 `json:"skProductId"` // 商品ID
	AddrId      int64 `json:"addrId"`      // 地址id
}

func (l *CreateOrderLogic) CreateOrder(in *__.CreateOrderReq) (*__.CreateOrderResp, error) {
	//获取商品evid
	gresp1, gerr1 := l.svcCtx.ProductRpc.GetSKProduct(l.ctx, &product.GetSKProductReq{
		SKProductID: int64(in.SkProductId),
	})
	if gerr1 != nil {
		logc.Errorf(l.ctx, "[productRpc] GetSKProduct err:%s", gerr1.Error())
		return &__.CreateOrderResp{}, gerr1
	}
	if gresp1.Base.Code == 1000 {
		return &__.CreateOrderResp{
			Base: &__.BaseResp{
				Code: 1000,
				Msg:  "商品不存在",
			},
		}, nil
	}
	//活动是否开启
	gresp, err := l.svcCtx.ProductRpc.GetSKEvStatus(l.ctx,
		&product.GetSKEvStatusReq{
			EvId: int64(in.SkProductId),
		})
	if err != nil {
		return nil, err
	}
	if gresp.Base.Code == 1 {
		return &__.CreateOrderResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "活动不存在",
			},
		}, nil
	}
	if gresp.Status == 1 {
		return &__.CreateOrderResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "活动未开始",
			},
		}, nil
	}
	if gresp.Status == 3 {
		return &__.CreateOrderResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "活动已结束",
			},
		}, nil
	}

	//一人一单 lua+redis set

	//已购买的用户集合
	stockKey := fmt.Sprintf("seckill:stock:%d", in.SkProductId)
	setKey := fmt.Sprintf("seckill:purchasedusers:%d", in.SkProductId)
	result, e1 := l.svcCtx.Rds.Eval(seckillLuaStr, []string{stockKey, setKey}, in.UserId)
	if e1 != nil {
		logc.Errorf(l.ctx, "[Redis] lua exec err:%s", e1.Error())
		return &__.CreateOrderResp{}, e1
	}
	resultNum, ok := result.(int64)
	if !ok {
		logc.Errorf(l.ctx, "[Redis] ok assert err")
		return &__.CreateOrderResp{}, errors.New("[Redis] ok assert err")
	} else {
		switch resultNum {
		case 0:
			logc.Errorf(l.ctx, "[Redis] err:key %s not exists", stockKey)
			return &__.CreateOrderResp{}, errors.New("[Redis] stock key not exists")
		case 1:
			return &__.CreateOrderResp{
				Base: &__.BaseResp{
					Code: 1012,
					Msg:  "库存不足",
				},
			}, nil
		case 2:
			return &__.CreateOrderResp{
				Base: &__.BaseResp{
					Code: 1013,
					Msg:  "用户已购买过",
				},
			}, nil
		case 3:
			//购买成功
			//orderNo
			var orderNo int64
			var e3 error
			for {
				orderNo, e3 = l.svcCtx.Node.Generate()
				if orderNo != 0 && e3 == nil {
					break
				}
			}

			//设置订单金额frozen信息到redis, seckill:payment:%d  SkProductId
			paymentKey := fmt.Sprintf("seckill:payment:%d", orderNo)
			//1800s半个小时订单过期
			e6 := l.svcCtx.Rds.SetexCtx(l.ctx, paymentKey, strconv.Itoa(int(in.SkProductId)), 1800)
			if e6 != nil {
				logc.Errorf(l.ctx, "[Redis] set payment key err:%s,orderNo:%d", e6.Error(), orderNo)
				return &__.CreateOrderResp{}, e6
			}
			//发送消息去mq

			var rmqChan *RMQUtils.ChannelWithConfirm
			var e2 error
			for {
				rmqChan, e2 = l.svcCtx.RMQ.Get()
				if errors.Is(e2, RMQUtils.ErrTimeout) {
					logc.Error(l.ctx, "[RMQ] get channel timeout")
					continue
				} else if e2 != nil {
					logc.Errorf(l.ctx, "[RMQ] get channel err:%s,orderNo:%d", e2.Error(), orderNo)
					return &__.CreateOrderResp{}, e2
				}
				break
			}
			//rmqChan.Confirm(false)
			orderMsg := OrderMessage{
				OrderNo:     orderNo,
				UserID:      int64(in.UserId),
				SKProductID: int64(in.SkProductId),
				AddrId:      int64(in.AddrId),
			}
			confirms := *(rmqChan.Confirm)
			//defer close(confirms)
			body, e4 := json.Marshal(orderMsg)
			if e4 != nil {
				logc.Errorf(l.ctx, "[json] marshal err:%s,orderNo:%d", e4.Error(), orderNo)
				return &__.CreateOrderResp{}, e4
			}
			e5 := rmqChan.Channel.PublishWithContext(
				l.ctx,
				"skmall.order.exc", // exchange
				"skmall.order",     // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					ContentType:  "application/json",
					Body:         body,
					DeliveryMode: amqp.Persistent, // 消息持久化
					Timestamp:    time.Now(),
				})
			//todo：这里的rmq如果重启的话，前面的channel全都失效了，优化可以先判断channel是否有用
			if e5 != nil {
				logc.Errorf(l.ctx, "[RMQ] send err:%s,orderNo:%d", e5.Error(), orderNo)
				return &__.CreateOrderResp{}, e5
			}
			select {
			case confirm := <-confirms:
				if confirm.Ack {
					logc.Info(l.ctx, "Message confirmed")
				} else {
					logc.Errorf(l.ctx, "[RMQ] send fail,orderNo:%d", orderNo)
				}
			case <-time.After(5 * time.Second): //超时时间
				logc.Errorf(l.ctx, "[RMQ] confirm timeout,orderNo:%d", orderNo)
				//超时直接关闭，不要放回channel池
				rmqChan.Channel.Close()
				return &__.CreateOrderResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "RMQ超时",
					},
				}, nil
			}
			l.svcCtx.RMQ.Put(rmqChan)
		}
	}

	return &__.CreateOrderResp{
		Base: &__.BaseResp{
			Code: 1,
			Msg:  "success",
		},
	}, nil
}
