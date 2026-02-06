package test

import (
	"context"
	"mcp_service/internal/rsi"
	"testing"

	"github.com/adshao/go-binance/v2/futures"
)

func TestGetKline(t *testing.T) {
	start, end := rsi.ParseDate("1h", 100)
	t.Logf("start: %s, end: %s", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	api := futures.NewClient("", "")
	klines, err := api.NewContinuousKlinesService().StartTime(start.UnixMilli()).EndTime(end.UnixMilli()).ContractType("PERPETUAL").Pair("ETHBTC").Interval("1h").Do(context.Background())
	if err != nil {
		t.Fatalf("获取K线数据失败: %v", err)
	}
	t.Logf("klines: %d", len(klines))
	for _, kline := range klines {
		t.Logf("kline: %v", kline.OpenTime)
		t.Logf("kline: %v", kline.CloseTime)
		t.Logf("kline: %v", kline.Open)
		t.Logf("kline: %v", kline.High)
		t.Logf("kline: %v", kline.Low)
		t.Logf("kline: %v", kline.Close)
		t.Logf("kline: %v", kline.Volume)
	}
}
