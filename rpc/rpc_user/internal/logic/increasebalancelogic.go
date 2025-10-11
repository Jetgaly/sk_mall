package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IncreaseBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIncreaseBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncreaseBalanceLogic {
	return &IncreaseBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IncreaseBalanceLogic) IncreaseBalance(in *__.IncreaseBalanceReq) (*__.IncreaseBalanceResq, error) {
	// todo: add your logic here and delete this line

	return &__.IncreaseBalanceResq{}, nil
}
