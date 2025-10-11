package main

import (
	"flag"
	"fmt"

	"sk_mall/rpc/cron/cache_loader/internal/config"
	"sk_mall/rpc/cron/cache_loader/internal/server"
	"sk_mall/rpc/cron/cache_loader/internal/svc"
	"sk_mall/rpc/cron/cache_loader/tasks"
	"sk_mall/rpc/cron/cache_loader/types"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/cacheloader.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		__.RegisterCacheLoaderServer(grpcServer, server.NewCacheLoaderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	scanner := tasks.Scanner{SvcCtx: ctx}
	scanner.Run()
	loader := tasks.Loader{SvcCtx: ctx}
	loader.Run()

	c1 := cron.New()
	c1.AddJob("00 1 * * *", &scanner) //晚上1点扫描mysql
	c1.AddJob("*/5 * * * *", &loader) //每5min扫描一次redis
	defer c1.Stop()
	c1.Start()

	defer s.Stop()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
