package product

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"sk_mall/api/gateway/internal/logic/product"
	"sk_mall/api/gateway/internal/middleware"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateProductReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//cover
		imgData, handler, e1 := r.FormFile("cover")
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
		req.CoverSuffix = suffix
		req.CoverImg, e1 = io.ReadAll(imgData)
		if e1 != nil {
			httpx.ErrorCtx(r.Context(), w, e1)
			return
		}
		//userid
		userId := r.Context().Value(middleware.UserIdKey).(int64)
		req.UserId = strconv.Itoa(int(userId))

		l := product.NewCreateProductLogic(r.Context(), svcCtx)
		resp, err := l.CreateProduct(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
