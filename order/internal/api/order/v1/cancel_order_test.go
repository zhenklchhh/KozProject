package api

import (
	"errors"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

func (s *ApiSuit) TestCancelOrderSuccess() {
	orderUUID := gofakeit.UUID()

	tc := &struct {
		req         orderV1.CancelOrderParams
		expectedRes orderV1.CancelOrderRes
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.CancelOrderNoContent{},
	}

	s.service.On("CancelOrder", s.ctx, tc.req.OrderUUID).Return(nil)
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp, tc.expectedRes)
	s.service.AssertCalled(s.T(), "CancelOrder", s.ctx, tc.req.OrderUUID)
}

func (s *ApiSuit) TestCancelOrderNotFound() {
	orderUUID := gofakeit.UUID()

	tc := &struct {
		req         orderV1.CancelOrderParams
		expectedRes *orderV1.NotFoundError
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.NotFoundError{
			Code:    404,
			Message: fmt.Sprintf("order service: %s\n", "not found"),
		},
	}

	s.service.On("CancelOrder", s.ctx, tc.req.OrderUUID).Return(model.ErrNotFound)
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.NotFoundError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "CancelOrder", s.ctx, tc.req.OrderUUID)
}

func (s *ApiSuit) TestCancelOrderInvalidStatus() {
	orderUUID := gofakeit.UUID()

	tc := &struct {
		req         orderV1.CancelOrderParams
		expectedRes *orderV1.ConflictError
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.ConflictError{
			Code:    409,
			Message: fmt.Sprintf("order service: %s\n", "conflict"),
		},
	}

	s.service.On("CancelOrder", s.ctx, tc.req.OrderUUID).Return(model.ErrConflict)
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.ConflictError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "CancelOrder", s.ctx, tc.req.OrderUUID)
}

func (s *ApiSuit) TestCancelOrderServiceError() {
	orderUUID := gofakeit.UUID()

	tc := &struct {
		req         orderV1.CancelOrderParams
		expectedRes *orderV1.InternalServerError
	}{
		req: orderV1.CancelOrderParams{
			OrderUUID: orderUUID,
		},
		expectedRes: &orderV1.InternalServerError{
			Code:    500,
			Message: fmt.Sprintf("order service: %s\n", "internal error"),
		},
	}

	s.service.On("CancelOrder", s.ctx, tc.req.OrderUUID).Return(errors.New("internal error"))
	resp, err := s.handler.CancelOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().Equal(resp.(*orderV1.InternalServerError), tc.expectedRes)
	s.service.AssertCalled(s.T(), "CancelOrder", s.ctx, tc.req.OrderUUID)
}
