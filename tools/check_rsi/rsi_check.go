package main

import (
	"log"
	"mcp_service/internal/db"
	"mcp_service/internal/rsi"
	"mcp_service/model"
)

func main() {
	go rsi.CalcRsiTask()
	go rsi.SaveRsiTask()
	db.InitPostgreSQL("host=pgm-bp140jpn9wct9u0two.pg.rds.aliyuncs.com user=wws password=Wws5201314 dbname=rsi port=5432 sslmode=disable TimeZone=UTC")

	// 检查每个 (symbol, interval) 是否有至少 100 条 RSI，不足则按周期补全
	rsiList := []model.Rsi{}
	rsi.SymbolList = []string{"ETHBTC"}
	rsi.IntervalList = map[string]string{
		// "1d": "0 0 * * *",
		// 	// "4h": "0 */4 * * *",
		"1h": "0 * * * *",
		// 	// "2h": "0 */2 * * *",
		// 	// "1w": "0 0 * * 1",
		// 	// "1M": "0 0 1 * *",
	}
	rsi.SaveDoneWG.Add(1)
	for _, symbol := range rsi.SymbolList {
		for interval := range rsi.IntervalList {
			err := db.GetPostgreSQL().DB.Where("symbol = ? AND interval = ?", symbol, interval).Order("time DESC").Limit(100).Find(&rsiList).Error
			if err != nil {
				log.Fatalf("获取RSI数据失败: %v", err)
			}
			if len(rsiList) < 100 {
				log.Printf("RSI数据不足: %s %s 当前%d条，需补全至100条", symbol, interval, len(rsiList))
				//一条数据没有则需要补全

				startTime, endTime := rsi.ParseDate(interval, 100)
				if startTime.IsZero() || endTime.IsZero() {
					log.Printf("日期格式错误: %s %s", symbol, interval)
					continue
				}
				log.Printf("%#v,%#v", startTime, endTime)
				rsi.GetKlineByDate(symbol, interval, startTime, endTime)

			}

		}
	}
	rsi.SaveDoneWG.Wait()
	log.Printf("RSI 补全完成，已全部落库")
}
