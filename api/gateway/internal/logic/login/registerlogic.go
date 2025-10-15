package login

import (
	"context"
	"encoding/json"
	"strings"
	"unicode/utf8"

	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	__ "sk_mall/rpc/rpc_user/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type userInfo struct {
	NickName string `json:"nickname"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	var usrModel userInfo
	e2 := json.Unmarshal([]byte(req.UserInfo), &usrModel)
	if e2 != nil {
		err = e2
		return
	}
	usrModel.NickName = strings.ReplaceAll(usrModel.NickName, " ", "")
	if utf8.RuneCountInString(usrModel.NickName) > 10 || utf8.RuneCountInString(usrModel.NickName) <= 0 {
		resp = &types.RegisterResp{
			Code: 11,
			Msg:  "昵称需要1-10个字符",
		}
		err = nil
		return
	}
	usrModel.UserName = strings.ReplaceAll(usrModel.UserName, " ", "")
	if utf8.RuneCountInString(usrModel.UserName) > 16 || utf8.RuneCountInString(usrModel.UserName) < 8 {
		resp = &types.RegisterResp{
			Code: 12,
			Msg:  "账号需要8-16个字符",
		}
		err = nil
		return
	}
	usrModel.Password = strings.ReplaceAll(usrModel.Password, " ", "")
	if utf8.RuneCountInString(usrModel.Password) > 24 || utf8.RuneCountInString(usrModel.Password) <= 10 {
		resp = &types.RegisterResp{
			Code: 13,
			Msg:  "密码需要10-24个字符",
		}
		err = nil
		return
	}
	usrModel.Email = strings.ReplaceAll(usrModel.Email, " ", "")
	if utf8.RuneCountInString(usrModel.Email) == 0 {
		resp = &types.RegisterResp{
			Code: 14,
			Msg:  "请输入邮箱",
		}
		err = nil
		return
	}
	usrModel.Code = strings.ReplaceAll(usrModel.Code, " ", "")
	if utf8.RuneCountInString(usrModel.Code) == 0 {
		resp = &types.RegisterResp{
			Code: 15,
			Msg:  "请输入验证码",
		}
		err = nil
		return
	}
	gresp, e1 := l.svcCtx.UserRpc.CreateUser(l.ctx, &__.CreateUserReq{
		NickName:   usrModel.NickName,
		UserName:   usrModel.UserName,
		Password:   usrModel.Password,
		Email:      usrModel.Email,
		Code:       usrModel.Code,
		AvatarName: req.Suffix,
		Avatar:     req.Avatar,
	})
	if e1 != nil {
		logc.Errorf(l.ctx, "[UserRpc] CreateUser err:%s", e1.Error())
		resp = &types.RegisterResp{
			Code: 999,
			Msg:  "server err",
		}
		err = nil
		return
	}
	resp = &types.RegisterResp{
		Code: int(gresp.Base.Code),
		Msg:  gresp.Base.Msg,
	}
	err = nil
	return
}
