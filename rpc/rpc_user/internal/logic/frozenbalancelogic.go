package logic

import (
	"context"
	"errors"
	"sk_mall/rpc/rpc_product/product"
	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FrozenBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFrozenBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FrozenBalanceLogic {
	return &FrozenBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var (
	ErrBalanceNotEnough = errors.New("balanceIsNotEnough")
)

func (l *FrozenBalanceLogic) FrozenBalance(in *__.FrozenBalanceReq) (*__.FrozenBalanceResq, error) {

	resp, err := l.svcCtx.ProductRpc.GetSKProduct(l.ctx, &product.GetSKProductReq{SKProductID: in.SkProductId})
	if err != nil {
		logc.Errorf(l.ctx, "[ProductRpc] err:%s", err.Error())
		return &__.FrozenBalanceResq{}, err
	}

	if resp.Base.Code == 1000 {
		return &__.FrozenBalanceResq{
			Base: &__.BaseResp{
				Code: 1000,
				Msg:  "商品不存在",
			},
		}, nil
	}
	if resp.Base.Code != 0 {
		return &__.FrozenBalanceResq{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "未知错误",
			},
		}, err
	}

	err = l.svcCtx.DBConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		r, e1 := session.ExecCtx(ctx, "update user_wallets set frozen_balance = frozen_balance+?,balance = balance-? where user_id=? and balance >= ?", resp.Info.SeckillPrice, resp.Info.SeckillPrice, in.UserId, resp.Info.SeckillPrice)
		if e1 != nil {
			return e1
		}
		ra, e2 := r.RowsAffected()
		if e2 != nil {
			return e2
		}
		if ra == 0 {
			return ErrBalanceNotEnough
		}
		_, e3 := session.ExecCtx(ctx, "insert into frozen_log(id,price,status) values(?,?,?)", in.OrderNo, resp.Info.SeckillPrice, 0)
		if e3 != nil {
			return e3
		}
		return nil
	})
	if errors.Is(err, ErrBalanceNotEnough) {
		return &__.FrozenBalanceResq{
			Base: &__.BaseResp{
				Code: 1001,
				Msg:  "余额不足",
			},
		}, nil
	} else if err != nil {
		//其他错误
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return &__.FrozenBalanceResq{
					Base: &__.BaseResp{
						Code: 1011,
						Msg:  "订单已存在",
					},
				}, nil
			}
		}
		logc.Errorf(l.ctx, "[Tranc] err:%s", err.Error())
		return &__.FrozenBalanceResq{}, err
	}
	return &__.FrozenBalanceResq{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		TotalAmount: resp.Info.SeckillPrice,
		MerchantId:  resp.Info.MerchantId,
	}, nil
}
