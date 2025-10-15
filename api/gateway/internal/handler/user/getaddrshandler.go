package user

import (
	"net/http"
	"strconv"

	"sk_mall/api/gateway/internal/logic/user"
	"sk_mall/api/gateway/internal/middleware"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAddrsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetAddrsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		userId := r.Context().Value(middleware.UserIdKey).(int64)
		req.Id = strconv.Itoa(int(userId))
		l := user.NewGetAddrsLogic(r.Context(), svcCtx)
		resp, err := l.GetAddrs(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
