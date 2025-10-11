package svc

import (
	"sk_mall/rpc/aggr_order/internal/config"
	"sk_mall/rpc/rpc_order/order"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	OrderRpc order.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		OrderRpc: order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
