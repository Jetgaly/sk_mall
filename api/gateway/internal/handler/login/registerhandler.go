package login

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"sk_mall/api/gateway/internal/logic/login"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)



func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//读取头像
		imgData, handler, e1 := r.FormFile("avatar")
		if e1 != nil {
			httpx.ErrorCtx(r.Context(), w, e1)
			return
		}
		defer imgData.Close()
		if handler.Size > 2*1024*1024 {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("文件过大"))
			return
		}
		nameList := strings.Split(handler.Filename, ".")
		suffix := strings.ToLower(nameList[len(nameList)-1])
		if !utils.Inlist(suffix, svc.EnableImageList) {
			//不在白名单
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("图片格式错误"))
			return
		}
		req.Suffix = suffix
		req.Avatar, e1 = io.ReadAll(imgData)
		if e1 != nil {
			httpx.ErrorCtx(r.Context(), w, e1)
			return
		}

		l := login.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			//400
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
			//200
		}
	}
}
