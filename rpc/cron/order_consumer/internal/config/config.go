package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource  string
		MaxOpen     int
		MaxIdle     int
		MaxLifetime int
	}
	RMQ struct {
		Dsn string
	}
	ProductRpcConf zrpc.RpcClientConf
	UserRpcConf    zrpc.RpcClientConf
}
