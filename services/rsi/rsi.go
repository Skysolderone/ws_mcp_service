package main

import (
	"flag"
	"os"

	"mcp_service/internal/setup"
	"mcp_service/pb/rsi"
	"mcp_service/services/rsi/internal/config"
	"mcp_service/services/rsi/internal/server"
	"mcp_service/services/rsi/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/rsi.yaml", "the config file")

func main() {
	flag.Parse()
	setup.Setup("rsi")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	logx.MustSetup(c.Log)
	logx.AddWriter(logx.NewWriter(os.Stdout))
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rsi.RegisterRsiServer(grpcServer, server.NewRsiServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	err := consul.RegisterService(c.ListenOn, c.Consul)
	if err != nil {
		logx.Errorf("Register service to consul failed: %v", err)
		return
	}
	logx.Infof("Register service to consul success")

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}
