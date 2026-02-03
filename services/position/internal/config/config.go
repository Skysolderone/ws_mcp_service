package config

import (
	"mcp_service/internal/setup"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Consul  setup.ConsulConf
	Binance setup.BinanceConf
	Redis   redis.RedisConf
}
