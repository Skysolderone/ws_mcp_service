package logic

import (
	"context"
	"strconv"

	"mcp_service/pb/order"
	"mcp_service/services/order/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderListLogic) GetOrderList(in *order.GetOrderListRequest) (*order.GetOrderListResponse, error) {
	orders, err := l.svcCtx.Binance.NewListOrdersService().Do(l.ctx)
	if err != nil {
		return nil, err
	}
	ordersList := make([]*order.OrderDetail, 0)
	for _, orderObject := range orders {
		ordersList = append(ordersList, &order.OrderDetail{
			OrderId:   strconv.FormatInt(orderObject.OrderID, 10),
			Symbol:    orderObject.Symbol,
			Side:      string(orderObject.Side),
			Type:      string(orderObject.Type),
			Quantity:  orderObject.OrigQuantity,
			Price:     orderObject.Price,
			StopPrice: orderObject.StopPrice,
		})
	}
	return &order.GetOrderListResponse{
		Orders:   ordersList,
		Total:    strconv.Itoa(len(orders)),
		Page:     "1",
		PageSize: "100",
		Code:     0,
	}, nil
}
