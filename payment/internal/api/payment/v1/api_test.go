package api

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zhenklchhh/KozProject/payment/internal/model"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

type testCase struct {
	req          *paymentV1.PayOrderRequest
	expectedUUID string
	expectedErr  error
}

func (s *ApiSuit) TestPayOrderValidRequest() {
	tc := &testCase{
		req: &paymentV1.PayOrderRequest{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
		},
		expectedUUID: gofakeit.UUID(),
	}
	convertedReq := &model.PayOrderRequest{
		OrderUuid:     tc.req.OrderUuid,
		UserUuid:      tc.req.UserUuid,
		PaymentMethod: model.PaymentMethodCard,
	}
	s.service.On("PayOrder", s.ctx, convertedReq).Return(&model.PayOrderResponse{TransactionUuid: tc.expectedUUID}, nil)
	response, err := s.handler.PayOrder(s.ctx, tc.req)
	s.Require().Equal(tc.expectedUUID, response.TransactionUuid)
	s.Require().NoError(err)
}

func (s *ApiSuit) TestPayOrderInvalidRequest() {
	tc := &testCase{
		req: &paymentV1.PayOrderRequest{
			OrderUuid:     "",
			UserUuid:      "",
			PaymentMethod: 404,
		},
		expectedErr: status.Errorf(codes.InvalidArgument, "payment service: validation error"),
	}
	response, err := s.handler.PayOrder(s.ctx, tc.req)
	s.Require().Error(err)
	s.Require().Equal(tc.expectedErr, err)
	s.Require().Nil(response)
}
