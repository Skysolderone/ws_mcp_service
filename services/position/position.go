package main

import (
	"flag"

	"mcp_service/internal/setup"
	"mcp_service/pb/position"
	"mcp_service/services/position/internal/config"
	"mcp_service/services/position/internal/server"
	"mcp_service/services/position/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/position.yaml", "the config file")

func main() {
	flag.Parse()
	setup.Setup("position")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		position.RegisterPositionServer(grpcServer, server.NewPositionServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	err := consul.RegisterService(c.Consul.ServiceAddress, c.Consul.Conf)
	if err != nil {
		logx.Errorf("Register service to consul failed: %v", err)
		return
	}
	logx.Infof("Register service to consul success")
	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}
