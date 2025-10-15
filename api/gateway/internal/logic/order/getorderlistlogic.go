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

type GetOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderListLogic) GetOrderList(req *types.GetOrderListReq) (resp *types.GetOrderListResp, err error) {
	userid, _ := strconv.Atoi(req.UserId)
	gresp, e1 := l.svcCtx.AggrOrderRpc.GetOrderList(l.ctx, &__.GetOrderListReq{
		UserId: int64(userid),
		Page:   req.Page,
		Limit:  req.Limit,
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[AggrOrderRpc] GetOrderList err:%s", e1.Error())
		resp = &types.GetOrderListResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	if gresp.Base.Code != 0 {
		resp = &types.GetOrderListResp{
			Code: int(gresp.Base.Code),
			Msg:  gresp.Base.Msg,
		}
		return
	}
	var List []types.OrderElem
	for _, v := range gresp.List {
		List = append(List, types.OrderElem{
			OrderNo:     strconv.Itoa(int(v.OrderNo)),
			Product:     types.OrderProduct{SKProductId: strconv.Itoa(int(v.SkProduck.SkProductId))},
			Quantity:    v.UnitPrice,
			UnitPrice:   v.UnitPrice,
			TotalAmount: v.TotalAmount,
			Status:      int64(v.Status),
			CreateAt:    v.CreatedAt.String(),
			ExpireAt:    v.ExpireAt.String(),
		})
	}
	resp = &types.GetOrderListResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
		List: List,
	}
	return
}
