package logic

import (
	"context"

	"sk_mall/rpc/rpc_merchant/internal/svc"
	"sk_mall/rpc/rpc_merchant/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComIncreaseBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewComIncreaseBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComIncreaseBalanceLogic {
	return &ComIncreaseBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ComIncreaseBalanceLogic) ComIncreaseBalance(in *__.IncreaseBalanceReq) (*__.IncreaseBalanceResp, error) {
	// todo: add your logic here and delete this line

	return &__.IncreaseBalanceResp{}, nil
}
