package binance

import (
	"github.com/adshao/go-binance/v2/futures"
)

var Client *futures.Client

func GetClient() *futures.Client {
	if Client != nil {
		return Client
	}
	Client = futures.NewClient("", "")
	return Client
}
