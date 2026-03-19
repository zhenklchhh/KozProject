package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/zhenklchhh/KozProject/order/internal/converter"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	"github.com/zhenklchhh/KozProject/order/internal/service"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

type api struct {
	service service.OrderService
}

func NewApi(s service.OrderService) *api {
	return &api{
		service: s,
	}
}

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	resp, err := a.service.Create(ctx, converter.CreateOrderRequestServiceToRepo(req))
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("order service: %v\n", err),
		}, nil
	}
	return converter.CreateOrderResponseRepoToService(resp), nil
}

func (a *api) PayOrder(ctx context.Context,
	req *orderV1.PayOrderRequest, params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	payResp, err := a.service.PayOrder(ctx, converter.PayOrderRequestServiceToRepo(req), params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBadRequest):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		case errors.Is(err, model.ErrNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		case errors.Is(err, model.ErrConflict):
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		}
	}
	return converter.PayOrderResponseRepoToService(payResp), nil
}

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.service.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		case errors.Is(err, model.ErrConflict):
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		default:
			return &orderV1.InternalServerError{
				Code:    500,
				Message: fmt.Sprintf("order service: %v\n", err),
			}, nil
		}
	}
	return &orderV1.CancelOrderNoContent{}, nil
}

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.service.Get(ctx, params.OrderUUID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order service: %v", err),
		}, nil
	}
	return converter.OrderRepoToService(order), nil
}
