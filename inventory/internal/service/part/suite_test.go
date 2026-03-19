package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMock "github.com/zhenklchhh/KozProject/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx     context.Context
	invRepo *repoMock.InventoryRepository
	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.invRepo = repoMock.NewInventoryRepository(s.T())
	s.service = NewService(s.invRepo)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceInteg(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
