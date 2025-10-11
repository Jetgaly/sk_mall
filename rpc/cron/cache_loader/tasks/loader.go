package tasks

import (
	"context"
	"errors"
	
	"sk_mall/rpc/cron/cache_loader/internal/svc"
	"strconv"
	"strings"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Loader struct {
	SvcCtx *svc.ServiceContext
}

var updateLuaStr string = `
local id = redis.call('HGET', KEYS[1], 'id')
local status = redis.call('HGET', KEYS[1], 'status')
local start_time = redis.call('HGET', KEYS[1], 'start_time')
local end_time = redis.call('HGET', KEYS[1], 'end_time')

if not id or not status or not start_time or not end_time then
    return 0
end

local current_time = tonumber(ARGV[1])
local start_time_num = tonumber(start_time)
local end_time_num = tonumber(end_time)
local status_num = tonumber(status)

if not start_time_num or not end_time_num or not status_num then
    return 0
end

if current_time > end_time_num and status_num ~= 3 then
    redis.call('HSET', KEYS[1], 'status', '3')
    return 3
elseif current_time >= start_time_num and current_time <= end_time_num and status_num ~= 2 then
    redis.call('HSET', KEYS[1], 'status', '2')
    return 2
elseif current_time < start_time_num and status_num ~= 1 then
    redis.call('HSET', KEYS[1], 'status', '1')
    return 1
end
return 4
`
var preloadCheckLuaStr string = `
local start_time = redis.call('HGET', KEYS[1], 'start_time')
local end_time = redis.call('HGET', KEYS[1], 'end_time')
local status = redis.call('HGET', KEYS[1], 'status')
if not start_time or not end_time or not status then
    return 0
end

local current_time = tonumber(ARGV[1])
local start_time_num = tonumber(start_time)
local end_time_num = tonumber(end_time)
local status_num = tonumber(status)

if not start_time_num or not end_time_num or not status_num then
    return 0
end

if current_time >= start_time_num and current_time <= end_time_num and status_num == 1 then
    return end_time_num-current_time
end
return 0

`

type SecProductStock struct {
	Id    int64 `db:"id"`
	Stock int64 `db:"seckill_stock"`
}

// 扫描redis的event
//
//seckill:event:id
func (l *Loader) Run() {
	logc.Info(context.Background(), "[Loader]:loading")
	ctx := context.Background()
	var cursor uint64
	count := int64(100) //扫描槽位的数量
	pattern := "seckill:event:*"
	lockCtx, cancel := context.WithCancel(ctx)
	var mutex *redsync.Mutex
	var lockErr error
	for {
		mutex, lockErr = l.SvcCtx.RLCreater.GetLock(lockCtx, "lock:seckill:loader",redsync.WithTries(32))
		if lockErr == nil {
			defer l.SvcCtx.RLCreater.ReleaseLock(mutex, cancel)
			break
		}
	}
	for {
		// 扫描Hash键
		keys, nextCursor, err := l.SvcCtx.Rds.Scan(cursor, pattern, count)
		if err != nil {
			logc.Errorf(ctx, "[Redis] scan err:%s", err.Error())
			break
		}

		// 处理批量Hash键
		if len(keys) > 0 {
			//preload
			l.preload(ctx, keys)
			//status
			l.processBatchHashes(ctx, keys)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

func (l *Loader) preload(c context.Context, keys []string) {
	keyPrefix := "seckill:stock:"
	//提前5分钟
	preTime := 5 * time.Minute
	now := time.Now().Add(preTime).Unix()
	for _, k := range keys {
		result, e1 := l.SvcCtx.Rds.Eval(preloadCheckLuaStr, []string{k}, now)
		if e1 != nil {
			logc.Errorf(c, "[Redis] lua exec err:%s", e1.Error())
			continue
		}
		res, ok := result.(int64)
	
		if ok && res > 0 {
			
			EvIdStr := strings.Split(k, ":")[2]
			EvId, _ := strconv.Atoi(EvIdStr)
			//todo :分批优化
			var StockList []SecProductStock
			e2 := l.SvcCtx.DBConn.QueryRows(&StockList, "select id,seckill_stock from seckill_products where event_id = ?", EvId)
			if errors.Is(e2, sqlx.ErrNotFound) {
				continue
			} else if e2 != nil {
				logc.Errorf(c, "[DBConn] query err:%s", e2.Error())
				continue
			}
	
			err := l.SvcCtx.Rds.Pipelined(func(p redis.Pipeliner) error {
				var perr error
				for _, s := range StockList {
					skey := keyPrefix + strconv.Itoa(int(s.Id))
					sval := strconv.Itoa(int(s.Stock))

					exp := int64(preTime.Seconds()) + res

					ret := p.SetNX(c, skey, sval, time.Duration(exp)*time.Second)
					perr = ret.Err()
					if perr != nil {
						logc.Errorf(c, "[Redis] err:%s", skey)
						continue
					}
				}
				return perr
			})
			if err != nil {
				logc.Errorf(c, "[Redis] pipeline err:%s", err.Error())
				continue
			}
		}
	}
}
func (l *Loader) processBatchHashes(c context.Context, keys []string) {
	now := time.Now().Add(5 * time.Second).Unix()

	for _, k := range keys {
		result, e1 := l.SvcCtx.Rds.Eval(updateLuaStr, []string{k}, now)
		if e1 != nil {
			logc.Errorf(c, "[Redis] lua exec err:%s", e1.Error())
			continue
		}

		res, ok := result.(int64)
		if ok && res != 0 && res != 4 {
			EvIdStr := strings.Split(k, ":")[2]
			EvId, _ := strconv.Atoi(EvIdStr)
			//修改数据库
			_, e2 := l.SvcCtx.DBConn.Exec("update seckill_events set status = ? where id = ?", res, EvId)
			if e2 != nil {
				logc.Errorf(c, "[DBConn] err:%s", e2.Error())
			}
		}
	}
}
