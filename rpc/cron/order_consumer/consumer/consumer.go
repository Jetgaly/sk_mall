package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"fmt"
	"sk_mall/rpc/cron/order_consumer/internal/svc"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_user/user"
	RMQUtils "sk_mall/utils/RabbitMQ"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Consumers struct {
	SvcCtx  *svc.ServiceContext
	Conn    *amqp.Connection
	Queue   string
	Handler func(SvcCtx *svc.ServiceContext, msg []byte) error
	Count   int //消费者数量
	Wg      sync.WaitGroup
	Cancel  context.CancelFunc
}

// Start 启动多个消费者协程
func (c *Consumers) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	c.Cancel = cancel

	for i := 0; i < c.Count; i++ {
		c.Wg.Add(1)
		go c.worker(ctx, i)
	}

	logc.Infof(context.Background(), "启动 %d 个消费者协程，队列: %s", c.Count, c.Queue)
	return nil
}

// NewConsumer 实例化一个消费者, 会单独用一个channel，设置每次只取一个消息
func (c *Consumers) worker(ctx context.Context, no int) {
	defer c.Wg.Done()

	logc.Infof(ctx, "worker:%d starts", no)

	err := c.consumeMessage(ctx, no)
	if err != nil {
		logc.Errorf(ctx, "worker-%d err: %s", no, err.Error())
		panic("worker init err")
	}

}
func (c *Consumers) consumeMessage(ctx context.Context, no int) error {

	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("new mq channel err: %s", err.Error())
	}
	defer ch.Close()
	// 设置 QoS：每次只取一个消息处理
	// 参数说明：
	// prefetchCount: 每次预取的消息数量，设为1表示每次只取一个
	// prefetchSize: 预取的消息总大小（字节），0表示不限制
	// global: 是否在连接级别应用，false表示只在当前channel生效
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("set qos err: %v", err)
	}

	deliveries, err := ch.Consume(c.Queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume err: %v, queue: %s", err, c.Queue)
	}

	for {
		select {
		case <-ctx.Done():
			logc.Infof(ctx, "worker:%d stop", no)
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				logc.Error(ctx, "msg queue has closed")
			} else {
				err = c.Handler(c.SvcCtx, delivery.Body)
				if err != nil {
					_ = delivery.Reject(true) // 处理失败，重新入队
					//logc
					logc.Errorf(context.Background(), "msg handler err:%s", err.Error())

				} else {
					_ = delivery.Ack(false) // 处理成功，确认消息
				}
			}
		}
	}
}

// Stop 优雅停止消费者
func (c *Consumers) Stop() {
	defer c.Conn.Close()
	logc.Info(context.Background(), "stopping consumers gracefully")

	if c.Cancel != nil {
		c.Cancel() // 发送关闭信号
	}

	// 等待所有消费者协程完成
	done := make(chan struct{})
	go func() {
		c.Wg.Wait()
		close(done)
	}()

	// 设置超时
	select {
	case <-done:
		logc.Info(context.Background(), "consumers have stopped")
	case <-time.After(30 * time.Second):
		logc.Info(context.Background(), "consumers stopping timeout 30 sec, quit forcefuly")
	}

}

type OrderMessage struct {
	OrderNo     int64 `json:"orderNo"`     // 雪花ID
	UserID      int64 `json:"userId"`      // 用户ID
	SKProductID int64 `json:"skProductId"` // 商品ID
	AddrId      int64 `json:"addrId"`      // 地址id
}
type DelayOrderMessage struct {
	OrderNo     int64 `json:"orderNo"`     // 雪花ID
	UserID      int64 `json:"userId"`      // 用户ID
	SKProductID int64 `json:"skProductId"` // 商品ID
}

func MsgHandler(SvcCtx *svc.ServiceContext, msg []byte) error {
	var orderMsg OrderMessage
	e1 := json.Unmarshal(msg, &orderMsg)
	if e1 != nil {
		logc.Errorf(context.Background(), "[json] unmarshal err:%s", e1.Error())
		return e1
	}
	//查询sk商品单价
	gresp, e2 := SvcCtx.ProductRpc.GetSKProduct(context.Background(), &product.GetSKProductReq{
		SKProductID: orderMsg.SKProductID,
	})
	if e2 != nil {
		logc.Errorf(context.Background(), "[ProductRpc] err:%s", e2.Error())
		return e2
	}
	//生成订单
	err := SvcCtx.DBConn.Transact(func(session sqlx.Session) error {
		sql := "insert into sk_orders(order_no,user_id,addr_id,sk_product_id,unit_price,total_amount,expire_time) values(?,?,?,?,?,?,?)"
		_, e3 := session.Exec(sql, orderMsg.OrderNo, orderMsg.UserID, orderMsg.AddrId, orderMsg.SKProductID, gresp.Info.SeckillPrice, gresp.Info.SeckillPrice, time.Now().Add(5*time.Minute))
		return e3
	})
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number != 1062 {
			logc.Errorf(context.Background(), "[DBConn] create order err:%s, orderNo:%d", err.Error(), orderMsg.OrderNo)
			return err
		}
	}
	//延迟队列：

	var channel *RMQUtils.ChannelWithConfirm
	var cerr error
	for {
		channel, cerr = SvcCtx.RMQ.Get()
		if errors.Is(cerr, RMQUtils.ErrTimeout) {
			logc.Error(context.Background(), "[RMQ] get channel timeout")
			continue
		} else if cerr != nil {
			logc.Errorf(context.Background(), "[RMQ] publish delayexc err:%s, orderNo:%d", cerr.Error(), orderMsg.OrderNo)
			return nil
		}
		break
	}
	delayOrderMsg := DelayOrderMessage{
		OrderNo:     orderMsg.OrderNo,
		SKProductID: orderMsg.SKProductID,
		UserID:      orderMsg.UserID,
	}
	delayMsgData, jerr := json.Marshal(delayOrderMsg)
	if jerr != nil {
		logc.Errorf(context.Background(), "[json] publish delayexc marshal err:%s orderNo:%d", jerr.Error(), orderMsg.OrderNo)
		return nil
	}
	e4 := channel.Channel.PublishWithContext(context.Background(),
		"sk.order.delayexc",
		"sk.order.delay",
		false,
		false,
		amqp.Publishing{
			Body: delayMsgData,
			//Expiration: "5000", // 消息级别TTL: 5秒 优先级更高
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // 消息持久化
			Timestamp:    time.Now(),
		})
	if e4 != nil {
		logc.Errorf(context.Background(), "[RMQ] publish delayexc err:%s, orderNo:%d", e4.Error(), orderMsg.OrderNo)
		return nil
	}
	confirms := *(channel.Confirm)
	select {
	case confirm := <-confirms:
		if confirm.Ack {
			logc.Info(context.Background(), "Message confirmed")
		} else {
			logc.Errorf(context.Background(), "[RMQ] send fail,orderNo:%d", orderMsg.OrderNo)
		}
	case <-time.After(5 * time.Second): //超时时间
		logc.Errorf(context.Background(), "[RMQ] confirm timeout,orderNo:%d", orderMsg.OrderNo)
		//超时直接关闭，不要放回channel池
		channel.Channel.Close()
		return nil
	}
	SvcCtx.RMQ.Put(channel)

	return nil
}

var comLuaStr string = `
	local ttl = redis.call('TTL', KEYS[1])
	local ttl_num = tonumber(ttl)
	if ttl_num == -2 then
	    return 0
	end
	if ttl_num <= 15 and ttl_num >=0 then
		return 0
	end

	local sttl = redis.call('TTL', KEYS[2])
	local sttl_num = tonumber(sttl)
	if sttl_num == -2 then
	    return 0
	end
	if sttl_num <= 15 and sttl_num >=0 then
		return 0
	end
	redis.call('incr',KEYS[1])
	redis.call('srem',KEYS[2],ARGV[1])
	return 1
`

func TimeoutHandler(SvcCtx *svc.ServiceContext, msg []byte) error {
	var orderMsg DelayOrderMessage
	e1 := json.Unmarshal(msg, &orderMsg)
	fmt.Println(orderMsg)
	//考虑加分布式锁+少量延时
	if e1 != nil {
		logc.Errorf(context.Background(), "[json] unmarshal err:%s", e1.Error())
		return e1
	}
	//检查是否冻结金额（是否支付）
	gresp, e2 := SvcCtx.UserRpc.GetFrozenInfo(context.Background(), &user.GetFrozenInfoReq{
		OrderNo: orderMsg.OrderNo,
	})
	if e2 != nil {
		logc.Errorf(context.Background(), "[UserRpc] GetFrozenInfo err:%s", e2.Error())
		return e2
	}
	if gresp.Base.Code == 0 {
		//已经冻结生成订单，直接返回
		return nil
	}

	//还没支付，真正过期
	//订单设置为已过期
	sql := `update sk_orders set status = 5 where order_no = ? and expire_time <= CURRENT_TIMESTAMP`
	r, e3 := SvcCtx.DBConn.Exec(sql, orderMsg.OrderNo)
	if e3 != nil {
		logc.Errorf(context.Background(), "[DBConn] update expire err:%s, orderNo:%d", e3.Error(), orderMsg.OrderNo)
		return e2
	}
	count, e4 := r.RowsAffected()
	if e4 != nil {
		logc.Errorf(context.Background(), "[DBConn] update expire err:%s, orderNo:%d", e4.Error(), orderMsg.OrderNo)
		return e4
	}
	if count == 0 {
		logc.Errorf(context.Background(), "[DBConn] update expire err:count == 0, orderNo:%d", orderMsg.OrderNo)
		return fmt.Errorf("update expire err,count == 0")
	}
	//redis库存归还，已购买用户去除
	stockKey := fmt.Sprintf("seckill:stock:%d", orderMsg.SKProductID)
	UserSetKey := fmt.Sprintf("seckill:purchasedusers:%d", orderMsg.SKProductID)

	_, e5 := SvcCtx.Rds.Eval(comLuaStr, []string{stockKey, UserSetKey}, orderMsg.UserID)
	if e5 != nil {
		logc.Errorf(context.Background(), "[Redis] lua exec err:%s, orderNo:%d", e5.Error(), orderMsg.OrderNo)
	}
	return e5
}
