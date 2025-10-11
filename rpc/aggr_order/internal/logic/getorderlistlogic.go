package logic

import (
	"context"

	"sk_mall/rpc/aggr_order/internal/svc"
	"sk_mall/rpc/aggr_order/types"
	"sk_mall/rpc/rpc_order/order"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderListLogic) GetOrderList(in *__.GetOrderListReq) (*__.GetOrderListResp, error) {
	// todo: add your logic here and delete this line
	gresp, err := l.svcCtx.OrderRpc.GetOrderList(l.ctx, &order.GetOrderListReq{
		UserId: in.UserId,
		Limit:  in.Limit,
		Page:   in.Page,
	})
	if err != nil {
		logc.Errorf(l.ctx, "[OrderRpc] GetOrderList err:%s", err.Error())
		return &__.GetOrderListResp{}, err
	}
	var list []*__.OrderListElem
	for _, e := range gresp.List {
		list = append(list, &__.OrderListElem{
			OrderNo:     e.OrderNo,
			SkProduck:   &__.SkProductInfo{SkProductId: e.SkProduckId},
			Quantity:    e.Quantity,
			UnitPrice:   e.UnitPrice,
			TotalAmount: e.TotalAmount,
			Status:      e.Status,
			CreatedAt:   e.CreatedAt,
			ExpireAt:    e.ExpireAt,
		})
	}
	return &__.GetOrderListResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		List: list,
	}, nil
}
