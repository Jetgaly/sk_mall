package order

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/aggr_order/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderInfoLogic {
	return &GetOrderInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderInfoLogic) GetOrderInfo(req *types.GetOrderInfoReq) (resp *types.GetOrderInfoResp, err error) {
	userid, _ := strconv.Atoi(req.UserId)
	gresp, e1 := l.svcCtx.AggrOrderRpc.GetOrderInfo(l.ctx, &__.GetOrderInfoReq{
		UserId:  int64(userid),
		OrderNo: req.OrderNo,
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[AggrOrderRpc] GetOrderInfo err:%s", e1.Error())
		resp = &types.GetOrderInfoResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	if gresp.Base.Code != 0 {
		resp = &types.GetOrderInfoResp{
			Code: int(gresp.Base.Code),
			Msg:  gresp.Base.Msg,
		}
		return
	}
	addrid := int(gresp.Info.Addr.AddrId)
	skpid := int(gresp.Info.SkProduck.SkProductId)
	quan := int(gresp.Info.Quantity)
	resp = &types.GetOrderInfoResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
		Info: types.OrderInfo{
			OrderNo: req.OrderNo,
			Addr: types.OrderAddr{
				AddrId: strconv.Itoa(addrid),
			},
			Product: types.OrderProduct{
				SKProductId: strconv.Itoa(skpid),
			},
			Quantity:    strconv.Itoa(quan),
			UnitPrice:   gresp.Info.UnitPrice,
			TotalAmount: gresp.Info.TotalAmount,
			Status:      int64(gresp.Info.Status),
			CreateAt:    gresp.Info.CreatedAt.String(),
			ExpireAt:    gresp.Info.ExpireAt.String(),
		},
	}
	return
}
