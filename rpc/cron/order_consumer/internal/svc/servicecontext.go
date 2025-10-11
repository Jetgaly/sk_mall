package svc

import (
	"context"
	"database/sql"

	"sk_mall/rpc/cron/order_consumer/internal/config"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_user/user"
	RMQUtils "sk_mall/utils/RabbitMQ"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	RMQConn    *amqp.Connection
	DBConn     sqlx.SqlConn
	RMQ        *RMQUtils.RMQ
	ProductRpc product.Product
	UserRpc    user.User
	Rds        *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.DB.DataSource)
	if err != nil {
		logc.Error(context.Background(), "数据库连接失败")
		panic(err)
	}

	if c.DB.MaxOpen > 0 {
		db.SetMaxOpenConns(c.DB.MaxOpen)
	}
	if c.DB.MaxIdle > 0 {
		db.SetMaxIdleConns(c.DB.MaxIdle)
	}
	if c.DB.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(c.DB.MaxLifetime) * time.Second)
	}

	// 包装成 sqlx.SqlConn
	conn := sqlx.NewSqlConnFromDB(db)

	rds := redis.MustNewRedis(c.Redis.RedisConf)
	redis_flag := rds.Ping()
	if !redis_flag {
		logc.Error(context.Background(), "[Redis]redis连接失败")
		panic("[Redis]redis连接失败")
	}

	//rmqconsumer
	rconn, e1 := amqp.Dial(c.RMQ.Dsn)
	if e1 != nil {
		logc.Errorf(context.Background(), "[RMQConn] dial err:%s", e1.Error())
	}

	//RMQ
	r, e2 := RMQUtils.NewRMQ(c.RMQ.Dsn, 3, 21, 9)
	if e2 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e2.Error())
		panic("[RMQ]初始化失败")
	}
	channelWithConfirm, e3 := r.Get()
	channel := channelWithConfirm.Channel
	defer r.Put(channelWithConfirm)
	if e3 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e3.Error())
		panic("[RMQ]初始化失败")
	}
	//在延迟队列过期之后会转发给死信交换机-->死信队列
	//死信交换机+死信队列(工作队列)
	e4 := channel.ExchangeDeclare(
		"sk.order.timeoutexc",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if e4 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e4.Error())
		panic("[RMQ]初始化失败")
	}
	_, e5 := channel.QueueDeclare(
		"sk.order.timeout",
		true,
		false,
		false, //允许多个消费者连接队列
		false,
		nil,
	)
	if e5 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e5.Error())
		panic("[RMQ]初始化失败")
	}
	e5 = channel.QueueBind(
		"sk.order.timeout",
		"sk.order.timeout",
		"sk.order.timeoutexc",
		false,
		nil,
	)
	if e5 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e5.Error())
		panic("[RMQ]初始化失败")
	}

	//延迟交换机+延迟队列
	e6 := channel.ExchangeDeclare(
		"sk.order.delayexc",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if e6 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e6.Error())
		panic("[RMQ]初始化失败")
	}
	args := amqp.Table{
		"x-dead-letter-exchange":    "sk.order.timeoutexc", // 死信转发到主交换机
		"x-dead-letter-routing-key": "sk.order.timeout",    // 死信路由键
		"x-message-ttl":             10000,                 // 队列级别TTL: 35.5分钟 2130000ms
	}
	_, e7 := channel.QueueDeclare(
		"sk.order.delay",
		true,
		false,
		false, //允许多个消费者连接队列
		false,
		args,
	)
	if e7 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e7.Error())
		panic("[RMQ]初始化失败")
	}
	e7 = channel.QueueBind(
		"sk.order.delay",
		"sk.order.delay",
		"sk.order.delayexc",
		false,
		nil,
	)
	if e7 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e7.Error())
		panic("[RMQ]初始化失败")
	}

	return &ServiceContext{
		Config:     c,
		RMQConn:    rconn,
		DBConn:     conn,
		RMQ:        r,
		Rds:        rds,
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpcConf)),
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
