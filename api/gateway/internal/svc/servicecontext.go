package svc

import (
	"sk_mall/api/gateway/internal/config"
	"sk_mall/rpc/aggr_order/aggrorder"
	"sk_mall/rpc/rpc_merchant/merchant"
	"sk_mall/rpc/rpc_payment/payment"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_seckill/seckill"
	"sk_mall/rpc/rpc_user/user"

	"sk_mall/api/gateway/internal/middleware"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// 图片白名单
var (
	EnableImageList = []string{
		"jpg",
		"jpeg",
		"png",
		"ico",
		"tiff",
		"gif",
		"svg",
		"webg",
		"webp",
	}
)

type ServiceContext struct {
	Config            config.Config
	JwtAuthMiddleware rest.Middleware
	UserRpc           user.User
	ProductRpc        product.Product
	MerchantRpc       merchant.Merchant
	AggrOrderRpc      aggrorder.AggrOrder
	SecKillRpc        seckill.Seckill
	PaymentRpc        payment.Payment
}

func NewServiceContext(c config.Config) *ServiceContext {
	u := user.NewUser(zrpc.MustNewClient(c.UserRpcConf))
	return &ServiceContext{
		Config:            c,
		UserRpc:           u,
		ProductRpc:        product.NewProduct(zrpc.MustNewClient(c.ProductRpcConf)),
		JwtAuthMiddleware: middleware.NewJwtAuthMiddleware(u).Handle,
		MerchantRpc:       merchant.NewMerchant(zrpc.MustNewClient(c.MerchantRpcConf)),
		AggrOrderRpc:      aggrorder.NewAggrOrder(zrpc.MustNewClient(c.AggrOrderRpcConf)),
		SecKillRpc:        seckill.NewSeckill(zrpc.MustNewClient(c.SecKillRpcConf)),
		PaymentRpc:        payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
	}
}
