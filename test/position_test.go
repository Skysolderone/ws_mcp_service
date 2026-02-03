package test

import (
	"context"
	"fmt"
	"mcp_service/pb/position"
	"testing"
)

func TestPosition(t *testing.T) {
	InitGrpcClient()
	client := position.NewPositionClient(GrpcConn)
	positions, err := client.GetPosition(context.Background(), &position.GetPositionRequest{})
	if err != nil {
		t.Fatalf("Failed to get positions: %v", err)
	}
	fmt.Println(positions)
}
