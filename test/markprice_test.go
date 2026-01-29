package test

import (
	"mcp_service/internal/websocket"
	"testing"
)

func TestMarkPrice(t *testing.T) {
	websocket.MarkPriceTask()
}
