package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeUserNicknameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeUserNicknameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserNicknameLogic {
	return &ChangeUserNicknameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeUserNicknameLogic) ChangeUserNickname(in *__.ChangeUserNicknameReq) (*__.ChangeUserNicknameResp, error) {
	// todo: add your logic here and delete this line
	var resp __.ChangeUserNicknameResp
	//只有登录的用户才能调用，调用这个接口前一定已经确认了数据库有这个用户
	_, err := l.svcCtx.DBConn.Exec("update users set nick_name=? where id=?", in.NewName, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "[DBConn] update err:%s", err.Error())
		resp = __.ChangeUserNicknameResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "内部错误",
			},
		}
	}
	resp = __.ChangeUserNicknameResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}
	return &resp, err
}
