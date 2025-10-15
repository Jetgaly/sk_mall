package user

import (
	"context"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_user/types"
	"sk_mall/rpc/rpc_user/user"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAddrLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAddrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAddrLogic {
	return &CreateAddrLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAddrLogic) CreateAddr(req *types.CreateAddrReq) (resp *types.CreateAddrResp, err error) {
	//todo:参数校验
	userId, _ := strconv.Atoi(req.UserId)
	gresp, e1 := l.svcCtx.UserRpc.CreateUserAddrsById(l.ctx, &user.CreateAddrsByIdReq{
		Id: int64(userId),
		Addr: &__.Address{
			ReceiverName:  req.Addr.Name,
			ReceiverPhone: req.Addr.Phone,
			Province:      req.Addr.Province,
			City:          req.Addr.City,
			District:      req.Addr.District,
			DetailAddress: req.Addr.Detail,
			IsDefault:     req.Addr.Default,
		},
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[UserRpc] CreateUserAddrsById err:%s", e1.Error())
		resp = &types.CreateAddrResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	resp = &types.CreateAddrResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	err = nil
	return
}
