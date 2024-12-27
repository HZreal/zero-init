package reqCtx

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/logx"
	"overall/common/thirdpart/jwts"
)

var (
	CtxKeyJwtUserId   = "userId"
	CtxKeyJwtRoleId   = "roleId"
	CtxKeyJwtUsername = "username"
)

type Jwt struct {
	UserId      int64  `ctx:"userId"`
	Username    string `ctx:"username"`
	RoleId      int64  `ctx:"roleId"`   // 角色ID
	RoleType    int64  `ctx:"roleType"` // 角色类型
	PrivilegeID int64  `ctx:"privilegeId"`
	Lang        string `ctx:"lang"`
}

func GetSessionFromCtx(ctx context.Context) *jwts.JwtPayload {
	var jwt Jwt
	err := contextx.For(ctx, &jwt)
	if err != nil {
		logx.WithContext(ctx).Errorf("GetSessionFromCtx err : %+v", err)
	}

	return &jwts.JwtPayload{
		UserId:      jwt.UserId,
		RoleId:      jwt.RoleId,
		RoleType:    jwt.RoleType,
		PrivilegeId: jwt.PrivilegeID,
	}
}

func GetUidFromCtx(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(CtxKeyJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}
	return uid
}

func GetRoleIdFromCtx(ctx context.Context) int64 {
	var cid int64
	if jsonUid, ok := ctx.Value(CtxKeyJwtRoleId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			cid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("CtxKeyJwtRoleId err : %+v", err)
		}
	}
	return cid
}

func GetUsernameFromCtx(ctx context.Context) string {
	name, _ := ctx.Value(CtxKeyJwtUsername).(string)
	return name
}
