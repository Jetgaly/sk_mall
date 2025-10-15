package product

import (
	"context"
	"strconv"
	"time"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_merchant/types"
	p "sk_mall/rpc/rpc_product/product"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEventLogic {
	return &CreateEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEventLogic) CreateEvent(req *types.CreateEventReq) (resp *types.CreateEventResp, err error) {
	err = nil
	st, e1 := strconv.Atoi(req.StartTime)
	et, e2 := strconv.Atoi(req.EndTime)
	if e1 != nil || e2 != nil {
		resp = &types.CreateEventResp{
			Code: 20,
			Msg:  "活动时间格式非法",
		}
		return
	}
	if st >= et {
		resp = &types.CreateEventResp{
			Code: 21,
			Msg:  "活动时间非法",
		}
		return
	}
	var ok bool
	req.Name, ok = utils.CheckString(req.Name, 1, 10)
	if !ok {
		resp = &types.CreateEventResp{
			Code: 22,
			Msg:  "name为1-10个字符",
		}
		return
	}
	uId, _ := strconv.Atoi(req.UserId)
	gmresp, e3 := l.svcCtx.MerchantRpc.GetMerchantId(l.ctx, &__.GetMerchantIdReq{
		UserId: uint64(uId),
	})
	if e3 != nil {
		logc.Errorf(l.ctx, "[MerchantRpc] GetMerchantId err:%s", e3.Error())
		resp = &types.CreateEventResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	if gmresp.Base.Code == 1 {
		resp = &types.CreateEventResp{
			Code: 23,
			Msg:  "账号非商铺账号",
		}
		return
	}

	gresp, e4 := l.svcCtx.ProductRpc.CreateSKEvent(l.ctx, &p.CreateSKEventReq{
		MerchantId: gmresp.MId,
		Name:       req.Name,
		StartTime:  timestamppb.New(time.Unix(int64(st), 0)),
		EndTime:    timestamppb.New(time.Unix(int64(et), 0)),
		Status:     1,
	})
	if e4 != nil {
		logc.Errorf(l.ctx, "[ProductRpc] CreateSKEvent err:%s", e4.Error())
		resp = &types.CreateEventResp{
			Code: 999,
			Msg:  "server err",
		}
		return
	}
	resp = &types.CreateEventResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	return
}
