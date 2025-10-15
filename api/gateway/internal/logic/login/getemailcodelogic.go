package login

import (
	"context"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	"sk_mall/rpc/rpc_user/user"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailCodeLogic {
	return &GetEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmailCodeLogic) GetEmailCode(req *types.EmailCodeReq) (resp *types.EmailCodeResp, err error) {
	// todo: 校验邮箱格式
	gresp, e1 := l.svcCtx.UserRpc.GetEmailCode(l.ctx, &user.GetEmailCodeReq{
		Email: req.Email,
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[UserRpc] GetEmailCode err:%s", e1.Error())
		resp = &types.EmailCodeResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	resp = &types.EmailCodeResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	err = nil
	return
}
