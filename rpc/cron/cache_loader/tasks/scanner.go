package tasks

import (
	"context"
	"errors"
	"fmt"
	"sk_mall/rpc/cron/cache_loader/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Scanner struct {
	SvcCtx *svc.ServiceContext
}

type SecEv struct {
	Id     int64     `db:"id"`
	Status int8      `db:"status"`
	STime  time.Time `db:"start_time"`
	ETime  time.Time `db:"end_time"`
}

var scanLuaStr string = `redis.call('hmset', KEYS[1], 'id', ARGV[1], 'status', ARGV[2], 'start_time', ARGV[3], 'end_time', ARGV[4])
						 redis.call('expireat', KEYS[1], ARGV[5])
						 return 1`

//  `
// -- 检查是否已存在
// if redis.call('exists', KEYS[1]) == 0 then
//     -- 创建活动数据
//     redis.call('hmset', KEYS[1], 'id', ARGV[1], 'status', ARGV[2], 'start_time', ARGV[3], 'end_time', ARGV[4])
//     -- 设置过期时间
//     redis.call('expireat', KEYS[1], ARGV[5])
//     return 1
// else
//     return 0
// end
// `

func (s *Scanner) Run() {
	//凌晨一点扫描
	var EvList []SecEv
	now := time.Now()
	startRange := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endRange := startRange.Add(28 * time.Hour) // 今天 + 明天额外4小时缓冲
	query := `SELECT id, status, start_time, end_time 
			  FROM seckill_events 
			  WHERE start_time >= ? AND start_time < ? 
			  AND status = 1`
	startTimeStr := startRange.Format("2006-01-02 15:04:05")
	endTimeStr := endRange.Format("2006-01-02 15:04:05")
	err := s.SvcCtx.DBConn.QueryRows(&EvList, query, startTimeStr, endTimeStr)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logc.Errorf(context.Background(), "[DBConn] query err:%s", err.Error())
		return
	}

	logc.Info(context.Background(), "[Scanner]:scanning")

	for _, ev := range EvList {
		key := fmt.Sprintf("seckill:event:%d", ev.Id)
		// 计算过期时间（活动结束后2小时）
		expireAt := ev.ETime.Add(2 * time.Hour).Unix()

		// 执行Lua脚本
		result, e1 := s.SvcCtx.Rds.Eval(scanLuaStr, []string{key},
			ev.Id, ev.Status, ev.STime.Unix(), ev.ETime.Unix(), expireAt)

		if e1 != nil {
			logc.Errorf(context.Background(), "lua script execution failed: %s", e1)
			continue
		}

		// 检查执行结果
		if resultInt, ok := result.(int64); ok {
			if resultInt == 1 {
				logc.Infof(context.Background(), "[Scanner] created cache for event %d", ev.Id)
				continue
			}
		}

	}
}
