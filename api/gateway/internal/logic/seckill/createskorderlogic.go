package seckill

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_seckill/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSkOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSkOrderLogic {
	return &CreateSkOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSkOrderLogic) CreateSkOrder(req *types.CreateSkOrderReq) (resp *types.CreateSkOrderResp, err error) {
	userid, _ := strconv.Atoi(req.UserId)
	skpid, _ := strconv.Atoi(req.SKProductId)
	addrid, _ := strconv.Atoi(req.AddrId)
	gresp, e1 := l.svcCtx.SecKillRpc.CreateOrder(l.ctx, &__.CreateOrderReq{
		UserId:      uint64(userid),
		SkProductId: uint64(skpid),
		AddrId:      uint64(addrid),
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[SecKillRpc] CreateOrder err:%s", e1.Error())
		resp = &types.CreateSkOrderResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}

	resp = &types.CreateSkOrderResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}

	return
}
