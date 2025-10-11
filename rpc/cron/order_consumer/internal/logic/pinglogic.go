package logic

import (
	"context"

	"sk_mall/rpc/cron/order_consumer/internal/svc"
	"sk_mall/rpc/cron/order_consumer/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *__.Empty) (*__.Pong, error) {
	// todo: add your logic here and delete this line

	return &__.Pong{}, nil
}
