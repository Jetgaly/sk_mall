package logic

import (
	"context"
	"database/sql"
	"fmt"

	"sk_mall/rpc/rpc_merchant/internal/svc"
	"sk_mall/rpc/rpc_merchant/types"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logc"
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

func (l *IncreaseBalanceLogic) IncreaseBalance(in *__.IncreaseBalanceReq) (*__.IncreaseBalanceResp, error) {
	fmt.Println(in)
	barrier, _ := dtmgrpc.BarrierFromGrpc(l.ctx)
	db, _ := l.svcCtx.DBConn.RawDB()
	err := barrier.CallWithDB(db, func(tx *sql.Tx) error {

		_, e := tx.Exec(`update merchant_accounts set balance = balance + ? where merchant_id = ?`, in.TotalAmount, in.MerchantId)
		if e != nil {
			return e
		}
		_, e = tx.Exec(`insert into pay_logs(id,price) values(?,?)`, in.OrderNo, in.TotalAmount)
		return e
	})
	if err != nil {
		logc.Errorf(l.ctx, "[DTM] Increase Balance err:%s ,orderNo:%d", err.Error(), in.OrderNo)
		return &__.IncreaseBalanceResp{}, err
	}
	return &__.IncreaseBalanceResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
