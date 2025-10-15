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

type DeleteAddrLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAddrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddrLogic {
	return &DeleteAddrLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAddrLogic) DeleteAddr(req *types.DeleteAddrReq) (resp *types.DeleteAddrResp, err error) {

	id, _ := strconv.Atoi(req.Id)
	gresp, e1 := l.svcCtx.UserRpc.DeleteAddrByAddrId(l.ctx, &user.DeleteAddrByAddrIdReq{
		Id: int64(id),
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[UserRpc] DeleteAddrByAddrId err:%s", e1.Error())
		resp = &types.DeleteAddrResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	resp = &types.DeleteAddrResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	err = nil
	return
}
