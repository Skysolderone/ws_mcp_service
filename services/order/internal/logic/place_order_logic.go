package logic

import (
	"context"
	"strconv"
	"time"

	"mcp_service/pb/order"
	"mcp_service/services/order/internal/svc"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/zeromicro/go-zero/core/logx"
)

type PlaceOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlaceOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlaceOrderLogic {
	return &PlaceOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PlaceOrderLogic) PlaceOrder(in *order.PlaceOrderRequest) (*order.PlaceOrderResponse, error) {
	orderService := l.svcCtx.Binance.NewCreateOrderService()
	orderService = orderService.Symbol(in.Symbol)
	orderService = orderService.Side(futures.SideType(in.Side))
	orderService = orderService.Type(futures.OrderType(in.Type))
	orderService = orderService.Quantity(in.Quantity)
	orderService = orderService.Price(in.Price)
	orderService = orderService.StopPrice(in.StopPrice)
	switch in.PositionSide {
	case "LONG":
		orderService = orderService.PositionSide(futures.PositionSideTypeLong)
	case "SHORT":
		orderService = orderService.PositionSide(futures.PositionSideTypeShort)
	case "BOTH":
		orderService = orderService.PositionSide(futures.PositionSideTypeBoth)
	}
	resp, err := orderService.Do(l.ctx)
	if err != nil {
		return nil, err
	}

	return &order.PlaceOrderResponse{
		OrderId:   strconv.FormatInt(resp.OrderID, 10),
		Status:    string(resp.Status),
		Message:   "下单成功",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
