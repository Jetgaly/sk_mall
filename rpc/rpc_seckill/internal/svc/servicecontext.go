package svc

import (
	"context"
	"database/sql"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_seckill/internal/config"
	"sk_mall/rpc/rpc_user/user"
	"sk_mall/utils"
	"sk_mall/utils/RabbitMQ"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	gr "github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config     config.Config
	DBConn     sqlx.SqlConn
	Rds        *redis.Redis
	Node       *utils.SafeSnowFlakeCreater
	RLCreater  *utils.RedLockCreater
	RMQ        *RMQUtils.RMQ
	Queue      amqp.Queue
	UserRpc    user.User
	ProductRpc product.Product
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

	//redis:
	rds := redis.MustNewRedis(c.Redis.RedisConf)
	redis_flag := rds.Ping()
	if !redis_flag {
		logc.Error(context.Background(), "[Redis]redis连接失败")
		panic("[Redis]redis连接失败")
	}

	//snowflakes
	nodeID, e1 := rds.Incr("SKMall:NodeId")
	if e1 != nil {
		logc.Errorf(context.Background(), "[Redis]SKMall:NodeId获取失败,%s", e1.Error())
		panic("[Redis]SKMall:NodeId获取失败")
	}
	node, e2 := utils.NewSafeSnowFlakeCreater(nodeID, time.Duration(500)*time.Millisecond)
	if e2 != nil {
		logc.Errorf(context.Background(), "[SnowFlake]初始化失败,%s", e2.Error())
		panic("[SnowFlake]初始化失败")
	}

	//RedLock
	var clis []*gr.Client
	for i, v := range c.DLockRedis.Hosts {
		r := gr.NewClient(&gr.Options{
			Addr:     v,
			Password: c.DLockRedis.Passes[i],
		})
		t := r.Ping(context.Background())
		if t.Err() != nil {
			logc.Errorf(context.Background(), "[RedLock]初始化失败,%s", t.Err().Error())
			panic("[RedLock]初始化失败")
		}
		clis = append(clis, r)
	}
	rlc, e3 := utils.NewRedLockCreater(clis)
	if e3 != nil {
		logc.Errorf(context.Background(), "[RedLock]初始化失败,%s", e3.Error())
		panic("[RedLock]初始化失败")
	}

	//RMQ
	r, e4 := RMQUtils.NewRMQ(c.RMQ.Dsn, 9, 180, 100)
	if e4 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e4.Error())
		panic("[RMQ]初始化失败")
	}
	channelWithConfirm, e5 := r.Get()
	channel := channelWithConfirm.Channel
	defer r.Put(channelWithConfirm)
	if e5 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e5.Error())
		panic("[RMQ]初始化失败")
	}
	e5 = channel.ExchangeDeclare(
		"skmall.order.exc",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if e5 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e5.Error())
		panic("[RMQ]初始化失败")
	}
	q, e6 := channel.QueueDeclare(
		"skmall.order.mq",
		true,
		false,
		false, //允许多个消费者连接队列
		false,
		nil,
	)

	/*
			"x-max-length": 10000, // 最大消息数
		    "x-overflow": "reject-publish", // 满时拒绝新消息
	*/
	if e6 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e6.Error())
		panic("[RMQ]初始化失败")
	}
	e5 = channel.QueueBind(
		"skmall.order.mq",
		"skmall.order",
		"skmall.order.exc",
		false,
		nil,
	)
	if e5 != nil {
		logc.Errorf(context.Background(), "[RMQ]初始化失败,%s", e5.Error())
		panic("[RMQ]初始化失败")
	}

	//延迟启动解决冲突相互依赖

	return &ServiceContext{
		Config:     c,
		DBConn:     conn,
		Rds:        rds,
		Node:       node,
		RLCreater:  rlc,
		RMQ:        r,
		Queue:      q,
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
		ProductRpc: product.NewProduct(zrpc.MustNewClient(c.ProductRpcConf)),
	}
}
