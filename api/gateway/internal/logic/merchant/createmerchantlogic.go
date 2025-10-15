package merchant

import (
	"context"
	"encoding/json"
	"strconv"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_merchant/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMerchantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMerchantLogic {
	return &CreateMerchantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type (
	MerchantInfo struct {
		Name string `json:"name"`
		Desc string `json:"description"`
		Type int64  `json:"type"`
	}
	AuthInfo struct {
		LegalName string `json:"legal_name"`
		IdCard    string `json:"idcard"`
	}
	AccountInfo struct {
		AccType  int64  `json:"acc_type"`
		AccNo    string `json:"acc_no"`
		AccName  string `json:"acc_name"`
		BankName string `json:"bank_name"`
	}
	Info struct {
		MInfo   MerchantInfo `json:"m_info"`
		AInfo   AuthInfo     `json:"a_info"`
		AccInfo AccountInfo  `json:"acc_info"`
	}
)

func (l *CreateMerchantLogic) CreateMerchant(req *types.CreateMerchantReq) (resp *types.CreateMerchantResp, err error) {
	err = nil
	var model Info
	e1 := json.Unmarshal([]byte(req.JsonInfo), &model)
	if e1 != nil {
		resp = &types.CreateMerchantResp{
			Code: 997,
			Msg:  e1.Error(),
		}
		err = e1
		return
	}
	userId, _ := strconv.Atoi(req.UserId)
	//验证参数
	var ok bool
	model.MInfo.Desc, ok = utils.CheckString(model.MInfo.Desc, 10, 100)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 20,
			Msg:  "desc为10-100个字符",
		}
		return
	}
	model.MInfo.Name, ok = utils.CheckString(model.MInfo.Name, 3, 10)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 21,
			Msg:  "merchant name为3-10个字符",
		}
		return
	}

	if model.MInfo.Type != 1 && model.MInfo.Type != 2 {
		resp = &types.CreateMerchantResp{
			Code: 22,
			Msg:  "merchant type必须为1或2",
		}
		return
	}
	model.AInfo.IdCard, ok = utils.CheckString(model.AInfo.IdCard, 18, 18)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 23,
			Msg:  "身份证号码不对,必须是18位",
		}
		return
	}
	model.AInfo.LegalName, ok = utils.CheckString(model.AInfo.LegalName, 1, 18)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 24,
			Msg:  "姓名非法",
		}
		return
	}
	model.AccInfo.AccName, ok = utils.CheckString(model.AccInfo.AccName, 1, 18)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 25,
			Msg:  "账户名非法",
		}
		return
	}
	model.AccInfo.AccNo, ok = utils.CheckString(model.AccInfo.AccName, 16, 19)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 26,
			Msg:  "银行卡号非法",
		}
		return
	}
	model.AccInfo.BankName, ok = utils.CheckString(model.AccInfo.AccName, 1, 16)
	if !ok {
		resp = &types.CreateMerchantResp{
			Code: 27,
			Msg:  "银行非法",
		}
		return
	}
	if model.AccInfo.AccType != 1 && model.AccInfo.AccType != 2 && model.AccInfo.AccType != 3 {
		resp = &types.CreateMerchantResp{
			Code: 28,
			Msg:  "账户类型非法",
		}
		return
	}
	gresp, e2 := l.svcCtx.MerchantRpc.CreateMerchant(l.ctx, &__.CreateMerchantReq{
		UserId: uint64(userId),
		MInfo: &__.MerchantInfo{
			Name:        model.MInfo.Name,
			Logo:        req.Logo,
			LogoSuffix:  req.LogoSuffix,
			Description: model.MInfo.Desc,
			Type:        int32(model.MInfo.Type),
		},
		AInfo: &__.AuthInfo{
			LegalName:     model.AInfo.LegalName,
			IdCard:        model.AInfo.IdCard,
			LicenseImg:    req.LincenseImg,
			LicenseSuffix: req.LincenseSuffix,
		},
		AcInfo: &__.AccInfo{
			AccNo:    model.AccInfo.AccNo,
			AccType:  int32(model.AccInfo.AccType),
			AccName:  model.AccInfo.AccName,
			BankName: model.AccInfo.BankName,
		},
	})

	if e2 != nil {
		logc.Errorf(l.ctx, "[MerchantRpc] CreateMerchant err:%s", e2.Error())
		resp = &types.CreateMerchantResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	resp = &types.CreateMerchantResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	err = nil

	return
}
