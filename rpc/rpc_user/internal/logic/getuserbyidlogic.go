package logic

import (
	"context"
	"database/sql"
	"errors"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type UserInfo struct {
	Id       int64   `db:"id"`
	NickName string  `db:"nick_name"`
	UserName string  `db:"user_name"`
	Balance  float64 `db:"balance"`
}

func (l *GetUserByIdLogic) GetUserById(in *__.GetUserByIdReq) (*__.GetUserByIdResp, error) {
	// todo: add your logic here and delete this line

	var userinfo UserInfo
	err := l.svcCtx.DBConn.QueryRow(&userinfo, "select id,nick_name,user_name,balance from users,user_wallets where id = ? and user_id = ?", in.Id, in.Id)

	var resp __.GetUserByIdResp
	if errors.Is(err, sql.ErrNoRows) {
		//找不到记录
		resp = __.GetUserByIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "用户不存在",
			},
		}
		err = nil
	} else if err == nil {
		resp = __.GetUserByIdResp{
			Base: &__.BaseResp{
				Code: 0,
				Msg:  "success",
			},
			Data: &__.UserInfo{
				Id:       userinfo.Id,
				NickName: userinfo.NickName,
				UserName: userinfo.UserName,
				Balance:  userinfo.Balance,
			},
		}
	} else {
		//其他错误
		logc.Errorf(l.ctx, "[ConnDB]err:%s", err.Error())
		resp = __.GetUserByIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "其他错误",
			},
		}
	}

	return &resp, err
}
