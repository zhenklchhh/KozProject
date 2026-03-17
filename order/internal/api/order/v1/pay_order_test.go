package api

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

func (s *ApiSuit) TestPayOrderSuccess() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	modelPaymentMethod := model.PaymentMethodCreditCard
	paymentMethodNum := paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(paymentMethod)])
	tc := &struct {
		req              *orderV1.PayOrderRequest
		params           orderV1.PayOrderParams
		payOrderRequest  *paymentV1.PayOrderRequest
		payOrderResponse *paymentV1.PayOrderResponse
		order            *model.Order
		payedOrder       *model.Order
		expectedRes      orderV1.PayOrderRes
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		payOrderRequest: &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethodNum,
		},
		payOrderResponse: &paymentV1.PayOrderResponse{
			TransactionUuid: transactionUUID,
		},
		order: &model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPendingPayment,
		},
		payedOrder: &model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUuids,
			TotalPrice:      15.99,
			Status:          model.OrderStatusPaid,
			TransactionUUID: &transactionUUID,
			PaymentMethod:   &modelPaymentMethod,
		},
		expectedRes: &orderV1.PayOrderResponse{
			TransactionUUID: transactionUUID,
		},
	}
	s.service.On("Get", s.ctx, orderUUID).Return(tc.order, nil)
	s.paymentClient.On("PayOrder", s.ctx, tc.payOrderRequest).Return(tc.payOrderResponse, nil)
	var payedOrder interface{}
	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
		payedOrder = args.Get(1)
	})
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.PayOrderResponse), tc.expectedRes)
	s.Require().NotNil(payedOrder)
	s.Require().Equal(payedOrder.(*model.Order).Status, model.OrderStatusPaid)
	s.Require().Equal(*payedOrder.(*model.Order).TransactionUUID, transactionUUID)
	s.Require().Equal(*payedOrder.(*model.Order).PaymentMethod, model.PaymentMethodCreditCard)
	s.service.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
	s.paymentClient.AssertCalled(s.T(), "PayOrder", s.ctx, tc.payOrderRequest)
}

func (s *ApiSuit) TestPayOrderNotFound() {
	orderUUID := gofakeit.UUID()
	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		expectedRes *orderV1.NotFoundError
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found: %s", orderUUID, "not found"),
		},
	}
	s.service.On("Get", s.ctx, orderUUID).Return(nil, errors.New("not found"))
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.NotFoundError), tc.expectedRes)
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
}

func (s *ApiSuit) TestPayOrderInvalidStatus() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	testCases := []struct {
		name   string
		status model.OrderStatus
	}{
		{
			name:   "already paid",
			status: model.OrderStatusPaid,
		},
		{
			name:   "cancelled",
			status: model.OrderStatusCancelled,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			orderUUID := gofakeit.UUID()
			req := &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: orderUUID,
			}
			order := &model.Order{
				OrderUUID:  orderUUID,
				UserUUID:   userUUID,
				PartUuids:  partUuids,
				TotalPrice: 15.99,
				Status:     tc.status,
			}
			expectedRes := &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("order %s has status '%s' and cannot be paid", orderUUID, tc.status),
			}

			s.service.On("Get", s.ctx, orderUUID).Return(order, nil)
			resp, err := s.handler.PayOrder(s.ctx, req, params)
			s.Require().NoError(err)
			s.Require().Equal(resp.(*orderV1.ConflictError), expectedRes)
			s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
			s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
		})
	}
}

func (s *ApiSuit) TestPayOrderInvalidPaymentMethod() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		order       *model.Order
		expectedRes *orderV1.BadRequestError
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethod("INVALID_METHOD"),
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		order: &model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPendingPayment,
		},
		expectedRes: &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("invalid payment method: %v\n", orderV1.PaymentMethod("INVALID_METHOD")),
		},
	}
	s.service.On("Get", s.ctx, orderUUID).Return(tc.order, nil)
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.BadRequestError), tc.expectedRes)
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
}

func (s *ApiSuit) TestPayOrderPaymentServiceError() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	paymentMethodNum := paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(paymentMethod)])
	tc := &struct {
		req             *orderV1.PayOrderRequest
		params          orderV1.PayOrderParams
		payOrderRequest *paymentV1.PayOrderRequest
		order           *model.Order
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		payOrderRequest: &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethodNum,
		},
		order: &model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPendingPayment,
		},
	}
	s.service.On("Get", s.ctx, orderUUID).Return(tc.order, nil)
	s.paymentClient.On("PayOrder", s.ctx, tc.payOrderRequest).Return(nil, errors.New("payment service error"))
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "payment service error")
	s.Require().Nil(resp)
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *ApiSuit) TestPayOrderSaveError() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	paymentMethodNum := paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(paymentMethod)])
	tc := &struct {
		req              *orderV1.PayOrderRequest
		params           orderV1.PayOrderParams
		payOrderRequest  *paymentV1.PayOrderRequest
		payOrderResponse *paymentV1.PayOrderResponse
		order            *model.Order
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		payOrderRequest: &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentMethodNum,
		},
		payOrderResponse: &paymentV1.PayOrderResponse{
			TransactionUuid: transactionUUID,
		},
		order: &model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPendingPayment,
		},
	}
	s.service.On("Get", s.ctx, orderUUID).Return(tc.order, nil)
	s.paymentClient.On("PayOrder", s.ctx, tc.payOrderRequest).Return(tc.payOrderResponse, nil)
	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("save error"))
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "payment service error")
	s.Require().Nil(resp)
}

func (s *ApiSuit) TestPayOrderDifferentPaymentMethods() {
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	
	testCases := []struct {
		name           string
		paymentMethod  orderV1.PaymentMethod
		expectedMethod model.PaymentMethod
	}{
		{
			name:           "credit card",
			paymentMethod:  orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
			expectedMethod: model.PaymentMethodCreditCard,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			orderUUID := gofakeit.UUID() // Generate unique UUID for each test case
			userUUID := gofakeit.UUID()
			transactionUUID := gofakeit.UUID()
			
			paymentMethodNum := paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[string(tc.paymentMethod)])
			req := &orderV1.PayOrderRequest{
				PaymentMethod: tc.paymentMethod,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: orderUUID,
			}
			payOrderRequest := &paymentV1.PayOrderRequest{
				OrderUuid:     orderUUID,
				UserUuid:      userUUID,
				PaymentMethod: paymentMethodNum,
			}
			payOrderResponse := &paymentV1.PayOrderResponse{
				TransactionUuid: transactionUUID,
			}
			order := &model.Order{
				OrderUUID:  orderUUID,
				UserUUID:   userUUID,
				PartUuids:  partUuids,
				TotalPrice: 15.99,
				Status:     model.OrderStatusPendingPayment,
			}

			s.service.On("Get", s.ctx, orderUUID).Return(order, nil)
			s.paymentClient.On("PayOrder", s.ctx, payOrderRequest).Return(payOrderResponse, nil)
			var payedOrder interface{}
			s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
				payedOrder = args.Get(1)
			})

			resp, err := s.handler.PayOrder(s.ctx, req, params)
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().NotNil(payedOrder, "payedOrder should not be nil - check if Create mock was called")
			s.Require().Equal(*payedOrder.(*model.Order).PaymentMethod, tc.expectedMethod)
		})
	}
}
