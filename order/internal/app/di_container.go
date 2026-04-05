package app

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	api "github.com/zhenklchhh/KozProject/order/internal/api/order/v1"
	invClient "github.com/zhenklchhh/KozProject/order/internal/client/inventory/v1"
	payClient "github.com/zhenklchhh/KozProject/order/internal/client/payment/v1"
	"github.com/zhenklchhh/KozProject/order/internal/config"
	"github.com/zhenklchhh/KozProject/order/internal/migrator"
	"github.com/zhenklchhh/KozProject/order/internal/repository"
	orderRepository "github.com/zhenklchhh/KozProject/order/internal/repository/order"
	"github.com/zhenklchhh/KozProject/order/internal/service"
	orderService "github.com/zhenklchhh/KozProject/order/internal/service/order"
	"github.com/zhenklchhh/KozProject/order/internal/transaction"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
 */
type diContainer struct {
	orderV1Api      orderV1.Server
	orderService    service.OrderService
	orderRepository repository.OrderRepository
	txManager       transaction.TransactionManager
	migratorRunner  *migrator.Migrator
	pool            *pgxpool.Pool

	onceApi            sync.Once
	oncePool           sync.Once
	onceService        sync.Once
	onceRepository     sync.Once
	onceTxManager      sync.Once
	onceMigratorRunner sync.Once
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) OrderV1Api(ctx context.Context) orderV1.Server {
	c.onceApi.Do(func() {
		apiServer, err := orderV1.NewServer(api.NewApi(c.OrderService(ctx)))
		if err != nil {
			panic(fmt.Sprintf("error initializing order api server: %w", err))
		}
		c.orderV1Api = *apiServer
	})
	return c.orderV1Api
}

func (c *diContainer) PgxPool(ctx context.Context) *pgxpool.Pool {
	c.oncePool.Do(func() {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres().URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %w\n", err))
		}
		c.pool = pool
		closer.AddNamed("Pgx Pool", func(context.Context) error {
			pool.Close()
			return nil
		})
	})
	return c.pool
}

func (c *diContainer) OrderService(ctx context.Context) service.OrderService {
	c.onceService.Do(func() {
		connInv, err := grpc.NewClient(config.AppConfig().InventoryClient().URI(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("failed to create inventory client connection: %v", err))
		}
		inventoryClient := invClient.NewClient(inventoryV1.NewInventoryServiceClient(connInv))
		connPay, err := grpc.NewClient(config.AppConfig().PaymentClient().URI(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("failed to create payment client connection: %v", err))
		}
		paymentClient := payClient.NewClient(paymentV1.NewPaymentServiceClient(connPay))
		c.orderService = orderService.NewService(
			c.OrderRepository(ctx),
			c.TxManager(ctx),
			paymentClient,
			inventoryClient,
		)
	})
	return c.orderService
}
func (c *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	c.onceRepository.Do(func() {
		c.orderRepository = orderRepository.NewPostgresRepository(c.PgxPool(ctx))
	})
	return c.orderRepository
}

func (c *diContainer) TxManager(ctx context.Context) transaction.TransactionManager {
	c.onceTxManager.Do(func() {
		c.txManager = transaction.NewManager(c.PgxPool(ctx))
	})
	return c.txManager
}

func (c *diContainer) MigratorRunner(ctx context.Context) *migrator.Migrator {
	c.onceApi.Do(func() {
		sqlDB, err := sql.Open("pgx", config.AppConfig().Postgres().URI())
		if err != nil {
			panic(fmt.Sprintf("Error connecting to database: %v\n", err))
		}
		closer.AddNamed("Sql Connection", func(ctx context.Context) error {
			sqlDB.Close()
			return nil
		})
		c.migratorRunner = migrator.NewMigrator(sqlDB, config.AppConfig().Migrations().Dir())
	})
	return c.migratorRunner
}
