package privlogic

import (
	"context"

	"overall/modules/rpc/auth/internal/svc"
	"overall/modules/rpc/auth/pb/priv"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckTokenExpireLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckTokenExpireLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckTokenExpireLogic {
	return &CheckTokenExpireLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查 token 过期
func (l *CheckTokenExpireLogic) CheckTokenExpire(in *priv.CheckTokenExpireReq) (*priv.CheckTokenExpireResp, error) {
	// todo: add your logic here and delete this line

	return &priv.CheckTokenExpireResp{}, nil
}
