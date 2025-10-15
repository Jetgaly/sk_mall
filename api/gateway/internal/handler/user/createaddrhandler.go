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

func CreateAddrHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateAddrReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		userId := r.Context().Value(middleware.UserIdKey).(int64)
		req.UserId = strconv.Itoa(int(userId))
		l := user.NewCreateAddrLogic(r.Context(), svcCtx)
		resp, err := l.CreateAddr(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
