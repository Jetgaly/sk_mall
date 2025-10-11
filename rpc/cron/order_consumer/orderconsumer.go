package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	"sk_mall/rpc/cron/order_consumer/consumer"
	"sk_mall/rpc/cron/order_consumer/internal/config"
	"sk_mall/rpc/cron/order_consumer/internal/server"
	"sk_mall/rpc/cron/order_consumer/internal/svc"
	"sk_mall/rpc/cron/order_consumer/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/orderconsumer.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	consumers := consumer.Consumers{
		SvcCtx:  ctx,
		Conn:    ctx.RMQConn,
		Queue:   "skmall.order.mq",
		Handler: consumer.MsgHandler,
		Count:   3,
	}
	timeoutConsumers := consumer.Consumers{
		SvcCtx:  ctx,
		Conn:    ctx.RMQConn,
		Queue:   "sk.order.timeout",
		Handler: consumer.TimeoutHandler,
		Count:   3,
	}
	proc.AddShutdownListener(func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			consumers.Stop()
		}()
		go func() {
			defer wg.Done()
			timeoutConsumers.Stop()
		}()

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(40 * time.Second):
			logc.Errorf(context.Background(), "消费者超时停止")
		}
	})

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterOrderConsumerServer(grpcServer, server.NewOrderConsumerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	defer s.Stop()

	consumers.Start()
	timeoutConsumers.Start()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
