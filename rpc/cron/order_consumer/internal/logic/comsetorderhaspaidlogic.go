package logic

import (
	"context"

	"sk_mall/rpc/cron/order_consumer/internal/svc"
	"sk_mall/rpc/cron/order_consumer/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComSetOrderHasPaidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewComSetOrderHasPaidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComSetOrderHasPaidLogic {
	return &ComSetOrderHasPaidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ComSetOrderHasPaidLogic) ComSetOrderHasPaid(in *__.SetOrderHasPaidReq) (*__.SetOrderHasPaidResp, error) {
	// todo: add your logic here and delete this line

	return &__.SetOrderHasPaidResp{}, nil
}
