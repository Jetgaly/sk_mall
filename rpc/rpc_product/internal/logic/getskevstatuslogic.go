package logic

import (
	"context"
	"math/rand/v2"

	"errors"
	"fmt"
	"strconv"
	"time"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/go-redsync/redsync/v4"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetSKEvStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSKEvStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSKEvStatusLogic {
	return &GetSKEvStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var luaStr = `if redis.call('exists', KEYS[1]) == 0 then
       redis.call('hmset', KEYS[1], 'id', ARGV[1], 'status', ARGV[2], 'start_time', ARGV[3], 'end_time', ARGV[4])
       redis.call('expire', KEYS[1], ARGV[5])
       return 1
   else
       return 0
   end`

type Ev struct {
	Id     int64     `db:"id"`
	Status int8      `db:"status"`
	STime  time.Time `db:"start_time"`
	ETime  time.Time `db:"end_time"`
}

func (l *GetSKEvStatusLogic) GetSKEvStatus(in *__.GetSKEvStatusReq) (*__.GetSKEvStatusResp, error) {
	// todo: add your logic here and delete this line
	//seckill:event:%d
	var evInfo Ev
	evKey := fmt.Sprintf("seckill:event:%d", in.EvId)
	for {
		statusStr, e1 := l.svcCtx.Rds.HgetCtx(l.ctx, evKey, "status")
		if e1 != nil && !errors.Is(e1, redis.Nil) {
			logc.Errorf(l.ctx, "[Redis] hget err:%s", e1.Error())
			return &__.GetSKEvStatusResp{}, e1
		}
		if errors.Is(e1, redis.Nil) {
			//获取分布式锁，更新缓存
			lockName := fmt.Sprintf("lock:skev:%d", in.EvId)
			lockCtx, cancel := context.WithCancel(l.ctx)
			var mutex *redsync.Mutex
			var lockErr error
			mutex, lockErr = l.svcCtx.RLCreater.GetLock(lockCtx, lockName, redsync.WithTries(3))
			if lockErr != nil {
				continue
			}
			//释放锁
			defer l.svcCtx.RLCreater.ReleaseLock(mutex, cancel)
			//添加缓存
			sql := `SELECT id, status, start_time, end_time 
			  FROM seckill_events 
			  WHERE id = ?`
			e2 := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &evInfo, sql, in.EvId)
			if errors.Is(e2, sqlx.ErrNotFound) {
				return &__.GetSKEvStatusResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "商品不存在",
					},
				}, nil
			}
			if e2 != nil {
				logc.Errorf(l.ctx, "[DBConn] query err:%s", e2.Error())
				return &__.GetSKEvStatusResp{}, e2
			}
			expire := 1200 + rand.IntN(10) //20min +- 10min
			_, e3 := l.svcCtx.Rds.EvalCtx(l.ctx, luaStr, []string{evKey}, evInfo.Id, evInfo.Status, evInfo.STime.Unix(), evInfo.ETime.Unix(), expire)
			if e3 != nil {
				logc.Errorf(l.ctx, "[Redis] eval lua err:%s", e3.Error())
				return &__.GetSKEvStatusResp{}, e3
			}
			break
		} else {
			//读到缓存
			t_status, _ := strconv.Atoi(statusStr)
			evInfo.Status = int8(t_status)
			break
		}
	}
	return &__.GetSKEvStatusResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Status: int32(evInfo.Status),
	}, nil
}
