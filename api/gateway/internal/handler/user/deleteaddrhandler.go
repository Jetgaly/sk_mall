package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sk_mall/api/gateway/internal/logic/user"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
)

func DeleteAddrHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteAddrReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewDeleteAddrLogic(r.Context(), svcCtx)
		resp, err := l.DeleteAddr(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
