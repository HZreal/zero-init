package loginlogic

import (
	"context"

	"overall/modules/rpc/auth/internal/svc"
	"overall/modules/rpc/auth/pb/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccountLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountLoginLogic {
	return &AccountLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登录
func (l *AccountLoginLogic) AccountLogin(in *login.AccountLoginReq) (*login.LoginResp, error) {
	// todo: add your logic here and delete this line

	return &login.LoginResp{}, nil
}
