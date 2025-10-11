package logic

import (
	"context"
	"database/sql"
	"fmt"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReduceSkProductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReduceSkProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReduceSkProductStockLogic {
	return &ReduceSkProductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReduceSkProductStockLogic) ReduceSkProductStock(in *__.ReduceSkProductStockReq) (*__.ReduceSkProductStockResp, error) {

	barrier, _ := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, _ := l.svcCtx.DBConn.RawDB()
	err := barrier.CallWithDB(db, func(tx *sql.Tx) error {
		// 在事务中执行扣减
		r, e := tx.Exec(`update seckill_products set seckill_stock = seckill_stock - 1
						 where id = ? && seckill_stock > 0 `, in.SKProductId)
		if e != nil {
			return e
		}
		c, _ := r.RowsAffected()
		if c == 0 {
			return fmt.Errorf("SKProductId:%d 库存不足", in.SKProductId)
		}
		return nil
	})
	if err != nil {
		logc.Errorf(l.ctx, "[DTM] ReduceSkProductStock err:%s ", err.Error())
		return &__.ReduceSkProductStockResp{}, err
	}
	return &__.ReduceSkProductStockResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
