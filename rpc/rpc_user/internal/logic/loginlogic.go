package logic

import (
	"context"
	"errors"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type user struct {
	Id       int64  `db:"id"`
	NickName string `db:"nick_name"`
	Avatar   string `db:"avatar"`
	Pwd      string `db:"password_hash"`
}

func (l *LoginLogic) Login(in *__.LoginReq) (*__.LoginResp, error) {
	var t_user user
	sql := `
		select id,nick_name,avatar,password_hash from users where user_name = ?
	`
	err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &t_user, sql, in.UserName)

	if errors.Is(err, sqlx.ErrNotFound) {
		return &__.LoginResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "用户不存在",
			},
		}, nil
	}
	if err != nil {
		logc.Errorf(l.ctx, "[DBConn] query err:%s", err.Error())
		return &__.LoginResp{}, err
	}
	if !utils.CheckHashStr(in.Pwd, t_user.Pwd) {
		return &__.LoginResp{
			Base: &__.BaseResp{
				Code: 2,
				Msg:  "密码错误",
			},
		}, nil
	}
	//获取token
	token, e1 := utils.GenerateToken(int(t_user.Id), t_user.NickName)
	if e1 != nil {
		logc.Errorf(l.ctx, "[Jwt] GenerateToken err:%s", e1.Error())
		return nil, e1
	}

	return &__.LoginResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg: "success",
		},
		Token: token,
		NickName: t_user.NickName,
		Avatar: t_user.Avatar,
	}, nil
}
