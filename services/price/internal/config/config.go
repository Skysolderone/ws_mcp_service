package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type ConsulConf struct {
	consul.Conf
	ServiceAddress string `json:",optional"`
}

type Config struct {
	zrpc.RpcServerConf
	Consul ConsulConf
}
