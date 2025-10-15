package product

import (
	"context"
	"encoding/json"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_merchant/types"
	p "sk_mall/rpc/rpc_product/product"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type product struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	Desc   string `json:"desc"`
	Status int    `json:"status"`
	Stock  int    `json:"stock"`
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductReq) (resp *types.CreateProductResp, err error) {
	var model product
	err = nil
	e1 := json.Unmarshal([]byte(req.JsonInfo), &model)
	if e1 != nil {
		resp = &types.CreateProductResp{
			Code: 997,
			Msg:  e1.Error(),
		}
		return
	}
	var ok bool
	model.Desc, ok = utils.CheckString(model.Desc, 10, 100)
	if !ok {
		resp = &types.CreateProductResp{
			Code: 21,
			Msg:  "desc为10-100个字符",
		}
		return
	}
	model.Name, ok = utils.CheckString(model.Name, 1, 10)
	if !ok {
		resp = &types.CreateProductResp{
			Code: 22,
			Msg:  "name为1-10个字符",
		}
		return
	}
	if model.Status != 1 && model.Status != 2 {
		resp = &types.CreateProductResp{
			Code: 23,
			Msg:  "status必须为1或0",
		}
		return
	}
	if model.Stock < 0 {
		resp = &types.CreateProductResp{
			Code: 24,
			Msg:  "stock非法",
		}
		return
	}
	model.Price, ok = utils.IsPositiveNumber(model.Price)
	if !ok {
		resp = &types.CreateProductResp{
			Code: 22,
			Msg:  "price非法",
		}
		return
	}
	//merchantid
	uId, _ := strconv.Atoi(req.UserId)
	gmresp, e2 := l.svcCtx.MerchantRpc.GetMerchantId(l.ctx, &__.GetMerchantIdReq{
		UserId: uint64(uId),
	})
	if e2 != nil {
		logc.Errorf(l.ctx, "[MerchantRpc] GetMerchantId err:%s", e2.Error())
		resp = &types.CreateProductResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	if gmresp.Base.Code == 1 {
		resp = &types.CreateProductResp{
			Code: 20,
			Msg:  "账号非商铺账号",
		}
		return
	}
	gresp, e3 := l.svcCtx.ProductRpc.CreateProduct(l.ctx, &p.CreateProductReq{
		MerchantId:       gmresp.MId,
		Name:             model.Name,
		Desc:             model.Desc,
		Price:            model.Price,
		Stock:            uint64(model.Stock),
		Status:           int32(model.Status),
		CoverImageSuffix: req.CoverSuffix,
		CoverImage:       req.CoverImg,
	})
	if e3 != nil {
		logc.Errorf(l.ctx, "[ProductRpc] CreateProduct err:%s", e3.Error())
		resp = &types.CreateProductResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	resp = &types.CreateProductResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	return
}
