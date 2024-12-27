package middleware

/**
 * @Author nico
 * @Date 2024-12-27
 * @File: jwtAuthMiddleware.go
 * @Description:
 */

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"overall/common/library/reqCtx"
	"overall/common/library/response"
	"overall/common/xerr"
	"overall/modules/rpc/auth/client/priv"
	"strings"
)

type AuthActionMiddleware struct {
	PrivilegeRpc priv.Priv
}

func NewAuthActionMiddleware(rpc priv.Priv) *AuthActionMiddleware {
	return &AuthActionMiddleware{
		PrivilegeRpc: rpc,
	}
}

func (m *AuthActionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		sessionData := reqCtx.GetSessionFromCtx(ctx)

		token, flag := GetTokenFromHeader(r)
		if !flag {
			httpx.WriteJson(w, http.StatusUnauthorized, response.Error(xerr.RequestParamError, xerr.MapErrMsg(xerr.RequestParamError), "", nil))
			return
		}
		// 检查 token redis 过期
		expire, err := m.PrivilegeRpc.CheckTokenExpire(ctx, &priv.CheckTokenExpireReq{
			UserId: sessionData.UserId,
			Token:  token,
		})
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, response.Error(xerr.ActionNoPermission, xerr.MapErrMsg(xerr.ActionNoPermission), "", nil))
			return
		}
		if expire.IsExpire {
			httpx.WriteJson(w, http.StatusUnauthorized, response.Error(xerr.ActionNoPermission, xerr.MapErrMsg(xerr.ActionNoPermission), "", nil))
			return
		}

		errCode := xerr.ServerBusy
		errMsg := "ServerBusy"

		slice := strings.Split(r.URL.Path, "/")
		l := len(slice)
		if l <= 4 {
			httpx.WriteJson(w, http.StatusOK, response.Error(errCode, errMsg, "", nil))
			return
		}

		// menukey := strings.Join(slice[3:l-1], "/")
		// action := slice[l-1]

		// 获取页面权限
		// result, err := m.PrivilegeRpc.GetPagePriv(ctx, &priv.GetPagePrivReq{
		// 	UserId:      session.UserId,
		// 	RoleID:      session.RoleId,
		// 	PageKey:     menukey,
		// })
		// if err != nil {
		// 	httpx.WriteJson(w, http.StatusOK, response.Error(errCode, i18n.TSArgs("ActionNoPermission1", session.Lang, i18n.TS(menukey, session.Lang), i18n.TS(action, session.Lang)), "", nil))
		// 	return
		// }
		//
		// // 是否有这一项权限
		// if _, ok := result.Actions[action]; !ok {
		// 	httpx.WriteJson(w, http.StatusOK, response.Error(errCode, i18n.TSArgs("ActionNoPermission1", session.Lang, i18n.TS(menukey, session.Lang), i18n.TS(action, session.Lang)), "", nil))
		// 	return
		// }
		//
		// // 这一项权限是否开启了
		// if !result.Actions[action] {
		// 	httpx.WriteJson(w, http.StatusOK, response.Error(errCode, i18n.TSArgs("ActionNoPermission1", session.Lang, i18n.TS(menukey, session.Lang), i18n.TS(action, session.Lang)), "", nil))
		// 	return
		// }

		next(w, r)
	}
}

// GetTokenFromHeader 获取请求 Token
func GetTokenFromHeader(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}
	// 去掉 "Bearer " 前缀，获取实际的 token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, true
}
