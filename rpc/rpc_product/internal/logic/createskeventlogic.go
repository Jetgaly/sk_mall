package logic

import (
	"context"
	"fmt"
	"time"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSKEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSKEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSKEventLogic {
	return &CreateSKEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var scanLuaStr string = `
-- 检查是否已存在
if redis.call('exists', KEYS[1]) == 0 then
    -- 创建活动数据
    redis.call('hmset', KEYS[1], 'id', ARGV[1], 'status', ARGV[2], 'start_time', ARGV[3], 'end_time', ARGV[4])
    -- 设置过期时间
    redis.call('expireat', KEYS[1], ARGV[5])
    return 1
else
    return 0
end
`

func (l *CreateSKEventLogic) CreateSKEvent(in *__.CreateSKEventReq) (*__.CreateSKEventResp, error) {
	// todo: add your logic here and delete this line
	st := in.StartTime.AsTime()
	et := in.EndTime.AsTime()

	r, err := l.svcCtx.DBConn.ExecCtx(l.ctx, "insert into seckill_events(merchant_id,name,start_time,end_time) values(?,?,?,?)", in.MerchantId, in.Name, st, et)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return &__.CreateSKEventResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  fmt.Sprintf("活动 %s 已存在", in.Name),
					},
				}, nil
			} else {
				logc.Errorf(l.ctx, "[DBConn] insert err:%s", err.Error())
				return &__.CreateSKEventResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "未知错误",
					},
				}, err
			}
		} else {
			logc.Errorf(l.ctx, "[DBConn] insert err:%s", err.Error())
			return &__.CreateSKEventResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "未知错误",
				},
			}, err
		}
	}

	now := time.Now()
	startRange := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endRange := startRange.Add(28 * time.Hour) // 今天 + 明天额外4小时缓冲

	if startRange.Before(st) && endRange.After(st) {
		//当天的秒杀ev放入redis中
		EvId, e1 := r.LastInsertId()
		if e1 != nil {
			logc.Errorf(l.ctx, "[DBConn] err:%s", e1.Error())
			return &__.CreateSKEventResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "未知错误",
				},
			}, err
		}
		key := fmt.Sprintf("seckill:event:%d", EvId)
		expireAt := et.Add(2 * time.Hour).Unix()

		// 执行Lua脚本
		result, e2 := l.svcCtx.Rds.Eval(scanLuaStr, []string{key},
			EvId, 1, st.Unix(), et.Unix(), expireAt)

		if e2 != nil {
			logc.Errorf(l.ctx, "lua script execution failed: %s", e2.Error())
			return &__.CreateSKEventResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "未知错误",
				},
			}, e2
		}

		// 检查执行结果
		if resultInt, ok := result.(int64); ok {
			switch resultInt {
			case 1:
				logc.Infof(l.ctx, "[Redis] created cache for event %d", EvId)

			default:
				logc.Infof(l.ctx, "[Redis] event %d cache already exists", EvId)
			}
		}

	}

	return &__.CreateSKEventResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, nil
}
