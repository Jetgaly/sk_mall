package logic

import (
	"context"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComReduceSkProductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewComReduceSkProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComReduceSkProductStockLogic {
	return &ComReduceSkProductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ComReduceSkProductStockLogic) ComReduceSkProductStock(in *__.ReduceSkProductStockReq) (*__.ReduceSkProductStockResp, error) {
	// todo: add your logic here and delete this line

	return &__.ReduceSkProductStockResp{}, nil
}
