package logic

import (
	"context"
	"strconv"

	"fmt"

	"sk_mall/rpc/cron/order_consumer/orderconsumer"
	"sk_mall/rpc/rpc_merchant/merchant"
	"sk_mall/rpc/rpc_payment/internal/svc"
	"sk_mall/rpc/rpc_payment/payment"
	"sk_mall/rpc/rpc_payment/types"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_user/user"

	"github.com/dtm-labs/client/dtmgrpc"
	_ "github.com/dtm-labs/driver-gozero"
	"github.com/go-redsync/redsync/v4"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayOrderLogic) PayOrder(in *__.PayOrderReq) (*__.PayOrderResp, error) {
	// todo: add your logic here and delete this line
	paymentKey := fmt.Sprintf("seckill:payment:%d", in.OrderNo)
	skProductId, e1 := l.svcCtx.Rds.GetCtx(l.ctx, paymentKey)
	fmt.Println(paymentKey, ":", skProductId)
	if e1 != nil {
		logc.Errorf(l.ctx, "[Redis] get seckill:payment:? err:%s", e1.Error())
		return &__.PayOrderResp{}, e1
	}
	if skProductId == "" {
		return &__.PayOrderResp{
			Base: &__.BaseResp{
				Code: 2001,
				Msg:  "订单过期",
			},
		}, nil
	}
	skProductIdNum, _ := strconv.Atoi(skProductId)
	//获取分布式锁，防止重复支付
	lockCtx, cancel := context.WithCancel(l.ctx)
	lockname := fmt.Sprintf("lock:order:%d", in.OrderNo)
	mutex, e2 := l.svcCtx.RLCreater.GetLock(lockCtx, lockname, redsync.WithTries(3))
	if e2 != nil {
		//拿不到锁
		cancel()
		return &__.PayOrderResp{
			Base: &__.BaseResp{
				Code: 2002,
				Msg:  "订单正在支付",
			},
		}, nil
	}
	//释放锁
	defer l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)

	gresp, err := l.svcCtx.UserRpc.FrozenBalance(l.ctx, &user.FrozenBalanceReq{UserId: int64(in.UserId), SkProductId: int64(skProductIdNum), OrderNo: in.OrderNo})
	if err != nil {
		return &__.PayOrderResp{}, err
	}
	switch gresp.Base.Code {
	case 1001:
		return &__.PayOrderResp{
			Base: &__.BaseResp{
				Code: 1001,
				Msg:  "余额不足",
			},
		}, nil
	case 1011:
		return &__.PayOrderResp{
			Base: &__.BaseResp{
				Code: 1011,
				Msg:  "订单已支付",
			},
		}, nil
	case 0:
		//Dtm
		gid := dtmgrpc.MustGenGid(l.svcCtx.Config.DTM)
		userRpcServer, e3 := l.svcCtx.Config.UserRpcConf.BuildTarget()
		if e3 != nil {
			logc.Errorf(l.ctx, "[DTM] l.svcCtx.Config.UserRpcConf.BuildTarget() err:%s", e3.Error())
			return &__.PayOrderResp{}, e3
		}
		productRpcServer, e4 := l.svcCtx.Config.ProductRpcConf.BuildTarget()
		if e4 != nil {
			logc.Errorf(l.ctx, "[DTM]  l.svcCtx.Config.ProductRpcConf.BuildTarget() err:%s", e4.Error())
			return &__.PayOrderResp{}, e3
		}
		orderRpcServer, e5 := l.svcCtx.Config.OrderRpcConf.BuildTarget()
		if e5 != nil {
			logc.Errorf(l.ctx, "[DTM]  l.svcCtx.Config.OrderRpcConf.BuildTarget() err:%s", e5.Error())
			return &__.PayOrderResp{}, e5
		}
		merchantRpcServer, e6 := l.svcCtx.Config.MerchantRpcConf.BuildTarget()
		if e6 != nil {
			logc.Errorf(l.ctx, "[DTM]  l.svcCtx.Config.MerchantRpcConf.BuildTarget() err:%s", e6.Error())
			return &__.PayOrderResp{}, e6
		}
		paymentRpcServer, e8 := l.svcCtx.Config.PaymentRpcConf.BuildTarget()
		if e8 != nil {
			logc.Errorf(l.ctx, "[DTM]  l.svcCtx.Config.PaymentRpcConf.BuildTarget() err:%s", e8.Error())
			return &__.PayOrderResp{}, e8
		}
		//记录gid和orderId
		paymentReq := &payment.SetPayLogReq{
			Gid:     gid,
			OrderNo: in.OrderNo,
		}
		//用户frozen余额扣减+frozen log//系统宕机可以根据这个log来恢复，同时这个log实现了幂等性
		userReq := &user.ReduceFrozenBalanceReq{
			OrderNo: in.OrderNo,
			UserId:  int64(in.UserId),
		}
		//扣减mysql库存
		productReq := &product.ReduceSkProductStockReq{
			SKProductId: int64(skProductIdNum),
		}

		//商户金额增加
		merchantReq := &merchant.IncreaseBalanceReq{
			OrderNo:     uint64(in.OrderNo),
			TotalAmount: gresp.TotalAmount,
			MerchantId:  gresp.MerchantId,
		}
		//修改订单状态
		orderReq := &orderconsumer.SetOrderHasPaidReq{
			OrderNo: in.OrderNo,
		}

		//开启dtm事务
		saga := dtmgrpc.NewSagaGrpc(l.svcCtx.Config.DTM, gid).
			Add(paymentRpcServer+"/payment.Payment/SetPayLog", paymentRpcServer+"/payment.Payment/ComSetPayLog", paymentReq).
			Add(userRpcServer+"/user.User/ReduceFrozenBalance", userRpcServer+"/user.User/ComReduceFrozenBalance", userReq).
			Add(productRpcServer+"/product.Product/ReduceSkProductStock", productRpcServer+"/product.Product/ComReduceSkProductStock", productReq).
			Add(merchantRpcServer+"/merchant.Merchant/IncreaseBalance", merchantRpcServer+"/merchant.Merchant/ComIncreaseBalance", merchantReq).
			Add(orderRpcServer+"/order_consumer.OrderConsumer/SetOrderHasPaid", orderRpcServer+"/order_consumer.OrderConsumer/SetOrderHasPaid", orderReq)
		retry := 0
		var e7 error
		//重试三次
		for retry <= 3 {
			e7 = saga.Submit()
			if e7 == nil {
				break
			}
			retry++
		}
		if e7 != nil {
			logc.Errorf(l.ctx, "[DTM] saga submit err:%s, orderNo: %d", e7.Error(), in.OrderNo)
			return &__.PayOrderResp{}, e7
		}
	}

	return &__.PayOrderResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
