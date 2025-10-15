package logic

import (
	"context"
	"strconv"

	"sk_mall/rpc/rpc_product/internal/svc"
	"sk_mall/rpc/rpc_product/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSkproductListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSkproductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSkproductListLogic {
	return &GetSkproductListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSkproductListLogic) GetSkproductList(in *__.GetSkproductListReq) (*__.GetSkproductListResp, error) {
	docs, err := utils.SearchDocuments(l.ctx, l.svcCtx.ESCli, "sk_products", in.Key, int(in.Page), int(in.Limit))
	if err != nil {
		logc.Errorf(l.ctx, "[ES] search err:%s", err.Error())
		return nil, err
	}
	var list []*__.SkproductElem
	for _, v := range docs.Hits.Hits {
		skpid, _ := strconv.Atoi(v.Id)
		evid, _ := strconv.Atoi(v.Source.EvID)
		pid, _ := strconv.Atoi(v.Source.ProductID)
		price := strconv.FormatFloat(v.Source.ProductPrice, 'f', 2, 64)
		list = append(list, &__.SkproductElem{
			EvId:  int64(evid),
			ProId: int64(pid),
			Stock: int64(v.Source.Stock),
			STime: v.Source.StartTime.String(),
			ETime: v.Source.EndTime.String(),
			Name:  v.Source.ProductName,
			Desc:  v.Source.ProductDesc,
			Price: price,
			Cover: v.Source.CoverPath,
			Id:    int64(skpid),
		})
	}
	return &__.GetSkproductListResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		List: list,
	}, nil
}
