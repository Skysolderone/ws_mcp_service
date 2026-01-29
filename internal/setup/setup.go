package setup

import (
	"context"
	"mcp_service/internal/rsi"
	"mcp_service/internal/websocket"
	"mcp_service/pkg/memcache"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

func Setup(service string) {
	memcache.InitMemcache()
	switch service {
	case "price":
		go websocket.MarkPriceTask()

	case "rsi":
		go websocket.MarkPriceTask()
		go rsi.CalcRsiTask()
		// 首次启动立即获取K线数据
		if err := rsi.GetKline("BTCUSDT"); err != nil {
			logx.WithContext(context.Background()).Errorf("首次获取K线数据失败", map[string]interface{}{
				"错误":  err.Error(),
				"交易对": "BTCUSDT",
			})
			return
		}
		timer := cron.New()
		timer.AddFunc("0 0 * * *", func() {
			logx.WithContext(context.Background()).Infof("触发定时任务：开始获取K线数据", map[string]interface{}{
				"交易对": "BTCUSDT",
			})
			if err := rsi.GetKline("BTCUSDT"); err != nil {
				logx.WithContext(context.Background()).Errorf("定时任务执行失败", map[string]interface{}{
					"错误":  err.Error(),
					"交易对": "BTCUSDT",
				})
				return
			}
		})
		timer.Start()

	}

}
