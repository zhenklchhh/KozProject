package order

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/zhenklchhh/KozProject/order/internal/client"
	"github.com/zhenklchhh/KozProject/order/internal/model"
	repoMock "github.com/zhenklchhh/KozProject/order/internal/repository/mocks"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

type OrderServiceSuit struct {
	suite.Suite

	ctx             context.Context
	repo            *repoMock.OrderRepository
	paymentClient   *client.MockPaymentClient
	inventoryClient *client.MockInventoryClient
	service         *service
}

func (s *OrderServiceSuit) SetupTest() {
	s.ctx = context.Background()
	s.repo = repoMock.NewOrderRepository(s.T())
	s.paymentClient = &client.MockPaymentClient{}
	s.inventoryClient = &client.MockInventoryClient{}
	s.service = NewService(s.repo, s.paymentClient, s.inventoryClient)
}

func (s *OrderServiceSuit) TearDownTest() {
}

func TestOrderServiceIntegration(t *testing.T) {
	suite.Run(t, new(OrderServiceSuit))
}

// Create Order Tests

func (s *OrderServiceSuit) TestCreateOrderSuccess() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	parts := []*inventoryV1.Part{
		{Uuid: partUuids[0], Price: 5.0},
		{Uuid: partUuids[1], Price: 7.5},
	}

	req := &model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUuids,
	}

	s.inventoryClient.On("ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids}).Return(parts, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	resp, err := s.service.Create(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(resp.TotalPrice, 12.5)
	s.Require().NotEmpty(resp.OrderUUID)
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids})
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *OrderServiceSuit) TestCreateOrderInventoryClientError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}

	req := &model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUuids,
	}

	s.inventoryClient.On("ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids}).Return(nil, errors.New("inventoryV1 error"))

	resp, err := s.service.Create(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().Contains(err.Error(), "inventory client error")
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids})
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestCreateOrderRepositoryError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	parts := []*inventoryV1.Part{
		{Uuid: partUuids[0], Price: 3.0},
		{Uuid: partUuids[1], Price: 4.0},
	}

	req := &model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUuids,
	}

	s.inventoryClient.On("ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids}).Return(parts, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("repo error"))

	resp, err := s.service.Create(s.ctx, req)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().Contains(err.Error(), "error creating order")
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids})
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *OrderServiceSuit) TestCreateOrderEmptyParts() {
	userUUID := gofakeit.UUID()
	partUuids := []string{}

	req := &model.CreateOrderRequest{
		UserUUID:  userUUID,
		PartUuids: partUuids,
	}

	s.inventoryClient.On("ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids}).Return([]*inventoryV1.Part{}, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	resp, err := s.service.Create(s.ctx, req)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(resp.TotalPrice, 0.0)
	s.Require().NotEmpty(resp.OrderUUID)
	s.inventoryClient.AssertCalled(s.T(), "ListParts", s.ctx, &inventoryV1.PartFilter{Uuids: partUuids})
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

// Get Order Tests

func (s *OrderServiceSuit) TestGetOrderSuccess() {
	orderUUID := gofakeit.UUID()
	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 10.0,
		Status:     model.OrderStatusPendingPayment,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)

	resp, err := s.service.Get(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(resp, order)
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
}

func (s *OrderServiceSuit) TestGetOrderNotFound() {
	orderUUID := gofakeit.UUID()

	s.repo.On("Get", s.ctx, orderUUID).Return(nil, errors.New("not found"))

	resp, err := s.service.Get(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
}

// Pay Order Tests

func (s *OrderServiceSuit) TestPayOrderSuccess() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPendingPayment,
	}

	req := &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, mock.Anything).Return(&paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	resp, err := s.service.PayOrder(s.ctx, req, orderUUID)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(resp.TransactionUUID, transactionUUID)
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.paymentClient.AssertCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *OrderServiceSuit) TestPayOrderNotFound() {
	orderUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard

	req := &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(nil, errors.New("not found"))

	resp, err := s.service.PayOrder(s.ctx, req, orderUUID)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().True(errors.Is(err, model.ErrNotFound))
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestPayOrderInvalidStatus() {
	orderUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPaid,
	}

	req := &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)

	resp, err := s.service.PayOrder(s.ctx, req, orderUUID)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().True(errors.Is(err, model.ErrConflict))
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestPayOrderInvalidPaymentMethod() {
	orderUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethod("INVALID_METHOD")

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPendingPayment,
	}

	req := &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)

	resp, err := s.service.PayOrder(s.ctx, req, orderUUID)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().True(errors.Is(err, model.ErrBadRequest))
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestPayOrderPaymentClientError() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPendingPayment,
	}

	req := &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, mock.Anything).Return(nil, errors.New("payment error"))

	resp, err := s.service.PayOrder(s.ctx, req, orderUUID)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().Contains(err.Error(), "payment client error")
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.paymentClient.AssertCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestPayOrderUpdateError() {
	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   userUUID,
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPendingPayment,
	}

	req := &model.PayOrderRequest{
		PaymentMethod: paymentMethod,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.On("PayOrder", s.ctx, mock.Anything).Return(&paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("update error"))

	resp, err := s.service.PayOrder(s.ctx, req, orderUUID)

	s.Require().Error(err)
	s.Require().Nil(resp)
	s.Require().Contains(err.Error(), "order service: failed to update order")
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.paymentClient.AssertCalled(s.T(), "PayOrder", s.ctx, mock.Anything)
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

// Cancel Order Tests

func (s *OrderServiceSuit) TestCancelOrderSuccess() {
	orderUUID := gofakeit.UUID()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPendingPayment,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(nil)

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}

func (s *OrderServiceSuit) TestCancelOrderNotFound() {
	orderUUID := gofakeit.UUID()

	s.repo.On("Get", s.ctx, orderUUID).Return(nil, errors.New("not found"))

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "order "+orderUUID+" not found")
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestCancelOrderInvalidStatus() {
	orderUUID := gofakeit.UUID()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPaid,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "can't be cancelled with status")
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.repo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
}

func (s *OrderServiceSuit) TestCancelOrderUpdateError() {
	orderUUID := gofakeit.UUID()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   gofakeit.UUID(),
		PartUuids:  []string{gofakeit.UUID()},
		TotalPrice: 15.99,
		Status:     model.OrderStatusPendingPayment,
	}

	s.repo.On("Get", s.ctx, orderUUID).Return(order, nil)
	s.repo.On("Create", s.ctx, mock.AnythingOfType("*model.Order")).Return(errors.New("update error"))

	err := s.service.CancelOrder(s.ctx, orderUUID)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed updating order")
	s.repo.AssertCalled(s.T(), "Get", s.ctx, orderUUID)
	s.repo.AssertCalled(s.T(), "Create", s.ctx, mock.AnythingOfType("*model.Order"))
}
