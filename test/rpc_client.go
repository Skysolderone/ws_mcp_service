package test

import (
	"log"

	"google.golang.org/grpc"
)

var GrpcConn *grpc.ClientConn

func InitGrpcClient() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	GrpcConn = conn
}
