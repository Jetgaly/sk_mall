package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAddrByAddrIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAddrByAddrIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddrByAddrIdLogic {
	return &DeleteAddrByAddrIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAddrByAddrIdLogic) DeleteAddrByAddrId(in *__.DeleteAddrByAddrIdReq) (*__.DeleteAddrByAddrIdResp, error) {
	// todo: add your logic here and delete this line
	var resp __.DeleteAddrByAddrIdResp
	_, err := l.svcCtx.DBConn.Exec("delete from user_addresses where id = ?", in.Id)

	//err := l.svcCtx.DBConn.QueryRow(&addr_flag, "select user_id from user_addresses where id = ?", in.Id)
	if err == nil {
		resp.Base = &__.BaseResp{
			Code: 0,
			Msg:  "success",
		}
	} else {
		logc.Errorf(l.ctx, "[DBConn] delete err:%s", err.Error())
	}
	return &resp, err
}
