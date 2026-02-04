package logic

import (
	"context"
	"strconv"
	"time"

	"mcp_service/pb/order"
	"mcp_service/services/order/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	orderID, err := strconv.ParseInt(in.OrderId, 10, 64)
	if err != nil {
		return nil, err
	}
	resp, err := l.svcCtx.Binance.NewCancelOrderService().OrderID(orderID).Do(l.ctx)
	if err != nil {
		return nil, err
	}
	orderIDStr := strconv.FormatInt(resp.OrderID, 10)
	return &order.CancelOrderResponse{
		OrderId:   orderIDStr,
		Status:    "success",
		Message:   "订单取消成功",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
