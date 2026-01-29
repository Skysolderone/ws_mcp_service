package rsi

import (
	"context"
	"fmt"
	"math"
	"mcp_service/model"
	"mcp_service/pkg/memcache"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	RsiChannel   = make(chan bool, 1)
	CloseChannel = make(chan bool, 1)
)

func CalcRsiTask() {
	logx.Info("RSI计算任务已启动")
	for {
		select {
		case <-RsiChannel:
			logx.Info("收到K线更新信号，开始计算RSI")
			CalcRsi()
			// 执行trade任务
			// TradeTask()
		case <-CloseChannel:
			logx.Info("收到退出信号，RSI计算任务停止")
			return
		}
	}
}

func TradeTask() {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	rsiValue := RsiMap[yesterday]

	logx.Infow("开始检查交易信号",
		logx.Field("日期", yesterday),
		logx.Field("RSI值", rsiValue),
		logx.Field("交易信号阈值", 30.0),
	)

	// 如果RsiMap中存在昨天的RSI值 并且小于30 则买入
	if rsiValue < 30 {
		logx.Infow("触发买入信号，准备下单",
			logx.Field("RSI值", rsiValue),
			logx.Field("操作", "买入"),
			logx.Field("交易对", "BTCUSDT"),
			logx.Field("数量", "0.001"),
			logx.Field("方向", "做多"),
			logx.Field("订单类型", "市价单"),
		)

		// 使用币安合约下单
		api := futures.NewClient("sAugoLUrKZUA5mRUeQIiL0CR0MaMFYkbhSeNrS3nZJDs9r5J4goXPxwUj2sOGQI7", "dXILNYaXZRdwjFnM17IKRltczkrlJwrLaADcJvCIsyYivfoPEopnI4iAjeSDFXGH")
		resp, err := api.NewCreateOrderService().Symbol("BTCUSDT").Side(futures.SideTypeBuy).Type(futures.OrderTypeMarket).Quantity("0.001").PositionSide(futures.PositionSideTypeLong).Do(context.Background())
		if err != nil {
			logx.Errorw("买入订单执行失败",
				logx.Field("错误", err.Error()),
				logx.Field("RSI值", rsiValue),
			)
			return
		}

		logx.Infow("买入订单执行成功",
			logx.Field("订单响应", resp),
			logx.Field("RSI值", rsiValue),
		)
	} else {
		logx.Infow("未触发交易信号",
			logx.Field("RSI值", rsiValue),
			logx.Field("原因", "RSI值未低于30"),
		)
	}
}

var RsiMap = make(map[string]float64)

func CalcRsi() {
	rsi := Rsi(model.KlineListModel.Klines, 14)
	// 使用昨天日期做key
	// yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// RsiMap[yesterday] = rsi
	memcache.SetMemcacheFloat("BTCUSDT", rsi)

	logx.Infow("RSI计算完成",
		logx.Field("交易对", "BTCUSDT"),
		logx.Field("RSI值", fmt.Sprintf("%.2f", rsi)),
		logx.Field("周期", 14),
		logx.Field("K线数量", len(model.KlineListModel.Klines)),
	)
}

// Rsi 计算RSI指标
// klines: K线数据切片
// period: RSI周期，通常使用14
// 返回: RSI值（0-100）
func Rsi(klines []model.Kline, period int) float64 {
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

func GetKline(symbol string) error {
	needInit := false
	var limit int = 2
	// 检查kline长度
	if model.KlineListModel.Len() == 0 {
		// 说明没有初始化
		needInit = true
		limit = 101
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

		klines, err := api.NewContinuousKlinesService().Limit(limit).ContractType("PERPETUAL").Pair(symbol).Interval("1d").Do(context.Background())
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
			model.KlineListModel.Add(model.Kline{
				OpenTime:  klinedata.OpenTime,
				CloseTime: klinedata.CloseTime,
				Open:      open,
				High:      high,
				Low:       low,
				Close:     close,
				Volume:    volume,
			})

		}
		logx.Infow("K线数据初始化完成",
			logx.Field("存储K线数量", model.KlineListModel.Len()),
			logx.Field("最早时间", time.UnixMilli(model.KlineListModel.Get(0).OpenTime).Format("2006-01-02")),
			logx.Field("最新时间", time.UnixMilli(model.KlineListModel.Get(model.KlineListModel.Len()-1).OpenTime).Format("2006-01-02")),
		)
		RsiChannel <- true
	} else {
		// 获取最新一条数据
		model.KlineListModel.RemoveFirst()
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
			model.KlineListModel.Add(model.Kline{
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
			logx.Field("存储K线数量", model.KlineListModel.Len()),
			logx.Field("最早时间", time.UnixMilli(model.KlineListModel.Get(0).OpenTime).Format("2006-01-02")),
			logx.Field("最新时间", time.UnixMilli(model.KlineListModel.Get(model.KlineListModel.Len()-1).OpenTime).Format("2006-01-02")),
		)
		RsiChannel <- true
	}
	return nil
}
