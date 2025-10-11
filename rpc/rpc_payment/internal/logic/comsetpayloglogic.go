package logic

import (
	"context"

	"sk_mall/rpc/rpc_payment/internal/svc"
	"sk_mall/rpc/rpc_payment/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ComSetPayLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewComSetPayLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ComSetPayLogLogic {
	return &ComSetPayLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ComSetPayLogLogic) ComSetPayLog(in *__.SetPayLogReq) (*__.SetPayLogResp, error) {
	// todo: add your logic here and delete this line

	return &__.SetPayLogResp{}, nil
}
