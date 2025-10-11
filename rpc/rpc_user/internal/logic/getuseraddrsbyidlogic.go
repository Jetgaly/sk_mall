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

type GetUserAddrsByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserAddrsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAddrsByIdLogic {
	return &GetUserAddrsByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type Addr struct {
	Id            int    `db:"id"`
	ReceiverName  string `db:"receiver_name"`
	ReceiverPhone string `db:"receiver_phone"`
	Province      string `db:"province"`
	City          string `db:"city"`
	District      string `db:"district"`
	DetailAddress string `db:"detail_address"`
	IsDefault     int32  `db:"is_default"`
}

func (l *GetUserAddrsByIdLogic) GetUserAddrsById(in *__.GetAddrsByIdReq) (*__.GetAddrsByIdResp, error) {
	// todo: add your logic here and delete this line
	var data []*__.Address
	var addrs []Addr
	var resp __.GetAddrsByIdResp
	err := l.svcCtx.DBConn.QueryRows(&addrs, "select id,receiver_name,receiver_phone,province,city,district,detail_address,is_default from user_addresses where user_id = ?", in.Id)
	if errors.Is(err, sql.ErrNoRows) || err == nil {

		resp = __.GetAddrsByIdResp{
			Base: &__.BaseResp{
				Code: 0,
				Msg:  "success",
			},
			Data: nil,
		}
		err = nil
		if len(addrs) != 0 {
			for _, a := range addrs {
				data = append(data, &__.Address{
					AddrId:        int64(a.Id),
					ReceiverName:  a.ReceiverName,
					ReceiverPhone: a.ReceiverPhone,
					Province:      a.Province,
					City:          a.City,
					District:      a.District,
					DetailAddress: a.DetailAddress,
					IsDefault:     a.IsDefault,
				})
			}
		}
		resp.Data = data
	} else {
		//其他错误
		logc.Errorf(l.ctx, "[DBConn] query err:%s", err.Error())
		resp = __.GetAddrsByIdResp{}
	}
	return &resp, err
}
