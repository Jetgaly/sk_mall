package user

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	"sk_mall/rpc/rpc_user/user"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddrsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAddrsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddrsLogic {
	return &GetAddrsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddrsLogic) GetAddrs(req *types.GetAddrsReq) (resp *types.GetAddrsResp, err error) {
	// todo: add your logic here and delete this line
	id, _ := strconv.Atoi(req.Id)
	gresp, e1 := l.svcCtx.UserRpc.GetUserAddrsById(l.ctx, &user.GetAddrsByIdReq{
		Id: int64(id),
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[UserRpc] GetUserAddrsById err:%s", e1.Error())
		resp = &types.GetAddrsResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	var list []types.AddressInfo
	for _, a := range gresp.Data {
		addrId := strconv.Itoa(int(a.AddrId))
		list = append(list, types.AddressInfo{
			Name:     a.ReceiverName,
			Phone:    a.ReceiverPhone,
			Province: a.Province,
			City:     a.City,
			District: a.District,
			Detail:   a.DetailAddress,
			Default:  a.IsDefault,
			Id:       addrId,
		})
	}

	resp = &types.GetAddrsResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
		List: list,
	}
	err = nil
	return
}
