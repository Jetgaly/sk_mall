package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sk_mall/rpc/rpc_merchant/internal/svc"
	"sk_mall/rpc/rpc_merchant/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMerchantIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMerchantIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMerchantIdLogic {
	return &GetMerchantIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMerchantIdLogic) GetMerchantId(in *__.GetMerchantIdReq) (*__.GetMerchantIdResp, error) {
	// todo: add your logic here and delete this line
	var m MerChantCheck
	err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &m, "select id from merchants where user_id=?", in.UserId)
	fmt.Println(err)
	if errors.Is(err, sql.ErrNoRows) {
		return &__.GetMerchantIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "不存在该MId",
			},
		}, nil
	} else if err != nil {
		logc.Errorf(l.ctx,"[DBConn] query err:%s",err.Error())
		return &__.GetMerchantIdResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "未知错误",
			},
		}, err
	}

	return &__.GetMerchantIdResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
		MId: uint64(m.Id),
	}, nil
}
