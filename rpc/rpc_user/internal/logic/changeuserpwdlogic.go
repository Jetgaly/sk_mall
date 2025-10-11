package logic

import (
	"context"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserPwd struct {
	PwdHash string `db:"password_hash"`
}
type ChangeUserPwdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeUserPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserPwdLogic {
	return &ChangeUserPwdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeUserPwdLogic) ChangeUserPwd(in *__.ChangeUserPwdReq) (*__.ChangeUserPwdResp, error) {
	// todo: add your logic here and delete this line

	var resp __.ChangeUserPwdResp

	var pwd UserPwd
	err := l.svcCtx.DBConn.QueryRow(&pwd, "select password_hash from users where id = ?", in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "[DBConn] query err:%s", err.Error())
		resp = __.ChangeUserPwdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "内部错误",
			},
		}
		return &resp, err
	}
	if !utils.CheckHashStr(in.OldPwd, pwd.PwdHash) {
		//密码不对
		resp = __.ChangeUserPwdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "旧密码错误",
			},
		}
		return &resp, err
	}
	newPwdHash := utils.GetHashStr(in.NewPwd)
	_, err = l.svcCtx.DBConn.Exec("update users set password_hash=? where id=?", newPwdHash, in.Id)
	if err != nil {
		logc.Errorf(l.ctx, "[DBConn] query err:%s", err.Error())
		resp = __.ChangeUserPwdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "内部错误",
			},
		}
		return &resp, err
	}
	resp = __.ChangeUserPwdResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}
	//@todo：重新输入密码登录
	return &resp, err
}
