package logic

import (
	"context"
	"errors"
	"time"

	"sk_mall/rpc/rpc_order/internal/svc"
	"sk_mall/rpc/rpc_order/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type orderListElement struct {
	OrderNo     int64     `db:"order_no"`
	SkProduckId int64     `db:"sk_product_id"`
	Quantity    int64     `db:"quantity"`
	UnitPrice   string    `db:"unit_price"`
	TotalAmount string    `db:"total_amount"`
	Status      int64     `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	ExpireAt    time.Time `db:"expire_time"`
}

func (l *GetOrderListLogic) GetOrderList(in *__.GetOrderListReq) (*__.GetOrderListResp, error) {
	// todo: add your logic here and delete this line
	var list []orderListElement
	sql := `
		select order_no,sk_product_id,quantity,unit_price,total_amount,status,created_at,expire_time
		from sk_orders 
		where user_id = ?
		order by created_at desc 
		limit ?,?
	`
	offset := (in.Page - 1) * in.Limit
	if offset < 0 {
		offset = 0
	}
	limit := in.Limit
	if in.Limit > 100{
		limit = 100
	}
	e1 := l.svcCtx.DBConn.QueryRowsCtx(l.ctx, &list, sql, in.UserId, offset, limit)
	if errors.Is(e1, sqlx.ErrNotFound) {
		return &__.GetOrderListResp{
			Base: &__.BaseResp{
				Code: 5000,
				Msg:  "暂无订单",
			},
		}, nil
	}
	if e1 != nil {
		logc.Errorf(l.ctx, "[DBConn] query err:%s", e1.Error())
		return &__.GetOrderListResp{}, e1
	}
	var orderlist []*__.OrderListInfo
	for _, e := range list {
		orderlist = append(orderlist, &__.OrderListInfo{
			OrderNo:     e.OrderNo,
			SkProduckId: e.SkProduckId,
			Quantity:    e.Quantity,
			UnitPrice:   e.UnitPrice,
			TotalAmount: e.TotalAmount,
			Status:      int32(e.Status),
			CreatedAt:   timestamppb.New(e.CreatedAt),
			ExpireAt:    timestamppb.New(e.ExpireAt),
		})
	}
	return &__.GetOrderListResp{
		Base: &__.BaseResp{
			Code: 1,
			Msg:  "success",
		},
		List: orderlist,
	}, nil
}
