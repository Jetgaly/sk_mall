package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComReduceFrozenBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewComReduceFrozenBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComReduceFrozenBalanceLogic {
	return &ComReduceFrozenBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ComReduceFrozenBalanceLogic) ComReduceFrozenBalance(in *__.ReduceFrozenBalanceReq) (*__.ReduceFrozenBalanceResp, error) {
	//dtm事务是最终完成的，不会触发回滚

	return &__.ReduceFrozenBalanceResp{}, nil
}
