package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpcConf      zrpc.RpcClientConf
	ProductRpcConf   zrpc.RpcClientConf
	MerchantRpcConf  zrpc.RpcClientConf
	AggrOrderRpcConf zrpc.RpcClientConf
	SecKillRpcConf   zrpc.RpcClientConf
	PaymentRpcConf   zrpc.RpcClientConf
}
