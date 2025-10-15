package merchant

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"sk_mall/api/gateway/internal/logic/merchant"
	"sk_mall/api/gateway/internal/middleware"
	"sk_mall/api/gateway/internal/svc"
	"sk_mall/api/gateway/internal/types"
	"sk_mall/utils"

	"github.com/zeromicro/go-zero/rest/httpx"
)



func CreateMerchantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateMerchantReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//lincense img
		imgData, handler, e1 := r.FormFile("license")
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
		req.LincenseSuffix = suffix
		req.LincenseImg, e1 = io.ReadAll(imgData)
		if e1 != nil {
			httpx.ErrorCtx(r.Context(), w, e1)
			return
		}
		//logo
		logoData, logohandler, e2 := r.FormFile("logo")
		if e2 != nil {
			httpx.ErrorCtx(r.Context(), w, e2)
			return
		}
		defer logoData.Close()
		if logohandler.Size > 2*1024*1024 {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("文件过大"))
			return
		}
		logonameList := strings.Split(logohandler.Filename, ".")
		logosuffix := strings.ToLower(logonameList[len(logonameList)-1])
		if !utils.Inlist(logosuffix, svc.EnableImageList) {
			//不在白名单
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("图片格式错误"))
			return
		}
		req.LogoSuffix = logosuffix
		req.Logo, e1 = io.ReadAll(logoData)
		if e1 != nil {
			httpx.ErrorCtx(r.Context(), w, e1)
			return
		}
		userId := r.Context().Value(middleware.UserIdKey).(int64)
		req.UserId = strconv.Itoa(int(userId))

		l := merchant.NewCreateMerchantLogic(r.Context(), svcCtx)
		resp, err := l.CreateMerchant(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
