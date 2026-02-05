package config

import (
	"mcp_service/internal/setup"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Consul     setup.ConsulConf
	PostgreSQL setup.PostgreSQLConf
}
