package websocket

import (
	"fmt"
	"mcp_service/pkg/memcache"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
)

func MarkPriceTask() {
	wsKlineHandler := func(event *futures.WsMarkPriceEvent) {
		markPrice, _ := strconv.ParseFloat(event.MarkPrice, 64)
		memcache.SetMemcacheFloat("BTCUSDT_MARK_PRICE", markPrice)
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := futures.WsCombinedMarkPriceServe([]string{"BTCUSDT"}, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC

}
