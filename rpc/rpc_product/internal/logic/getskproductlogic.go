package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/go-redsync/redsync/v4"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetSKProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSKProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSKProductLogic {
	return &GetSKProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type SKInfo struct {
	Id           int64          `db:"id"`
	EvId         int64          `db:"event_id"`
	ProductId    int64          `db:"product_id"`
	SeckillPrice sql.NullString `db:"seckill_price"`
	SeckillStock int64          `db:"seckill_stock"`
	MerchantId   int64          `db:"merchant_id"`
}

// rpc GetProduct(GetProductReq) returns(GetProductResp);
func (l *GetSKProductLogic) GetSKProduct(in *__.GetSKProductReq) (*__.GetSKProductResp, error) {
	var skInfo SKInfo
	proKey := fmt.Sprintf("sk:product:%d", in.SKProductID)
	for {
		proMap, e1 := l.svcCtx.Rds.HgetallCtx(l.ctx, proKey)
		if e1 != nil {
			logc.Errorf(l.ctx, "[Redis] query err:%s", e1.Error())
			return &__.GetSKProductResp{}, e1
		}
		if len(proMap) == 0 {
			//获取分布式锁
			lockCtx, cancel := context.WithCancel(l.ctx)
			var mutex *redsync.Mutex
			var lockErr error
			lockName := fmt.Sprintf("lock:skproduct:%d", in.SKProductID)
			mutex, lockErr = l.svcCtx.RLCreater.GetLock(lockCtx, lockName, redsync.WithTries(3))
			if lockErr != nil {
				continue
			}
			//拿到锁，直接更新缓存
			sql := `SELECT 
   					sp.id,
    				sp.event_id,
    				sp.product_id,
    				sp.seckill_price,
    				sp.seckill_stock,
   					(SELECT merchant_id FROM seckill_events WHERE id = sp.event_id) as merchant_id
					FROM seckill_products sp
					WHERE sp.id = ?`
			err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &skInfo, sql, in.SKProductID)
			if err != nil {
				if errors.Is(err, sqlx.ErrNotFound) {
					l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)
					return &__.GetSKProductResp{
						Base: &__.BaseResp{
							Code: 1000,
							Msg:  "商品不存在",
						},
					}, nil
				}
				l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)
				logc.Errorf(l.ctx, "[DBConn] query err:%s", err.Error())
				return &__.GetSKProductResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "内部错误",
					},
				}, err
			}
			t_evid := strconv.Itoa(int(skInfo.EvId))
			t_productid := strconv.Itoa(int(skInfo.ProductId))
			t_merchantid := strconv.Itoa(int(skInfo.MerchantId))
			t_stock := strconv.Itoa(int(skInfo.SeckillStock))

			fields := map[string]string{
				"event_id":      t_evid,
				"product_id":    t_productid,
				"seckill_price": skInfo.SeckillPrice.String,
				"merchant_id":   t_merchantid,
				"stock":         t_stock,
			}
			e2 := l.svcCtx.Rds.HmsetCtx(l.ctx, proKey, fields)
			if e2 != nil {
				l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)
				logc.Errorf(l.ctx, "[Redis] set err:%s", e2.Error())
				return &__.GetSKProductResp{}, e2
			}
			e2 = l.svcCtx.Rds.ExpireCtx(l.ctx, proKey, 1200) //20min
			if e2 != nil {
				l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)
				logc.Errorf(l.ctx, "[Redis] set err:%s", e2.Error())
				return &__.GetSKProductResp{}, e2
			}
			//解锁
			l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)
			break
		} else {
			t_evid, _ := strconv.Atoi(proMap["event_id"])
			t_proid, _ := strconv.Atoi(proMap["product_id"])
			t_merchantid, _ := strconv.Atoi(proMap["merchant_id"])
			t_skprice := proMap["seckill_price"]
			t_stock, _ := strconv.Atoi(proMap["stock"])
			skInfo.EvId = int64(t_evid)
			skInfo.ProductId = int64(t_proid)
			skInfo.SeckillPrice.String = t_skprice
			skInfo.MerchantId = int64(t_merchantid)
			skInfo.SeckillStock = int64(t_stock)
			skInfo.Id = in.SKProductID

			break
		}
	}
	// 查询实时库存，如果不存在说明活动未开始，直接使用sk:product:%d的stock，
	// 因为活动未开始，所以不会数据不一致
	stockKey := fmt.Sprintf("seckill:stock:%d", in.SKProductID)
	stockStr, e3 := l.svcCtx.Rds.GetCtx(l.ctx, stockKey)
	if e3 != nil {
		logc.Errorf(l.ctx, "[Redis] get err:%s", e3.Error())
		return &__.GetSKProductResp{}, e3
	}
	if stockStr != "" {
		t_stock, _ := strconv.Atoi(stockStr)
		skInfo.SeckillStock = int64(t_stock)

	}
	return &__.GetSKProductResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Info: &__.SKProductInfo{
			Id:           skInfo.Id,
			EvId:         skInfo.EvId,
			ProductId:    skInfo.ProductId,
			SeckillPrice: skInfo.SeckillPrice.String,
			SeckillStock: skInfo.SeckillStock,
			MerchantId:   skInfo.MerchantId,
		},
	}, nil
}
