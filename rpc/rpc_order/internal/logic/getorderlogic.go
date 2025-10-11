package logic

import (
	"context"
	"errors"
	"time"

	"sk_mall/rpc/rpc_order/internal/svc"
	"sk_mall/rpc/rpc_order/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type orderInfo struct {
	OrderNo     int64     `db:"order_no"`
	AddrId      int64     `db:"addr_id"`
	SkProduckId int64     `db:"sk_product_id"`
	Quantity    int64     `db:"quantity"`
	UnitPrice   string    `db:"unit_price"`
	TotalAmount string    `db:"total_amount"`
	Status      int64     `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	ExpireAt    time.Time `db:"expire_time"`
}

func (l *GetOrderLogic) GetOrder(in *__.GetOrderReq) (*__.GetOrderResp, error) {
	// todo: add your logic here and delete this line
	var info orderInfo
	sql := `
		select order_no,addr_id,sk_product_id,quantity,unit_price,total_amount,status,created_at,expire_time
		from sk_orders 
		where order_no = ? and user_id = ?
	`
	err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &info, sql, in.OrderId, in.UserId)
	if errors.Is(err, sqlx.ErrNotFound) {
		return &__.GetOrderResp{
			Base: &__.BaseResp{
				Code: 10,
				Msg:  "订单不存在",
			},
		}, nil
	}
	if err != nil {
		return &__.GetOrderResp{}, err
	}

	return &__.GetOrderResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Info: &__.OrderInfo{
			OrderNo:     info.OrderNo,
			AddrId:      info.AddrId,
			SkProduckId: info.SkProduckId,
			Quantity:    info.Quantity,
			UnitPrice:   info.UnitPrice,
			TotalAmount: info.TotalAmount,
			Status:      int32(info.Status),
			CreatedAt:   timestamppb.New(info.CreatedAt),
			ExpireAt:    timestamppb.New(info.ExpireAt),
		},
	}, nil
}
