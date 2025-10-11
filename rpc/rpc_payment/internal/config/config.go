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

	UserRpcConf     zrpc.RpcClientConf
	ProductRpcConf  zrpc.RpcClientConf
	OrderRpcConf    zrpc.RpcClientConf
	MerchantRpcConf zrpc.RpcClientConf
	PaymentRpcConf  zrpc.RpcClientConf
	DTM             string
}
