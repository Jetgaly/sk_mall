package logic

import (
	"context"
	"database/sql"
	"errors"

	"sk_mall/rpc/cron/order_consumer/internal/svc"
	"sk_mall/rpc/cron/order_consumer/types"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetOrderHasPaidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetOrderHasPaidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetOrderHasPaidLogic {
	return &SetOrderHasPaidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	ErrOrderNotExists error = errors.New("订单还未生成")
)

func (l *SetOrderHasPaidLogic) SetOrderHasPaid(in *__.SetOrderHasPaidReq) (*__.SetOrderHasPaidResp, error) {
	barrier, _ := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, _ := l.svcCtx.DBConn.RawDB()
	err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		r, e := tx.Exec(`update sk_orders set pay_status = 1, status = 1 where order_no = ?`, in.OrderNo)
		if e != nil {
			return e
		}
		c, _ := r.RowsAffected()
		if c == 0 {
			return ErrOrderNotExists
		}
		return nil
	})
	if err != nil {
		if err != ErrOrderNotExists {
			logc.Errorf(l.ctx, "[DTM] SetOrderHasPaid err:%s ", err.Error())
		}
		return &__.SetOrderHasPaidResp{}, err
	}
	return &__.SetOrderHasPaidResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil

}
