package api

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

func (s *ApiSuit) TestCancelOrderSuccess() {
	OrderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}

	tc := &struct {
		req            orderV1.CancelOrderParams
		order          *model.Order
		convertedOrder *orderV1.Order
		expectedRes    orderV1.CancelOrderRes
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: OrderUUID,
		},
		expectedRes: &orderV1.CancelOrderNoContent{},
		order: &model.Order{
			OrderUUID:  OrderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPendingPayment,
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	var cancelledOrder interface{}
	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
		cancelledOrder = args.Get(1)
	})
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().Equal(resp, tc.expectedRes)
	s.Require().NoError(err)
	s.Require().Equal(cancelledOrder.(*model.Order).Status, model.OrderStatusCancelled)
	s.service.AssertCalled(s.T(), "Get", s.ctx, tc.req.OrderUUID)
}

func (s *ApiSuit) TestCancelOrderServiceError() {
	OrderUUID := gofakeit.UUID()
	tc := &struct {
		req         orderV1.CancelOrderParams
		expectedRes *orderV1.NotFoundError
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: OrderUUID,
		},
		expectedRes: &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found: %s", OrderUUID, "not found"),
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(nil, errors.New("not found"))
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.NotFoundError), tc.expectedRes)
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *ApiSuit) TestCancelOrderNotPendingPaymentStatus() {
	OrderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	tc := &struct {
		req         orderV1.CancelOrderParams
		order       *model.Order
		expectedRes *orderV1.ConflictError
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: OrderUUID,
		},
		expectedRes: &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order %s can't be cancelled, order status = %s", OrderUUID, model.OrderStatusPaid),
		},
		order: &model.Order{
			OrderUUID:  OrderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPaid,
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.ConflictError), tc.expectedRes)
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *ApiSuit) TestCancelOrderSaveServiceError() {
	OrderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	tc := &struct {
		req            orderV1.CancelOrderParams
		order          *model.Order
		cancelledOrder *model.Order
		expectedRes    *orderV1.InternalServerError
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: OrderUUID,
		},
		expectedRes: &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("internal error: %s", "err"),
		},
		order: &model.Order{
			OrderUUID:  OrderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPendingPayment,
		},
		cancelledOrder: &model.Order{
			OrderUUID:  OrderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusCancelled,
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	s.service.On("Create", s.ctx, tc.cancelledOrder).Return(errors.New("err"))
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.InternalServerError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "Get", s.ctx, OrderUUID)
}
