package svc

import (
	"mcp_service/services/position/internal/config"

	"github.com/adshao/go-binance/v2/futures"
)

type ServiceContext struct {
	Config config.Config
	Cli    *futures.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Cli:    futures.NewClient(c.Binance.ApiKey, c.Binance.ApiSecret),
	}
}
