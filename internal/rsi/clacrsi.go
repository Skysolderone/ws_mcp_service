package rsi

import (
	"context"
	"fmt"
	"math"
	"mcp_service/internal/db"
	"mcp_service/model"
	"mcp_service/pkg/memcache"
	"strconv"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	RsiChannel      = make(chan model.KlineList, 1)
	CloseChannel    = make(chan bool, 1)
	SaveRsiChannel  = make(chan model.Rsi, 1)
	RsiCheckChannel = make(chan model.KlineList, 1)
	// SaveDoneWG 用于在工具程序中等所有 RSI 落库后再退出
	SaveDoneWG sync.WaitGroup
)

func CalcRsiTask() {
	logx.Info("RSI计算任务已启动")
	for {
		select {
		case data := <-RsiChannel:
			logx.Info("收到K线更新信号，开始计算RSI")

			calcRsi(data)
			// 执行trade任务
			// TradeTask()
		case data := <-RsiCheckChannel:
			logx.Info("收到RSI检查信号，开始检查RSI")
			calcRsiCheck(data)

		case <-CloseChannel:
			logx.Info("收到退出信号，RSI计算任务停止")
			return
		}
	}
}

func SaveRsiTask() {
	logx.Info("RSI保存任务已启动")
	for {
		select {
		case data := <-SaveRsiChannel:
			logx.Info("收到保存RSI信号，开始保存RSI")
			saveRsi(data)

		case <-CloseChannel:
			logx.Info("收到退出信号，RSI保存任务停止")
			return
		}
	}
}

var RsiMap = make(map[string]float64)

func calcRsiCheck(klineList model.KlineList) {
	symbol := klineList.Symbol
	interval := klineList.Interval
	klines := klineList.Klines
	if len(klines) == 0 {
		return
	}

	for i := 0; i < 102; i++ {

		newklines := klines[i : i+100]
		rsiVal := rsi(newklines, 14)
		// 该 RSI 对应最后一根 K 线的周期结束时间，存为 2006-01-02 15:04:05
		lastCloseMs := newklines[len(newklines)-1].CloseTime
		timeStr := time.Unix(lastCloseMs/1000, 0).UTC().Format("2006-01-02 15:04:05")
		memcache.SetMemcacheFloat(fmt.Sprintf("%s_%s", symbol, interval), rsiVal)

		logx.Infow("RSI计算完成",
			logx.Field("交易对", symbol),
			logx.Field("RSI值", fmt.Sprintf("%.2f", rsiVal)),
			logx.Field("周期", 14),
			logx.Field("K线数量", len(newklines)),
		)
		// SaveRsiChannel <- model.Rsi{
		// 	Symbol:   symbol,
		// 	Value:    rsiVal,
		// 	Time:     timeStr,
		// 	Interval: interval,
		// }
		saveRsi(model.Rsi{
			Symbol:   symbol,
			Value:    rsiVal,
			Time:     timeStr,
			Interval: interval,
		})
	}
	SaveDoneWG.Done()
}

func calcRsi(klineList model.KlineList) {
	symbol := klineList.Symbol
	interval := klineList.Interval
	klines := klineList.Klines
	if len(klines) == 0 {
		return
	}
	rsiVal := rsi(klines, 14)
	// 该 RSI 对应最后一根 K 线的周期结束时间，存为 2006-01-02 15:04:05
	lastCloseMs := klines[len(klines)-1].CloseTime
	timeStr := time.Unix(lastCloseMs/1000, 0).UTC().Format("2006-01-02 15:04:05")
	memcache.SetMemcacheFloat(fmt.Sprintf("%s_%s", symbol, interval), rsiVal)

	logx.Infow("RSI计算完成",
		logx.Field("交易对", symbol),
		logx.Field("RSI值", fmt.Sprintf("%.2f", rsiVal)),
		logx.Field("周期", 14),
		logx.Field("K线数量", len(klines)),
	)
	SaveRsiChannel <- model.Rsi{
		Symbol:   symbol,
		Value:    rsiVal,
		Time:     timeStr,
		Interval: interval,
	}
}

// Rsi 计算RSI指标
// klines: K线数据切片
// period: RSI周期，通常使用14
// 返回: RSI值（0-100）
func rsi(klines []model.Kline, period int) float64 {
	if len(klines) < period+1 {
		return 0 // 数据不足，无法计算
	}

	var gains, losses float64

	// 计算第一个周期的平均涨跌幅
	for i := 1; i <= period; i++ {
		change := klines[i].Close - klines[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += math.Abs(change)
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// 使用Wilder平滑方法计算后续周期
	for i := period + 1; i < len(klines); i++ {
		change := klines[i].Close - klines[i-1].Close
		if change > 0 {
			avgGain = (avgGain*float64(period-1) + change) / float64(period)
			avgLoss = (avgLoss * float64(period-1)) / float64(period)
		} else {
			avgGain = (avgGain * float64(period-1)) / float64(period)
			avgLoss = (avgLoss*float64(period-1) + math.Abs(change)) / float64(period)
		}
	}

	// 避免除零错误
	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}

func GetKline(symbol string, interval string) error {
	needInit := false
	var limit int = 2
	// 检查kline长度
	klinelist := model.KlineListModel[symbol]
	if len(klinelist.Klines) == 0 {
		// 说明没有初始化
		needInit = true
		limit = 101
		model.KlineListModel[symbol] = model.KlineList{
			Symbol:   symbol,
			Interval: interval,
			Klines:   make([]model.Kline, 0),
		}
		logx.Infow("开始初始化K线数据",
			logx.Field("交易对", symbol),
			logx.Field("获取数量", limit),
		)
	} else {
		logx.Infow("开始获取最新K线数据",
			logx.Field("交易对", symbol),
			logx.Field("获取数量", limit),
		)
	}
	api := futures.NewClient("", "")
	if needInit {
		// 初始化kline
		// 使用币安客户端获取合约历史一百条数据

		klines, err := api.NewContinuousKlinesService().Limit(limit).ContractType("PERPETUAL").Pair(symbol).Interval(interval).Do(context.Background())
		if err != nil {
			logx.Errorw("初始化K线数据失败",
				logx.Field("错误", err.Error()),
				logx.Field("交易对", symbol),
			)
			return err
		}
		logx.Infow("成功获取初始K线数据",
			logx.Field("交易对", symbol),
			logx.Field("数据条数", len(klines)),
		)
		for _, klinedata := range klines {
			// 如果openTime大于time.Now().AddDate(0, 0, -1).Unix()，则跳过
			if klinedata.OpenTime > time.Now().AddDate(0, 0, -1).UnixMilli() {
				fmt.Println("openTime大于time.Now().AddDate(0, 0, -1).UnixMilli()", klinedata.OpenTime, time.Now().AddDate(0, 0, -1).UnixMilli())
				continue
			}
			open, _ := strconv.ParseFloat(klinedata.Open, 64)
			high, _ := strconv.ParseFloat(klinedata.High, 64)
			low, _ := strconv.ParseFloat(klinedata.Low, 64)
			close, _ := strconv.ParseFloat(klinedata.Close, 64)
			volume, _ := strconv.ParseFloat(klinedata.Volume, 64)
			klinelist := model.KlineListModel[symbol]
			klinelist.Add(model.Kline{
				OpenTime:  klinedata.OpenTime,
				CloseTime: klinedata.CloseTime,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    volume,
			})
			model.KlineListModel[symbol] = klinelist

		}
		logx.Infow("K线数据初始化完成",
			logx.Field("存储K线数量", len(model.KlineListModel[symbol].Klines)),
			logx.Field("最早时间", time.UnixMilli(model.KlineListModel[symbol].Klines[0].OpenTime).Format("2006-01-02")),
			logx.Field("最新时间", time.UnixMilli(model.KlineListModel[symbol].Klines[len(model.KlineListModel[symbol].Klines)-1].OpenTime).Format("2006-01-02")),
		)
		SaveDoneWG.Add(1)
		RsiChannel <- model.KlineList{
			Symbol:   symbol,
			Interval: interval,
			Klines:   model.KlineListModel[symbol].Klines,
		}
	} else {
		// 获取最新一条数据
		klinelist := model.KlineListModel[symbol]
		klinelist.RemoveFirst()
		model.KlineListModel[symbol] = klinelist
		klines, err := api.NewContinuousKlinesService().Limit(limit).ContractType("PERPETUAL").Pair(symbol).Interval("1d").Do(context.Background())
		if err != nil {
			logx.Errorw("获取最新K线数据失败",
				logx.Field("错误", err.Error()),
				logx.Field("交易对", symbol),
			)
			return err
		}

		open, _ := strconv.ParseFloat(klines[0].Open, 64)
		high, _ := strconv.ParseFloat(klines[0].High, 64)
		low, _ := strconv.ParseFloat(klines[0].Low, 64)
		close, _ := strconv.ParseFloat(klines[0].Close, 64)
		volume, _ := strconv.ParseFloat(klines[0].Volume, 64)
		for _, klinedata := range klines {
			// 如果openTime大于time.Now().AddDate(0, 0, -1).Unix()，则跳过
			if klinedata.OpenTime > time.Now().AddDate(0, 0, -1).UnixMilli() {
				fmt.Println("openTime大于time.Now().AddDate(0, 0, -1).UnixMilli()", klinedata.OpenTime, time.Now().AddDate(0, 0, -1).UnixMilli())
				continue
			}
			klinelist := model.KlineListModel[symbol]
			klinelist.Add(model.Kline{
				OpenTime:  klinedata.OpenTime,
				CloseTime: klinedata.CloseTime,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    volume,
			})
		}
		logx.Infow("K线数据更新完成",
			logx.Field("存储K线数量", len(model.KlineListModel[symbol].Klines)),
			logx.Field("最早时间", time.UnixMilli(model.KlineListModel[symbol].Klines[0].OpenTime).Format("2006-01-02")),
			logx.Field("最新时间", time.UnixMilli(model.KlineListModel[symbol].Klines[len(model.KlineListModel[symbol].Klines)-1].OpenTime).Format("2006-01-02")),
		)
		SaveDoneWG.Add(1)
		RsiChannel <- model.KlineList{
			Symbol:   symbol,
			Interval: interval,
			Klines:   klinelist.Klines,
		}
	}
	return nil
}

func saveRsi(rsi model.Rsi) {
	err := db.GetPostgreSQL().DB.Create(&rsi).Error
	if err != nil {
		logx.Errorw("保存RSI失败",
			logx.Field("错误", err.Error()),
			logx.Field("RSI", rsi),
		)
	}
	logx.Infow("RSI保存成功",
		logx.Field("RSI", rsi),
		logx.Field("交易对", rsi.Symbol),
		logx.Field("周期", rsi.Interval),
		logx.Field("时间", rsi.Time),
	)
}

func ParseDate(interval string, i int) (time.Time, time.Time) {
	now := time.Now().UTC()
	switch interval {
	case "1d":
		// 上一自然日结束 = 今天 0 点，再往前 i 天；startTime = 再往前 100 天
		date := now.AddDate(0, 0, -1).Format("2006-01-02")                 //这里是昨天
		endTime, err := time.ParseInLocation("2006-01-02", date, time.UTC) //
		if err != nil {
			logx.Errorw("日期格式错误", logx.Field("错误", err.Error()), logx.Field("日期", date))
			return time.Time{}, time.Time{}
		}
		return endTime.AddDate(0, 0, -200), endTime
	case "4h":
		// 与 1d 一样按“日期”处理，只是回溯天数不同
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		endTime, err := time.ParseInLocation("2006-01-02", date, time.UTC)
		if err != nil {
			logx.Errorw("日期格式错误", logx.Field("错误", err.Error()), logx.Field("日期", date))
			return time.Time{}, time.Time{}
		}
		// 约覆盖最近 20 天的 4h 数据
		startTime := endTime.AddDate(0, 0, -20)
		return startTime, endTime
	case "1h":
		date := now.Truncate(time.Hour).Format("2006-01-02 15:04:05")
		endTime, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.UTC)
		if err != nil {
			logx.Errorw("日期格式错误", logx.Field("错误", err.Error()), logx.Field("日期", date))
			return time.Time{}, time.Time{}
		}
		// 覆盖最近 5 天的 1h 数据
		startTime := endTime.Truncate(time.Hour * 100)
		return startTime, endTime
	case "2h":
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		endTime, err := time.ParseInLocation("2006-01-02", date, time.UTC)
		if err != nil {
			logx.Errorw("日期格式错误", logx.Field("错误", err.Error()), logx.Field("日期", date))
			return time.Time{}, time.Time{}
		}
		// 覆盖最近 10 天的 2h 数据
		startTime := endTime.AddDate(0, 0, -10)
		return startTime, endTime
	case "1w":
		// 与 1d 一样先按天得到 endTime，再往前回溯 100 周对应的天数
		date := now.AddDate(0, 0, -i*7).Format("2006-01-02")
		endTime, err := time.ParseInLocation("2006-01-02", date, time.UTC)
		if err != nil {
			logx.Errorw("日期格式错误", logx.Field("错误", err.Error()), logx.Field("日期", date))
			return time.Time{}, time.Time{}
		}
		startTime := endTime.AddDate(0, 0, -100*7)
		return startTime, endTime
	case "1M":
		// 与 1d 一样使用日期格式，但按月回溯
		date := now.AddDate(0, -i, 0).Format("2006-01-02")
		endTime, err := time.ParseInLocation("2006-01-02", date, time.UTC)
		if err != nil {
			logx.Errorw("日期格式错误", logx.Field("错误", err.Error()), logx.Field("日期", date))
			return time.Time{}, time.Time{}
		}
		startTime := endTime.AddDate(0, -100, 0)
		return startTime, endTime
	}
	return time.Time{}, time.Time{}
}

func GetKlineByDate(symbol string, interval string, startTime time.Time, endTime time.Time) {

	api := futures.NewClient("", "")
	//这里已经拿到所有的rsi数据 可以计算 一百天的rsi数据
	klines, err := api.NewContinuousKlinesService().StartTime(startTime.UnixMilli()).ContractType("PERPETUAL").EndTime(endTime.UnixMilli()).Pair(symbol).Interval(interval).Do(context.Background())
	if err != nil {
		logx.Errorw("初始化K线数据失败",
			logx.Field("错误", err.Error()),
			logx.Field("交易对", symbol),
			logx.Field("日期", endTime.Format("2006-01-02")),
		)
		return
	}
	logx.Infow("成功获取初始K线数据",
		logx.Field("交易对", symbol),
		logx.Field("数据条数", len(klines)),
	)
	if len(klines) == 0 {
		err := fmt.Errorf("该日期无K线数据: %s", endTime.Format("2006-01-02"))
		logx.Errorw("K线数据为空", logx.Field("错误", err.Error()), logx.Field("交易对", symbol), logx.Field("日期", endTime.Format("2006-01-02")))
		return
	}

	for _, klinedata := range klines {

		open, _ := strconv.ParseFloat(klinedata.Open, 64)
		high, _ := strconv.ParseFloat(klinedata.High, 64)
		low, _ := strconv.ParseFloat(klinedata.Low, 64)
		close, _ := strconv.ParseFloat(klinedata.Close, 64)
		volume, _ := strconv.ParseFloat(klinedata.Volume, 64)
		klinelist := model.KlineListModel[symbol]
		klinelist.Add(model.Kline{
			OpenTime:  klinedata.OpenTime,
			CloseTime: klinedata.CloseTime,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    volume,
		})
		model.KlineListModel[symbol] = klinelist

	}
	stored := model.KlineListModel[symbol].Klines
	if len(stored) == 0 {
		err := fmt.Errorf("过滤后无有效K线数据: %s", endTime.Format("2006-01-02"))
		logx.Errorw("K线数据为空", logx.Field("错误", err.Error()), logx.Field("交易对", symbol), logx.Field("日期", endTime.Format("2006-01-02")))
		return
	}
	logx.Infow("K线数据初始化完成",
		logx.Field("存储K线数量", len(stored)),
		logx.Field("最早时间", time.UnixMilli(stored[0].OpenTime).Format("2006-01-02")),
		logx.Field("最新时间", time.UnixMilli(stored[len(stored)-1].OpenTime).Format("2006-01-02")),
	)
	RsiCheckChannel <- model.KlineList{
		Symbol:   symbol,
		Interval: interval,
		Klines:   stored,
	}

}
