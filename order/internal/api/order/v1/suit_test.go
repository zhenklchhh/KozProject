package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	serviceMock "github.com/zhenklchhh/KozProject/order/internal/service/mocks"
	invClientMock "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1/mocks"
	paymentClientMock "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1/mocks"
)

type ApiSuit struct {
	suite.Suite

	ctx             context.Context
	service         *serviceMock.OrderService
	inventoryClient *invClientMock.InventoryServiceClient
	paymentClient   *paymentClientMock.PaymentServiceClient
	handler         *api
}

func (s *ApiSuit) SetupTest() {
	s.ctx = context.Background()
	s.service = serviceMock.NewOrderService(s.T())
	s.inventoryClient = invClientMock.NewInventoryServiceClient(s.T())
	s.paymentClient = paymentClientMock.NewPaymentServiceClient(s.T())
	s.handler = NewApi(s.service, s.inventoryClient, s.paymentClient)
}

func (s *ApiSuit) TearDownTest() {
}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuit))
}
