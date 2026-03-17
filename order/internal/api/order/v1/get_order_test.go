package api

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

func (s *ApiSuit) TestGetOrderSuccess() {
	OrderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}

	tc := &struct {
		req           *orderV1.GetOrderParams
		order         *model.Order
		expectedOrder *orderV1.Order
		expectedErr   error
	}{
		req: &orderV1.GetOrderParams{
			OrderUUID: OrderUUID,
		},
		order: &model.Order{
			OrderUUID:  OrderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     model.OrderStatusPaid,
		},
		expectedOrder: &orderV1.Order{
			OrderUUID:  OrderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUuids,
			TotalPrice: 15.99,
			Status:     orderV1.OrderStatusPAID,
			TransactionUUID: orderV1.OptNilString{
				Value: "",
				Set:   true,
				Null:  true,
			},
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(tc.order, nil)
	response, err := s.handler.GetOrder(s.ctx, *tc.req)
	s.Require().NoError(err)
	s.Require().Equal(tc.expectedOrder, response)
	s.service.AssertCalled(s.T(), "Get", s.ctx, OrderUUID)
}

func (s *ApiSuit) TestGetOrderNotFound() {
	OrderUUID := gofakeit.UUID()
	tc := &struct {
		req           *orderV1.GetOrderParams
		order         *model.Order
		expectedOrder *orderV1.Order
		expectedResp  orderV1.GetOrderRes
	}{
		req: &orderV1.GetOrderParams{
			OrderUUID: OrderUUID,
		},
		expectedResp: &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order %s not found: %s", OrderUUID, "not found"),
		},
	}
	s.service.On("Get", s.ctx, tc.req.OrderUUID).Return(nil, errors.New("not found"))
	response, err := s.handler.GetOrder(s.ctx, *tc.req)
	s.Require().NoError(err)
	s.Require().Equal(tc.expectedResp, response)
	s.service.AssertCalled(s.T(), "Get", s.ctx, OrderUUID)
}
