package loginlogic

import (
	"context"

	"overall/modules/rpc/auth/internal/svc"
	"overall/modules/rpc/auth/pb/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 退出登录
func (l *LogoutLogic) Logout(in *login.LogoutReq) (*login.LogoutResq, error) {
	// todo: add your logic here and delete this line

	return &login.LogoutResq{}, nil
}
