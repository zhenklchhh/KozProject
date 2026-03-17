package api

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

var (
	parts = setUpParts()
)

func setUpParts() []*inventoryV1.Part {
	return []*inventoryV1.Part{
		&inventoryV1.Part{
			Uuid:  gofakeit.UUID(),
			Price: 4.3,
		},
		&inventoryV1.Part{
			Uuid:  gofakeit.UUID(),
			Price: 5.7,
		},
	}
}

func (s *ApiSuit) TestCreateOrderSuccess() {
	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
		expectedErr  error
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID: gofakeit.UUID(),
			PartUuids: []string{
				parts[0].Uuid,
				parts[1].Uuid,
			},
		},
		expectedResp: &orderV1.CreateOrderResponse{
			TotalPrice: 10.0,
		},
	}
	listPartsReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: tc.req.PartUuids,
		},
	}
	s.inventoryClient.On("ListParts", s.ctx, listPartsReq).Return(
		&inventoryV1.ListPartsResponse{Parts: parts}, nil,
	)
	var capturedOrder interface{}
	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
		capturedOrder = args.Get(1)
	})

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	createResp := response.(*orderV1.CreateOrderResponse)
	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().Equal(createResp.TotalPrice, tc.expectedResp.TotalPrice)
	s.Require().NotNil(createResp.OrderUUID)
	s.Require().NotNil(capturedOrder)
	s.Require().Equal(createResp.OrderUUID, capturedOrder.(*model.Order).OrderUUID)
	s.Require().Equal(capturedOrder.(*model.Order).UserUUID, tc.req.UserUUID)
	s.Require().Equal(capturedOrder.(*model.Order).PartUuids, tc.req.PartUuids)
	s.Require().Equal(capturedOrder.(*model.Order).Status, model.OrderStatusPendingPayment)
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, listPartsReq)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *ApiSuit) TestCreateOrderClientError() {
	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
		expectedErr  error
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID: gofakeit.UUID(),
			PartUuids: []string{
				parts[0].Uuid,
				parts[1].Uuid,
			},
		},
	}
	s.inventoryClient.On("ListParts", s.ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: tc.req.PartUuids,
		},
	}).Return(nil, errors.New("client error"))

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	s.Require().Error(err)
	s.Require().Empty(response)
	s.Require().Contains(err.Error(), "failed to get list parts from inventory service")
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, mock.Anything)
}

func (s *ApiSuit) TestCreateOrderMismatchNumberOfParts() {
	tc := &struct {
		req          *orderV1.CreateOrderRequest
		expectedResp *orderV1.CreateOrderResponse
		expectedErr  error
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID: gofakeit.UUID(),
			PartUuids: []string{
				parts[0].Uuid,
				parts[1].Uuid,
			},
		},
	}
	listPartsReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: tc.req.PartUuids,
		},
	}
	s.inventoryClient.On("ListParts", s.ctx, listPartsReq).Return(
		&inventoryV1.ListPartsResponse{Parts: []*inventoryV1.Part{parts[0]}}, nil,
	)
	response, _ := s.handler.CreateOrder(s.ctx, tc.req)
	resp := response.(*orderV1.BadRequestError)
	s.Require().Equal(resp.Code, 400)
	s.Require().Equal(resp.Message, "some parts don't exist in inventory service")
	s.service.AssertNotCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, listPartsReq)
}

func (s *ApiSuit) TestCreateOrderServiceSaveError() {
	tc := &struct {
		req *orderV1.CreateOrderRequest
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID: gofakeit.UUID(),
			PartUuids: []string{
				parts[0].Uuid,
				parts[1].Uuid,
			},
		},
	}
	listPartsReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: tc.req.PartUuids,
		},
	}
	s.inventoryClient.On("ListParts", s.ctx, listPartsReq).Return(
		&inventoryV1.ListPartsResponse{Parts: parts}, nil,
	)
	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("database error"))

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	s.Require().Error(err)
	s.Require().Empty(response)
	s.Require().Contains(err.Error(), "order service error")
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, listPartsReq)
	s.service.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *ApiSuit) TestCreateOrderEmptyPartUuids() {
	tc := &struct {
		req *orderV1.CreateOrderRequest
	}{
		req: &orderV1.CreateOrderRequest{
			UserUUID:  gofakeit.UUID(),
			PartUuids: []string{},
		},
	}
	listPartsReq := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartFilter{
			Uuids: tc.req.PartUuids,
		},
	}
	s.inventoryClient.On("ListParts", s.ctx, listPartsReq).Return(
		&inventoryV1.ListPartsResponse{Parts: []*inventoryV1.Part{}}, nil,
	)
	var capturedOrder interface{}
	s.service.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Run(func(args mock.Arguments) {
		capturedOrder = args.Get(1)
	})

	response, err := s.handler.CreateOrder(s.ctx, tc.req)
	createResp := response.(*orderV1.CreateOrderResponse)
	s.Require().NoError(err)
	s.Require().NotNil(createResp)
	s.Require().Equal(createResp.TotalPrice, 0.0)
	s.Require().NotNil(capturedOrder)
	s.Require().Equal(capturedOrder.(*model.Order).PartUuids, []string{})
}
