package logic

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"
	"sk_mall/utils"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type UserCheck struct {
	Id int64 `db:"id"`
}

func (l *CreateUserLogic) CreateUser(in *__.CreateUserReq) (*__.CreateUserResp, error) {
	// todo: add your logic here and delete this line

	var u UserCheck
	err := l.svcCtx.DBConn.QueryRowCtx(context.Background(), &u, "select id from users where email = ? or user_name = ? limit 1", in.Email, in.UserName)
	if err == nil {
		logc.Infof(l.ctx, "[DBConn]用户已存在:email=%s or user_name=%s", in.Email, in.UserName)
		return &__.CreateUserResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "email/user_name已存在",
			},
		}, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		//验证验证码
		key := l.svcCtx.EmailCodePrefix + in.Email
		code, re := l.svcCtx.Rds.Get(key)
		if re != nil {
			return &__.CreateUserResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "验证码不存在",
				},
			}, nil
		}
		if code != in.Code {
			return &__.CreateUserResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "验证码错误",
				},
			}, nil
		}
		//参数校验在api网关做
		err = nil
		//密码sha256加密
		hashpwd := utils.GetHashStr(in.Password)

		AvatarName := uuid.NewString() + "." + in.AvatarName //api传后缀
		avatarPath := filepath.Join(l.svcCtx.Config.Avatar.UploadPath, AvatarName)

		//如果写文件失败头像为空字符串
		if err := os.WriteFile(avatarPath, in.Avatar, 0644); err != nil {
			logc.Errorf(l.ctx, "[AvatarUpload]文件写入失败: %s", err.Error())
			avatarPath = ""
		}

		err := l.svcCtx.DBConn.Transact(
			func(session sqlx.Session) error {
				var r sql.Result
				var e1 error
				if avatarPath != "" {
					r, e1 = session.Exec("insert into users(nick_name,user_name,password_hash,email,avatar) values(?,?,?,?,?)", in.NickName, in.UserName, hashpwd, in.Email, avatarPath)
					if e1 != nil {
						return e1
					}
				} else {
					r, e1 = session.Exec("insert into users(nick_name,user_name,password_hash,email) values(?,?,?,?)", in.NickName, in.UserName, hashpwd, in.Email)
					if e1 != nil {
						return e1
					}
				}
				id, e2 := r.LastInsertId()
				if e2 != nil {
					return e2
				}
				_, e3 := session.Exec("insert into user_wallets(user_id) values(?)", id)
				if e3 != nil {
					return e3
				}
				return nil
			},
		)
		if err != nil {
			logc.Errorf(l.ctx, "[DBConn]用户创建失败:%s", err.Error())
			return &__.CreateUserResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "创建失败",
				},
			}, err
		}
		return &__.CreateUserResp{
			Base: &__.BaseResp{
				Code: 0,
				Msg:  "创建成功",
			},
		}, err
	}

	logc.Errorf(context.Background(), "[DBConn]其他错误:%s", err.Error())
	return &__.CreateUserResp{
		Base: &__.BaseResp{
			Code: 1,
			Msg:  "创建失败",
		},
	}, err
	//返回err累计熔断数据

}
