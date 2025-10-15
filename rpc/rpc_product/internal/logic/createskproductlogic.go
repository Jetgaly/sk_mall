package logic

import (
	"context"
	"errors"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateSKProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSKProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSKProductLogic {
	return &CreateSKProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSKProductLogic) CreateSKProduct(in *__.CreateSKProductReq) (*__.CreateSKProductResp, error) {
	// todo: add your logic here and delete this line
	priceDec, err := decimal.NewFromString(in.SeckillPrice)
	if err != nil {
		return &__.CreateSKProductResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "Price 解析错误",
			},
		}, nil
	}
	//验证merchant event
	var evModel struct {
		Name string `db:"name"`
	}
	e1 := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &evModel, "SELECT name FROM seckill_events WHERE merchant_id=? AND id=?", in.MerchantId, in.EventId)
	if errors.Is(e1, sqlx.ErrNotFound) {
		return &__.CreateSKProductResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "找不到活动",
			},
		}, nil
	}
	if e1 != nil {
		logc.Errorf(l.ctx, "[DBConn] query err:%s", e1.Error())
		return &__.CreateSKProductResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "未知错误",
			},
		}, err
	}
	_, err = l.svcCtx.DBConn.ExecCtx(l.ctx, "insert into seckill_products(event_id,product_id,seckill_price,seckill_stock) values(?,?,?,?)", in.EventId, in.ProductId, priceDec.String(), in.SeckillStock)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			// if mysqlErr.Number == 1062 {
			// 	return &__.CreateSKProductResp{
			// 		Base: &__.BaseResp{
			// 			Code: 1,
			// 			Msg:  "活动商品 %s 已存在",
			// 		},
			// 	}, nil
			// } else if mysqlErr.Number == 1452 { //外键约束
			// 	return &__.CreateSKProductResp{
			// 		Base: &__.BaseResp{
			// 			Code: 1,
			// 			Msg:  "找不到商品/活动",
			// 		},
			// 	}, nil
			// } else {
			// 	return &__.CreateSKProductResp{
			// 		Base: &__.BaseResp{
			// 			Code: 1,
			// 			Msg:  "未知错误",
			// 		},
			// 	}, err
			// }
			switch mysqlErr.Number {
			case 1062:
				return &__.CreateSKProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "活动商品 %s 已存在",
					},
				}, nil
			case 1452:
				return &__.CreateSKProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "找不到商品/活动",
					},
				}, nil
			default:
				logc.Errorf(l.ctx, "[DBConn] insert err:%s", err.Error())
				return &__.CreateSKProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "未知错误",
					},
				}, err
			}
		} else {
			logc.Errorf(l.ctx, "[DBConn] insert err:%s", err.Error())
			return &__.CreateSKProductResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "未知错误",
				},
			}, err
		}
	}
	return &__.CreateSKProductResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
