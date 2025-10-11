package logic

import (
	"context"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
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
				return &__.CreateSKProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "未知错误",
					},
				}, err
			}
		} else {
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
