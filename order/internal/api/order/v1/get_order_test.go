package api

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

func (s *ApiSuit) TestGetOrderSuccess() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	transactionUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard

	tc := &struct {
		req           *orderV1.GetOrderParams
		order         *model.Order
		expectedOrder *orderV1.Order
	}{
		req: &orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		},
		order: &model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUuids,
			TotalPrice:      15.99,
			Status:          model.OrderStatusPaid,
			TransactionUUID: &transactionUUID,
			PaymentMethod:   &paymentMethod,
		},
		expectedOrder: &orderV1.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     orderV1.OrderStatusPAID,
			TransactionUUID: orderV1.OptNilString{
				Value: transactionUUID,
				Set:   true,
				Null:  false,
			},
			PaymentMethod: orderV1.OptPaymentMethod{
				Value: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
				Set:   true,
			},
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	response, err := s.handler.GetOrder(s.ctx, *tc.req)
	s.Require().NoError(err)
	s.Require().Equal(tc.expectedOrder, response)
	s.service.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
}

func (s *ApiSuit) TestGetOrderNotFound() {
	orderUUID := gofakeit.UUID()

	tc := &struct {
		req          *orderV1.GetOrderParams
		expectedResp orderV1.GetOrderRes
	}{
		req: &orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		},
		expectedResp: &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order service: %v", errors.New("not found")),
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(nil, errors.New("not found"))
	response, err := s.handler.GetOrder(s.ctx, *tc.req)
	s.Require().NoError(err)
	s.Require().Equal(tc.expectedResp, response)
	s.service.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
}

func (s *ApiSuit) TestGetOrderPendingPaymentStatus() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}

	tc := &struct {
		req           *orderV1.GetOrderParams
		order         *model.Order
		expectedOrder *orderV1.Order
	}{
		req: &orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		},
		order: &model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 10.50,
			Status:     model.OrderStatusPendingPayment,
		},
		expectedOrder: &orderV1.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 10.50,
			Status:     orderV1.OrderStatusPENDINGPAYMENT,
			TransactionUUID: orderV1.OptNilString{
				Value: "",
				Set:   true,
				Null:  true,
			},
			PaymentMethod: orderV1.OptPaymentMethod{
				Set: false,
			},
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	response, err := s.handler.GetOrder(s.ctx, *tc.req)
	s.Require().NoError(err)
	s.Require().Equal(tc.expectedOrder, response)
	s.service.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
}

func (s *ApiSuit) TestGetOrderCancelledStatus() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}

	tc := &struct {
		req           *orderV1.GetOrderParams
		order         *model.Order
		expectedOrder *orderV1.Order
	}{
		req: &orderV1.GetOrderParams{
			OrderUUID: orderUUID,
		},
		order: &model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 25.99,
			Status:     model.OrderStatusCancelled,
		},
		expectedOrder: &orderV1.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 25.99,
			Status:     orderV1.OrderStatusCANCELLED,
			TransactionUUID: orderV1.OptNilString{
				Value: "",
				Set:   true,
				Null:  true,
			},
			PaymentMethod: orderV1.OptPaymentMethod{
				Set: false,
			},
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	response, err := s.handler.GetOrder(s.ctx, *tc.req)
	s.Require().NoError(err)
	s.Require().Equal(tc.expectedOrder, response)
	s.service.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
}
