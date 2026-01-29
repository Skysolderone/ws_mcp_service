package logic

import (
	"context"

	"mcp_service/pb/rsi"
	"mcp_service/pkg/memcache"
	"mcp_service/services/rsi/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetRsiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRsiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRsiLogic {
	return &GetRsiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRsiLogic) GetRsi(in *rsi.GetRsiRequest) (*rsi.GetRsiResponse, error) {
	if in.Symbol == "" {
		return nil, status.Error(codes.InvalidArgument, "symbol is required")
	}
	value := memcache.GetMemcacheFloat(in.Symbol)
	if value != 0 {
		return &rsi.GetRsiResponse{
			Symbol: in.Symbol,
			Rsi:    float32(value),
		}, nil
	}

	return &rsi.GetRsiResponse{
		Symbol: in.Symbol,
		Rsi:    -1,
	}, nil
}
