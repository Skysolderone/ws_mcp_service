package setup

import "github.com/zeromicro/zero-contrib/zrpc/registry/consul"

type ConsulConf struct {
	consul.Conf
	ServiceAddress string `json:",optional"`
}

type BinanceConf struct {
	ApiKey    string `json:"ApiKey"`
	ApiSecret string `json:"ApiSecret"`
}

type PostgreSQLConf struct {
	DSN string `json:"DSN"`
}
