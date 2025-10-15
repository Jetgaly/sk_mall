package product

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_product/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSkProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSkProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSkProductListLogic {
	return &GetSkProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSkProductListLogic) GetSkProductList(req *types.GetSkProductListReq) (resp *types.GetSkProductListResp, err error) {
	gresp, e1 := l.svcCtx.ProductRpc.GetSkproductList(l.ctx, &__.GetSkproductListReq{
		Key:   req.Key,
		Limit: req.Limit,
		Page:  req.Page,
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[ProductRpc] GetSkproductList err:%s", e1.Error())
		resp = &types.GetSkProductListResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	var list []types.SkproductElem
	for _, v := range gresp.List {
		list = append(list, types.SkproductElem{
			EvId:  strconv.Itoa(int(v.EvId)),
			ProId: strconv.Itoa(int(v.ProId)),
			Stock: v.Stock,
			STime: v.STime,
			ETime: v.ETime,
			Name:  v.Name,
			Desc:  v.Desc,
			Price: v.Price,
			Cover: v.Cover,
			Id:    strconv.Itoa(int(v.Id)),
		})
	}
	resp = &types.GetSkProductListResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
		List: list,
	}
	return
}
