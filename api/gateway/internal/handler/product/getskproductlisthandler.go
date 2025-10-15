package product

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"sk_mall/api/gateway/internal/logic/product"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
)

func GetSkProductListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSkProductListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewGetSkProductListLogic(r.Context(), svcCtx)
		resp, err := l.GetSkProductList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
