package logic

import (
	"context"
	"strconv"

	"sk_mall/rpc/aggr_order/internal/svc"
	"sk_mall/rpc/aggr_order/types"
	"sk_mall/rpc/rpc_order/order"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderInfoLogic {
	return &GetOrderInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderInfoLogic) GetOrderInfo(in *__.GetOrderInfoReq) (*__.GetOrderInfoResp, error) {
	orderNo, _ := strconv.Atoi(in.OrderNo)
	orderResp, e1 := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderReq{
		OrderId: int64(orderNo),
		UserId:  in.UserId,
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[OrderRpc] GetOrder err:%s", e1.Error())
		return &__.GetOrderInfoResp{}, e1
	}
	if orderResp.Base.Code == 10 {
		return &__.GetOrderInfoResp{
			Base: &__.BaseResp{
				Code: 10,
				Msg:  "订单不存在",
			},
		}, nil
	}
	//todo:其他信息查询聚合
	return &__.GetOrderInfoResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Info: &__.OrderInfo{
			OrderNo: int64(orderNo),
			Addr: &__.AddrInfo{
				AddrId: orderResp.Info.AddrId,
			},
			SkProduck: &__.SkProductInfo{
				SkProductId: orderResp.Info.SkProduckId,
			},
			Quantity:    orderResp.Info.Quantity,
			UnitPrice:   orderResp.Info.UnitPrice,
			TotalAmount: orderResp.Info.TotalAmount,
			CreatedAt:   orderResp.Info.CreatedAt,
			ExpireAt:    orderResp.Info.ExpireAt,
		},
	}, nil
}
