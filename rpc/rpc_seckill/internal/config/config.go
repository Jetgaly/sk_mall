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
	DLockRedis struct {
		Hosts  []string
		Passes []string
	}

	RMQ struct {
		Dsn string
	}
	UserRpcConf zrpc.RpcClientConf
	ProductRpcConf zrpc.RpcClientConf
}
