package login

import (
	"context"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	"sk_mall/rpc/rpc_user/user"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	gresp, e := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		UserName: req.Username,
		Pwd:      req.Password,
	})
	if e != nil {
		logc.Errorf(l.ctx, "[UserRpc] Login err:%s", e.Error())
		resp = &types.LoginResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	return &types.LoginResp{
		Code:     int(gresp.Base.Code),
		Msg:      gresp.Base.Msg,
		NickName: gresp.NickName,
		Token:    gresp.Token,
		Avatar:   gresp.Avatar,
	}, nil
}
