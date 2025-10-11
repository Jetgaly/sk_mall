package logic

import (
	"context"
	"errors"
	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetSKEvStatusByPIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSKEvStatusByPIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSKEvStatusByPIdLogic {
	return &GetSKEvStatusByPIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type EventStatus struct {
	Status int32 `db:"status"`
}

func (l *GetSKEvStatusByPIdLogic) GetSKEvStatusByPId(in *__.GetSKEvStatusByPIdReq) (*__.GetSKEvStatusByPIdResp, error) {
	
	//seckill:event:%d
	
	var es EventStatus
	err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &es, "SELECT e.status FROM seckill_events e JOIN seckill_products p ON e.id = p.event_id WHERE p.id = ?", in.PId)

	if errors.Is(err, sqlx.ErrNotFound) {
		return &__.GetSKEvStatusByPIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "活动不存在",
			},
		}, nil
	} else if err != nil {
		return nil, err
	}

	return &__.GetSKEvStatusByPIdResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		Status: es.Status,
	}, nil
}
