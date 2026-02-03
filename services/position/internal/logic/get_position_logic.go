package logic

import (
	"context"

	"mcp_service/pb/position"
	"mcp_service/services/position/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPositionLogic {
	return &GetPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPositionLogic) GetPosition(in *position.GetPositionRequest) (*position.PositionList, error) {

	positions, err := l.svcCtx.Cli.NewGetPositionRiskV3Service().Do(l.ctx)
	if err != nil {
		return nil, err
	}
	positionsList := make([]*position.PositionDetail, 0)
	for _, p := range positions {
		positionsList = append(positionsList, &position.PositionDetail{
			Symbol:           p.Symbol,
			Quantity:         p.PositionAmt,
			EntryPrice:       p.EntryPrice,
			UnrealizedPnl:    p.UnRealizedProfit,
			Side:             p.PositionSide,
			LiquidationPrice: p.LiquidationPrice,
		})
	}
	return &position.PositionList{
		Positions: positionsList,
	}, nil
}
