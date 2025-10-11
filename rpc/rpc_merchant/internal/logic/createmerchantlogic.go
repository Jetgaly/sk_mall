package logic

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"sk_mall/rpc/rpc_merchant/internal/svc"
	"sk_mall/rpc/rpc_merchant/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateMerchantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMerchantLogic {
	return &CreateMerchantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

type MerChantCheck struct {
	Id int `db:"id"`
}

func (l *CreateMerchantLogic) CreateMerchant(in *__.CreateMerchantReq) (*__.CreateMerchantResp, error) {
	// todo: add your logic here and delete this line
	//查找是否已经创建了merchant
	var flag MerChantCheck
	err := l.svcCtx.DBConn.QueryRowCtx(l.ctx, &flag, "select id from merchants where user_id = ?", in.UserId)
	if errors.Is(err, sqlx.ErrNotFound) {
		err = nil
	} else if err == nil {
		return &__.CreateMerchantResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "该账号已经创建了商户",
			},
		}, nil
	} else {
		logc.Errorf(l.ctx, "[DBConn] query err:%s", err.Error())
		return &__.CreateMerchantResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "未知错误",
			},
		}, err
	}
	logoName := uuid.NewString() + "." + in.MInfo.LogoSuffix //api传后缀
	logoPath := filepath.Join(l.svcCtx.Config.Logo.UploadPath, logoName)

	//如果写文件失败logo为空字符串
	if err := os.WriteFile(logoPath, in.MInfo.Logo, 0644); err != nil {
		logc.Errorf(l.ctx, "[LogoUpload]文件写入失败: %s", err.Error())
		logoPath = ""
	}

	licenseName := uuid.NewString() + "." + in.AInfo.LicenseSuffix //api传后缀
	licensePath := filepath.Join(l.svcCtx.Config.License.UploadPath, licenseName)

	//如果写文件失败logo为空字符串
	if err := os.WriteFile(licensePath, in.AInfo.LicenseImg, 0644); err != nil {
		logc.Errorf(l.ctx, "[LicenseUpload]文件写入失败: %s", err.Error())
		return &__.CreateMerchantResp{
			Base: &__.BaseResp{
				Code: 1,
				Msg:  "未知错误",
			},
		}, err
	}

	err = l.svcCtx.DBConn.Transact(func(session sqlx.Session) error {

		//merchant
		r, e1 := session.ExecCtx(l.ctx, "insert into merchants(user_id,name,logo,description) values(?,?,?,?)", in.UserId, in.MInfo.Name, logoPath, in.MInfo.Description)
		if e1 != nil {
			logc.Errorf(l.ctx, "[DBConn] insert err:%s", e1.Error())
			return e1
		}
		mId, e := r.LastInsertId()
		if e != nil {
			logc.Errorf(l.ctx, "[DBConn] err:%s", e.Error())
			return e
		}
		//acc
		_, e2 := session.ExecCtx(l.ctx, "insert into merchant_accounts(merchant_id,account_type,account_no,account_name,bank_name) values(?,?,?,?,?)", mId, in.AcInfo.AccType, in.AcInfo.AccNo, in.AcInfo.AccName, in.AcInfo.BankName)
		if e2 != nil {
			logc.Errorf(l.ctx, "[DBConn] insert err:%s", e2.Error())
			return e2
		}

		//auth
		_, e3 := session.ExecCtx(l.ctx, "insert into merchant_auths(merchant_id,legal_name,id_card,license_img) values(?,?,?,?)", mId, in.AInfo.LegalName, in.AInfo.IdCard, licensePath)
		if e3 != nil {
			logc.Errorf(l.ctx, "[DBConn] insert err:%s", e3.Error())
			return e3
		}
		
		return nil
	})
	if err != nil {
		logc.Errorf(l.ctx, "[DBConn] err:%s", err.Error())
	}
	return &__.CreateMerchantResp{
		Base: &__.BaseResp{
			Code: 0,
			Msg:  "success",
		},
	}, err
}
