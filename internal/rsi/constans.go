package rsi

var IntervalList = map[string]string{
	"1h": "0 * * * *",
	"2h": "0 */2 * * *",
	"4h": "0 */4 * * *",
	"1d": "0 0 * * *",
	"1w": "0 0 * * 1",
	"1M": "0 0 1 * *",
}

var SymbolList = []string{"BTCUSDT", "ETHUSDT", "XRPUSDT", "YFIUSDT", "WLDUSDT", "ADAUSDT", "ETHBTC"}
