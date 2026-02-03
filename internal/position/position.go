package position

import (
	"context"
	"encoding/json"
	"mcp_service/pb/position"
	"mcp_service/pkg/binance"
	"mcp_service/pkg/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

func PullPosition() {
	positions, err := binance.GetClient().NewGetPositionRiskV3Service().Do(context.Background())
	if err != nil {
		logx.Errorw("获取仓位失败", logx.Field("错误", err.Error()))
		return
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
	jsonData, err := json.Marshal(positionsList)
	if err != nil {
		logx.Errorw("序列化仓位失败", logx.Field("错误", err.Error()))
		return
	}
	redis.GetRedis().Set("POSITIONS", string(jsonData))
	logx.Infow("获取仓位成功", logx.Field("仓位数量", len(positionsList)))
}
