package logic

import (
	"context"
	"fmt"

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
	if in.Interval == "" {
		in.Interval = "1d"
	}
	key := fmt.Sprintf("%s_%s", in.Symbol, in.Interval)
	rsiValue := memcache.GetMemcacheFloat(key)
	if rsiValue == 0 {
		return nil, status.Error(codes.NotFound, "RSI not found")
	}
	return &rsi.GetRsiResponse{
		Symbol: in.Symbol,
		Rsi:    float32(rsiValue),
	}, nil
}
