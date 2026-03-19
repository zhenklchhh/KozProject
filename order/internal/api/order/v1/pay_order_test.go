package api

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

func (s *ApiSuit) TestPayOrderSuccess() {
	orderUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD

	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		expectedRes *orderV1.PayOrderResponse
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: paymentMethod,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.PayOrderResponse{
			TransactionUUID: transactionUUID,
		},
	}

	s.service.On("PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID).Return(&model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil)

	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.PayOrderResponse), tc.expectedRes)
	s.service.AssertCalled(s.T(), "PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID)
}

func (s *ApiSuit) TestPayOrderNotFound() {
	orderUUID := gofakeit.UUID()
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD

	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		expectedRes *orderV1.NotFoundError
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: paymentMethod,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order service: %s\n", "not found"),
		},
	}

	s.service.On("PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID).Return(nil, model.ErrNotFound)
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.NotFoundError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID)
}

func (s *ApiSuit) TestPayOrderInvalidStatus() {
	orderUUID := gofakeit.UUID()
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD

	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		expectedRes *orderV1.ConflictError
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: paymentMethod,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order service: %s\n", "conflict"),
		},
	}

	s.service.On("PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID).Return(nil, model.ErrConflict)
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.ConflictError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID)
}

func (s *ApiSuit) TestPayOrderInvalidPaymentMethod() {
	orderUUID := gofakeit.UUID()
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD

	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		expectedRes *orderV1.BadRequestError
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: paymentMethod,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("order service: %s\n", "bad request"),
		},
	}

	s.service.On("PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID).Return(nil, model.ErrBadRequest)
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.BadRequestError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID)
}

func (s *ApiSuit) TestPayOrderServiceError() {
	orderUUID := gofakeit.UUID()
	paymentMethod := orderV1.PaymentMethodPAYMENTMETHODCREDITCARD

	tc := &struct {
		req         *orderV1.PayOrderRequest
		params      orderV1.PayOrderParams
		expectedRes *orderV1.InternalServerError
	}{
		req: &orderV1.PayOrderRequest{
			PaymentMethod: paymentMethod,
		},
		params: orderV1.PayOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("order service: %s\n", "internal error"),
		},
	}

	s.service.On("PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID).Return(nil, errors.New("internal error"))
	resp, err := s.handler.PayOrder(s.ctx, tc.req, tc.params)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.InternalServerError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID)
}

func (s *ApiSuit) TestPayOrderDifferentPaymentMethods() {
	testCases := []struct {
		name           string
		paymentMethod  orderV1.PaymentMethod
	}{
		{
			name:          "credit card",
			paymentMethod: orderV1.PaymentMethodPAYMENTMETHODCREDITCARD,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			orderUUID := gofakeit.UUID()
			transactionUUID := gofakeit.UUID()

			req := &orderV1.PayOrderRequest{
				PaymentMethod: tc.paymentMethod,
			}
			params := orderV1.PayOrderParams{
				OrderUUID: orderUUID,
			}

			s.service.On("PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID).Return(&model.PayOrderResponse{
				TransactionUUID: transactionUUID,
			}, nil)

			resp, err := s.handler.PayOrder(s.ctx, req, params)
			s.Require().NoError(err)
			s.Require().NotNil(resp)
			
			payResp := resp.(*orderV1.PayOrderResponse)
			s.Require().Equal(payResp.TransactionUUID, transactionUUID)
			s.service.AssertCalled(s.T(), "PayOrder", s.ctx, mock.AnythingOfType("*model.PayOrderRequest"), orderUUID)
		})
	}
}

