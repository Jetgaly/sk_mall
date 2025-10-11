package logic

import (
	"context"
	"database/sql"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReduceFrozenBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReduceFrozenBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReduceFrozenBalanceLogic {
	return &ReduceFrozenBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReduceFrozenBalanceLogic) ReduceFrozenBalance(in *__.ReduceFrozenBalanceReq) (*__.ReduceFrozenBalanceResp, error) {

	barrier, _ := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, _ := l.svcCtx.DBConn.RawDB()
	err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 在事务中执行扣减
		_, e := tx.Exec(`update user_wallets set frozen_balance = frozen_balance - COALESCE(
        		(SELECT price FROM frozen_log WHERE id = ?), 0 ) where user_id = ?`, in.OrderNo, in.UserId)
		return e
	})
	if err != nil {
		logc.Errorf(l.ctx, "[DTM] reduce frozen balance err:%s ,orderNo:%d", err.Error(), in.OrderNo)
		return &__.ReduceFrozenBalanceResp{}, err
	}
	return &__.ReduceFrozenBalanceResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
