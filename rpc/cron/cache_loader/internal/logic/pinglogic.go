package logic

import (
	"context"
	"fmt"
	"time"

	"sk_mall/rpc/cron/cache_loader/internal/svc"
	"sk_mall/rpc/cron/cache_loader/types"

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
	time.Sleep(10 * time.Second)
	fmt.Println("after 10 sec")
	return &__.Pong{
		Msg: "pong",
	}, nil
}
