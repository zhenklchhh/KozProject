package app

import (
	"context"
	"sync"

	api "github.com/zhenklchhh/KozProject/payment/internal/api/payment/v1"
	"github.com/zhenklchhh/KozProject/payment/internal/repository"
	paymentRepo "github.com/zhenklchhh/KozProject/payment/internal/repository/payment"
	"github.com/zhenklchhh/KozProject/payment/internal/service"
	paymentServiceImpl "github.com/zhenklchhh/KozProject/payment/internal/service/payment"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1Api      paymentV1.PaymentServiceServer
	paymentService    service.PaymentService
	paymentRepository repository.PaymentRepository

	onceApi        sync.Once
	onceService    sync.Once
	onceRepository sync.Once
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) PaymentV1Api(ctx context.Context) paymentV1.PaymentServiceServer {
	c.onceApi.Do(func() {
		c.paymentV1Api = api.NewApi(c.PaymentService(ctx))
	})
	return c.paymentV1Api
}

func (c *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	c.onceService.Do(func() {
		c.paymentService = paymentServiceImpl.NewService(c.PaymentRepository(ctx))
	})
	return c.paymentService
}

func (c *diContainer) PaymentRepository(ctx context.Context) repository.PaymentRepository {
	c.onceRepository.Do(func() {
		c.paymentRepository = paymentRepo.NewRepository()
	})
	return c.paymentRepository
}