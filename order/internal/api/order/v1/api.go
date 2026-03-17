package api

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/zhenklchhh/KozProject/order/internal/converter"
	"github.com/zhenklchhh/KozProject/order/internal/service"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	paymentServiceURI   = "localhost:50052"
	inventoryServiceURI = "localhost:50051"
)

type api struct {
	service         service.OrderService
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewApi(s service.OrderService, invClient inventoryV1.InventoryServiceClient, payClient paymentV1.PaymentServiceClient) *api {
	inventoryClient, paymentClient := invClient, payClient
	if inventoryClient == nil {
		inventoryClient = initGrpcInventoryClient()
	}
	if paymentClient == nil {
		paymentClient = initGrpcPaymentClient()
	}
	return &api{
		service: s,
		inventoryClient: inventoryClient,
		paymentClient: paymentClient,
	}
}

func initGrpcInventoryClient() inventoryV1.InventoryServiceClient {
	conn, err := grpc.NewClient(inventoryServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create client: %v\n", err)
	}
	return inventoryV1.NewInventoryServiceClient(conn)
}

func initGrpcPaymentClient() paymentV1.PaymentServiceClient {
	conn, err := grpc.NewClient(paymentServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create client: %v\n", err)
	}
	return paymentV1.NewPaymentServiceClient(conn)
}

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	invResp, err := a.inventoryClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: req.PartUuids,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get list parts from inventory service: %v\n", err)
	}
	if len(invResp.GetParts()) != len(req.PartUuids) {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "some parts don't exist in inventory service",
		}, nil
	}
	totalPrice := 0.0
	for _, part := range invResp.GetParts() {
		totalPrice += part.Price
	}
	newUUID := uuid.New()
	resp := &orderV1.CreateOrderResponse{
		OrderUUID:  newUUID.String(),
		TotalPrice: totalPrice,
	}
	order := &orderV1.Order{
		OrderUUID:  newUUID.String(),
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}
	err = a.service.Create(ctx, converter.OrderServiceToRepo(order))
	if err != nil {
		return nil, fmt.Errorf("order service error: %v\n", err)
	}
	return resp, nil
}

func (a *api) PayOrder(ctx context.Context,
	req *orderV1.PayOrderRequest, params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	orderRepo, err := a.service.Get(ctx, params.OrderUUID)
	order := converter.OrderRepoToService(orderRepo)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found: %s", params.OrderUUID, err.Error()),
		}, nil
	}
	if order.GetStatus() != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order %s has status '%s' and cannot be paid", params.OrderUUID, order.GetStatus()),
		}, nil
	}
	paymentMethodNum, ok := paymentV1.PaymentMethod_value[string(req.PaymentMethod)]
	if !ok {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("invalid payment method: %v\n", req.PaymentMethod),
		}, nil
	}
	payResp, err := a.paymentClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentMethodNum),
	})
	if err != nil {
		return nil, fmt.Errorf("payment service error: %v\n", err)
	}
	order.SetStatus(orderV1.OrderStatusPAID)
	order.SetTransactionUUID(orderV1.NewOptNilString(payResp.GetTransactionUuid()))
	order.SetPaymentMethod(orderV1.NewOptPaymentMethod(req.PaymentMethod))
	err = a.service.Create(ctx, converter.OrderServiceToRepo(order))
	if err != nil {
		return nil, fmt.Errorf("payment service error: %v\n", err)
	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: payResp.TransactionUuid,
	}, nil
}

func (a *api) CancelOrder(ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	orderRepo, err := a.service.Get(ctx, params.OrderUUID)
	order := converter.OrderRepoToService(orderRepo)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found: %s", params.OrderUUID, err.Error()),
		}, nil
	}
	if order.GetStatus() != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order %s can't be cancelled, order status = %s", params.OrderUUID, order.GetStatus()),
		}, nil
	}
	order.SetStatus(orderV1.OrderStatusCANCELLED)
	err = a.service.Create(ctx, converter.OrderServiceToRepo(order))
	if err != nil {
		return &orderV1.InternalServerError{
			Code: 500,
			Message: fmt.Sprintf("internal error: %s", err.Error()),
		}, nil
	}	
	return &orderV1.CancelOrderNoContent{}, nil
}

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.service.Get(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found: %s", params.OrderUUID, err.Error()),
		}, nil
	}
	return converter.OrderRepoToService(order), nil
}
