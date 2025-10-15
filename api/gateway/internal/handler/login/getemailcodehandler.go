package login

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sk_mall/api/gateway/internal/logic/login"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
)

func GetEmailCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmailCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewGetEmailCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetEmailCode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
