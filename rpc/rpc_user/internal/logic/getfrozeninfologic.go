package logic

import (
	"context"
	"database/sql"
	"errors"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetFrozenInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFrozenInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFrozenInfoLogic {
	return &GetFrozenInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type frozenLog struct {
	Price sql.NullString `db:"price"`
}

func (l *GetFrozenInfoLogic) GetFrozenInfo(in *__.GetFrozenInfoReq) (*__.GetFrozenInfoResp, error) {
	// todo: add your logic here and delete this line
	sql := `select * from frozen_log where id = ?`
	var info frozenLog
	err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &info, sql, in.OrderNo)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logc.Errorf(l.ctx, "[DBConn] query err: %s", err.Error())
		return &__.GetFrozenInfoResp{}, err
	}
	if errors.Is(err, sqlx.ErrNotFound) {
		return &__.GetFrozenInfoResp{
			Base: &__.BaseResp{
				Code: 3001,
				Msg:  "frozenlog not exists",
			},
		}, nil
	}
	return &__.GetFrozenInfoResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
