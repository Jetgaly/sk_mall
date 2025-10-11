package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailCodeLogic {
	return &GetEmailCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmailCodeLogic) GetEmailCode(in *__.GetEmailCodeReq) (*__.GetEmailCodeResp, error) {
	// todo: add your logic here and delete this line
	var resp __.GetEmailCodeResp
	resp.Base = &__.BaseResp{
		Code: 0,
		Msg:  "success",
	}
	code := utils.GenerateCode()
	key := l.svcCtx.EmailCodePrefix + in.Email
	err := l.svcCtx.Rds.Setex(key, code, 300) //5min

	if err != nil {
		logc.Errorf(l.ctx, "[Redis] set err:%s", err.Error())
		resp.Base = &__.BaseResp{
			Code: 1,
			Msg:  "其他错误",
		}
	}
	err = l.svcCtx.EmailSender.SendEmail(in.Email, code, "SkMall验证码")
	if err != nil {
		logc.Errorf(l.ctx, "[Email] send err:%s", err.Error())
		resp.Base = &__.BaseResp{
			Code: 1,
			Msg:  "其他错误",
		}
	}

	return &resp, err
}
