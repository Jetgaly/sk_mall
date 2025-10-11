package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type JwtAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJwtAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JwtAuthLogic {
	return &JwtAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JwtAuthLogic) JwtAuth(in *__.JwtAuthReq) (*__.JwtAuthResp, error) {
	// todo: add your logic here and delete this line
	resp, err := utils.ParseToken(in.Token)
	if err != nil {
		return &__.JwtAuthResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "auth failed",
			},
		}, nil
	}
	return &__.JwtAuthResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		UserId: int64(resp.UserId),
	}, nil
}
