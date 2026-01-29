package model

// 由于rsi不需要全部数据即可计算  所以不需要数据库存储
// 只需要历史100条数据 就可以计算出rsi 使用内存存储
type Kline struct {
	OpenTime  int64
	CloseTime int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type KlineList struct {
	Klines []Kline
}

var KlineListModel = NewKlineList()

func (k *KlineList) Add(kline Kline) {
	k.Klines = append(k.Klines, kline)
}

func (k *KlineList) Get(index int) Kline {
	return k.Klines[index]
}

func (k *KlineList) Len() int {
	return len(k.Klines)
}

// 删除第一条数据
func (k *KlineList) RemoveFirst() {
	k.Klines = k.Klines[1:]
}

func (k *KlineList) RemoveLast() {
	k.Klines = k.Klines[:len(k.Klines)-1]
}

func (k *KlineList) Remove(index int) {
	k.Klines = append(k.Klines[:index], k.Klines[index+1:]...)
}

func NewKlineList() KlineList {
	return KlineList{
		Klines: make([]Kline, 0),
	}
}
