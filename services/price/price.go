package main

import (
	"flag"
	"fmt"

	"mcp_service/internal/setup"
	"mcp_service/pb/price"
	"mcp_service/services/price/internal/config"
	"mcp_service/services/price/internal/server"
	"mcp_service/services/price/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/price.yaml", "the config file")

func main() {
	flag.Parse()
	setup.Setup("price")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		price.RegisterPriceServiceServer(grpcServer, server.NewPriceServiceServer(ctx))

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
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
