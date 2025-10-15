package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sk_mall/rpc/rpc_user/user"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JwtAuthMiddleware struct {
	UserRpc user.User
}

func NewJwtAuthMiddleware(u user.User) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		UserRpc: u,
	}
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type contextKey string

const (
	UserIdKey contextKey = "UserId"
)

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			httpx.OkJsonCtx(r.Context(), w,
				resp{
					Code: 998,
					Msg:  "未携带Authorization",
				})
			return
		}
		gresp, e1 := m.UserRpc.JwtAuth(r.Context(), &user.JwtAuthReq{Token: token})
		if e1 != nil {
			logc.Errorf(r.Context(), "[UserRpc] JwtAuth err:%s", e1.Error())
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("server err"))
			return
		}
		if gresp.Base.Code != 0 {
			httpx.OkJsonCtx(r.Context(), w,
				resp{
					Code: int(gresp.Base.Code),
					Msg:  gresp.Base.Msg,
				})
			return
		}
		reqCtx := r.Context()
		val := gresp.UserId
		ctx := context.WithValue(reqCtx, UserIdKey, val)
		newReq := r.WithContext(ctx)
		next(w, newReq)
		// Passthrough to next handler if need
	}
}
