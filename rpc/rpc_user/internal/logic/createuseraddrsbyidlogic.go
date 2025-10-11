package logic

import (
	"context"
	"errors"

	"sk_mall/rpc/rpc_user/internal/svc"
	"sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateUserAddrsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserAddrsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserAddrsByIdLogic {
	return &CreateUserAddrsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type Addrs_count struct {
	Count int `db:"addrs_count"`
}

func (l *CreateUserAddrsByIdLogic) CreateUserAddrsById(in *__.CreateAddrsByIdReq) (*__.CreateAddrsByIdResp, error) {
	// todo: add your logic here and delete this line
	var ac Addrs_count
	var resp __.CreateAddrsByIdResp

	err := l.svcCtx.DBConn.QueryRow(&ac, "select addrs_count from users where id = ?", in.Id)
	if errors.Is(err, sqlx.ErrNotFound) {
		resp = __.CreateAddrsByIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "用户不存在",
			},
		}
		//不启动熔断
		err = nil
	} else if err != nil {
		logc.Errorf(l.ctx, "[DBConn]err:%s", err.Error())
		resp = __.CreateAddrsByIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "其他错误",
			},
		}
	} else {
		//err == nil
		if ac.Count >= 10 {
			resp = __.CreateAddrsByIdResp{
				Base: &__.BaseResp{
					Code: 1,
					Msg:  "地址数量超过限制，限制为10",
				},
			}
		} else {
			if ac.Count == 0 {
				in.Addr.IsDefault = 1
			}
			err = l.svcCtx.DBConn.Transact(func(session sqlx.Session) error {
				if ac.Count > 0 && in.Addr.IsDefault == 1 {
					_, e1 := session.Exec("update user_addresses set is_default=0 where user_id = ? and is_default = 1", in.Id)
					if e1 != nil {
						logc.Errorf(l.ctx, "[DBConn] update err:%s", e1.Error())
						return e1
					}
				}

				_, e1 := session.Exec("insert into user_addresses(user_id,receiver_name,receiver_phone,province,city,district,detail_address,is_default) values(?,?,?,?,?,?,?,?)", in.Id, in.Addr.ReceiverName, in.Addr.ReceiverPhone, in.Addr.Province, in.Addr.City, in.Addr.District, in.Addr.DetailAddress, in.Addr.IsDefault)
				if e1 != nil {
					logc.Errorf(l.ctx, "[DBConn] insert err:%s", e1.Error())
					return e1
				}
				_, e1 = session.Exec("update users set addrs_count=addrs_count+1 where id = ?", in.Id)
				if e1 != nil {
					logc.Errorf(l.ctx, "[DBConn] update err:%s", e1.Error())
				}
				return e1
			})
			if err != nil {
				logc.Errorf(l.ctx, "[DBConn] err:%s", err.Error())
				resp = __.CreateAddrsByIdResp{
					Base: &__.BaseResp{
						Code: 1,
						Msg:  "其他错误",
					},
				}
			} else {
				resp = __.CreateAddrsByIdResp{
					Base: &__.BaseResp{
						Code: 0,
						Msg:  "success",
					},
				}
			}
		}
	}
	return &resp, err
}
