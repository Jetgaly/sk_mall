package logic

import (
	"context"
	"database/sql"
	"errors"

	"sk_mall/rpc/rpc_merchant/internal/svc"
	"sk_mall/rpc/rpc_merchant/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantByMIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMerchantByMIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantByMIdLogic {
	return &GetMerchantByMIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type MInfo struct{
	MId uint64 `db:"id"`
	Name string `db:"name"`
	Logo string `db:"logo"`
	Desc string `db:"description"`
	Type int32 `db:"type"`
	Status int32 `db:"status"`
}

func (l *GetMerchantByMIdLogic) GetMerchantByMId(in *__.GetMerchantByMIdReq) (*__.GetMerchantByMIdResp, error) {
	// todo: add your logic here and delete this line
	var info MInfo
	err:=l.svcCtx.DBConn.QueryRowCtx(l.ctx,&info,"select id,name,logo,description,type,status from merchants where id = ?",in.MerchantId)
	if errors.Is(err,sql.ErrNoRows){
		return &__.GetMerchantByMIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg: "merchant 不存在",
			},
		}, nil
	}else if err!=nil{
		logc.Errorf(l.ctx,"[DBConn] query err:%s",err.Error())
		return &__.GetMerchantByMIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg: "未知错误",
			},
		}, err
	}
	return &__.GetMerchantByMIdResp{
		Base: &__.BaseResp{
				Code: 0,
				Msg: "success",
			},
		Info: &__.GMerchantInfo{
			MerchantId: info.MId,
			Name: info.Name,
			Logo: info.Logo,
			Desc: info.Desc,
			Type: info.Type,
			Status: info.Status,
		},
	}, nil
}
