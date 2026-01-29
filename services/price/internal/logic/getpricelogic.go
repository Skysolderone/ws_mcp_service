package logic

import (
	"context"

	"mcp_service/pb/price"
	"mcp_service/pkg/memcache"
	"mcp_service/services/price/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetPriceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPriceLogic {
	return &GetPriceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPriceLogic) GetPrice(in *price.Symbol) (*price.Price, error) {
	memprice := memcache.GetMemcacheFloat("BTCUSDT_MARK_PRICE")
	if memprice == 0 {
		return nil, status.Error(codes.NotFound, "price not found")
	}
	return &price.Price{
		Symbol: in.Symbol,
		Price:  (memprice),
	}, nil
}
