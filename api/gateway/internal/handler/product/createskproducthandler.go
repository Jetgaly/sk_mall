package product

import (
	"net/http"
	"strconv"

	"sk_mall/api/gateway/internal/logic/product"
	"sk_mall/api/gateway/internal/middleware"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateSkProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateSkProductReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//userid
		userId := r.Context().Value(middleware.UserIdKey).(int64)
		req.UserId = strconv.Itoa(int(userId))

		l := product.NewCreateSkProductLogic(r.Context(), svcCtx)
		resp, err := l.CreateSkProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
