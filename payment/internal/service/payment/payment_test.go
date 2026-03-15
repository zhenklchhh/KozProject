package payment

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/zhenklchhh/KozProject/payment/internal/model"
)

type testCase struct {
	name         string
	req          *model.PayOrderRequest
	isValidReq   bool
	expectedUUID string
	repoError    error
}

func (s *ServiceSuit) TestPayOrderValidRequest() {
	tc := &testCase{
		req: &model.PayOrderRequest{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: model.PaymentMethodCard,
		},
		expectedUUID: gofakeit.UUID(),
	}
	s.paymentRepo.On("PayOrder", s.ctx, tc.req).Return(tc.expectedUUID, nil)
	response, err := s.service.PayOrder(s.ctx, tc.req)
	s.Require().Equal(tc.expectedUUID, response.TransactionUuid)
	s.Require().NoError(err)
}

func (s *ServiceSuit) TestPayOrderInvalidRequest() {
	tc := &testCase{
		req: &model.PayOrderRequest{
			OrderUuid:     "",
			UserUuid:      "",
			PaymentMethod: 404,
		},
	}
	s.paymentRepo.On("PayOrder", s.ctx, tc.req).Return("", model.ErrValidation)
	response, err := s.service.PayOrder(s.ctx, tc.req)
	s.Require().Error(err)
	s.Require().Nil(response)
}
