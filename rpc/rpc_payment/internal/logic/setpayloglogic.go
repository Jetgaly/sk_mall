package logic

import (
	"context"

	"sk_mall/rpc/rpc_payment/internal/svc"
	"sk_mall/rpc/rpc_payment/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetPayLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetPayLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetPayLogLogic {
	return &SetPayLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetPayLogLogic) SetPayLog(in *__.SetPayLogReq) (*__.SetPayLogResp, error) {
	// todo: add your logic here and delete this line
	sql := `insert into payment_logs(order_no,dtm_gid) values(?,?)`
	_, e := l.svcCtx.DBConn.ExecCtx(l.ctx, sql, in.OrderNo, in.Gid)
	return &__.SetPayLogResp{}, e
}
