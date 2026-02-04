package svc

import (
	"mcp_service/services/order/internal/config"

	"github.com/adshao/go-binance/v2/futures"
)

type ServiceContext struct {
	Config  config.Config
	Binance *futures.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		Binance: futures.NewClient(c.Binance.ApiKey, c.Binance.ApiSecret),
	}
}
