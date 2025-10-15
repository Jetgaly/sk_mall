package product

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_merchant/types"
	p "sk_mall/rpc/rpc_product/product"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSkProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSkProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSkProductLogic {
	return &CreateSkProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSkProductLogic) CreateSkProduct(req *types.CreateSkProductReq) (resp *types.CreateSkProductResp, err error) {
	//参数验证
	var ok bool
	req.EvId, ok = utils.IsPositiveNumber(req.EvId)
	evid, e3 := strconv.Atoi(req.EvId)
	if !ok || e3 != nil {
		resp = &types.CreateSkProductResp{
			Code: 20,
			Msg:  "EvId非法",
		}
		return
	}
	req.Price, ok = utils.IsPositiveNumber(req.Price)
	if !ok {
		resp = &types.CreateSkProductResp{
			Code: 21,
			Msg:  "Price非法",
		}
		return
	}
	req.Stock, ok = utils.IsPositiveNumber(req.Stock)
	stock, e1 := strconv.Atoi(req.Stock)
	if !ok || e1 != nil {
		resp = &types.CreateSkProductResp{
			Code: 22,
			Msg:  "Stock非法",
		}
		return
	}
	req.ProductId, ok = utils.IsPositiveNumber(req.ProductId)
	pid, e2 := strconv.Atoi(req.ProductId)
	if !ok || e2 != nil {
		resp = &types.CreateSkProductResp{
			Code: 23,
			Msg:  "ProductId非法",
		}
		return
	}

	userid, _ := strconv.Atoi(req.UserId)
	gmresp, e4 := l.svcCtx.MerchantRpc.GetMerchantId(l.ctx, &__.GetMerchantIdReq{
		UserId: uint64(userid),
	})
	if e4 != nil {
		logc.Errorf(l.ctx, "[MerchantRpc] GetMerchantId err:%s", e4.Error())
		resp = &types.CreateSkProductResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	if gmresp.Base.Code != 0 {
		resp = &types.CreateSkProductResp{
			Code: 24,
			Msg:  "账号为非商铺账号",
		}
		return
	}
	gresp, e5 := l.svcCtx.ProductRpc.CreateSKProduct(l.ctx, &p.CreateSKProductReq{
		EventId:      uint64(evid),
		ProductId:    uint64(pid),
		SeckillPrice: req.Price,
		SeckillStock: uint64(stock),
		MerchantId:   int64(gmresp.MId),
	})
	if e5 != nil {
		logc.Errorf(l.ctx, "[ProductRpc] CreateSKProduct err:%s", e5.Error())
		resp = &types.CreateSkProductResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	resp = &types.CreateSkProductResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	return
}
