package payment

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_payment/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderLogic) PayOrder(req *types.PayOrderReq) (resp *types.PayOrderResp, err error) {
	userid, _ := strconv.Atoi(req.UserId)
	orderno, _ := strconv.Atoi(req.OrderNo)
	gresp, e1 := l.svcCtx.PaymentRpc.PayOrder(l.ctx, &__.PayOrderReq{
		UserId:  uint64(userid),
		OrderNo: int64(orderno),
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[PaymentRpc] PayOrder err:%s", e1.Error())
		resp = &types.PayOrderResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	resp = &types.PayOrderResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	return
}
