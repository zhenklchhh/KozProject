package api

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
)

func (s *ApiSuite) TestCreateOrderSuccess() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	orderUUID := gofakeit.UUID()
	totalPrice := 10.0

	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		},
		expectedResp: &orderV1.CreateOrderResponse{
			OrderUUID:  orderUUID,
			TotalPrice: totalPrice,
		},
	}

	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest")).Return(&model.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil)

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	createResp := response.(*orderV1.CreateOrderResponse)

	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().Equal(createResp.TotalPrice, tc.expectedResp.TotalPrice)
	s.Require().Equal(createResp.OrderUUID, tc.expectedResp.OrderUUID)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest"))
}

func (s *ApiSuite) TestCreateOrderServiceError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}

	tc := &struct {
		req *orderV1.CreateOrderRequest
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		},
	}

	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest")).Return(nil, errors.New("service error"))

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	s.Require().NoError(err)
	s.Require().NotNil(response)

	errorResp := response.(*orderV1.InternalServerError)
	s.Require().Equal(errorResp.Code, 500)
	s.Require().Contains(errorResp.Message, "order service: service error")
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest"))
}

func (s *ApiSuite) TestCreateOrderEmptyPartUuids() {
	userUUID := gofakeit.UUID()
	orderUUID := gofakeit.UUID()

	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []string{},
		},
		expectedResp: &orderV1.CreateOrderResponse{
			OrderUUID:  orderUUID,
			TotalPrice: 0.0,
		},
	}

	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest")).Return(&model.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: 0.0,
	}, nil)

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	createResp := response.(*orderV1.CreateOrderResponse)

	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().Equal(createResp.TotalPrice, 0.0)
	s.Require().Equal(createResp.OrderUUID, orderUUID)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest"))
}

func (s *ApiSuite) TestCreateOrderSinglePart() {
	userUUID := gofakeit.UUID()
	partUuid := gofakeit.UUID()
	orderUUID := gofakeit.UUID()
	totalPrice := 5.99

	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: []string{partUuid},
		},
		expectedResp: &orderV1.CreateOrderResponse{
			OrderUUID:  orderUUID,
			TotalPrice: totalPrice,
		},
	}

	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest")).Return(&model.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil)

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	createResp := response.(*orderV1.CreateOrderResponse)

	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().Equal(createResp.TotalPrice, totalPrice)
	s.Require().Equal(createResp.OrderUUID, orderUUID)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest"))
}

func (s *ApiSuite) TestCreateOrderMultipleParts() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}
	orderUUID := gofakeit.UUID()
	totalPrice := 25.99

	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUuids: partUuids,
		},
		expectedResp: &orderV1.CreateOrderResponse{
			OrderUUID:  orderUUID,
			TotalPrice: totalPrice,
		},
	}

	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest")).Return(&model.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil)

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	createResp := response.(*orderV1.CreateOrderResponse)

	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().Equal(createResp.TotalPrice, totalPrice)
	s.Require().Equal(createResp.OrderUUID, orderUUID)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.CreateOrderRequest"))
}
