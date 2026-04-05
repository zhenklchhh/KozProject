package app

import (
	"context"
	"fmt"
	"sync"

	invV1Api "github.com/zhenklchhh/KozProject/inventory/internal/api/inventory/v1"
	"github.com/zhenklchhh/KozProject/inventory/internal/config"
	"github.com/zhenklchhh/KozProject/inventory/internal/repository"
	invRepository "github.com/zhenklchhh/KozProject/inventory/internal/repository/part"
	"github.com/zhenklchhh/KozProject/inventory/internal/service"
	invService "github.com/zhenklchhh/KozProject/inventory/internal/service/part"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository
	mongoDBClient       *mongo.Client
	mongoDBHandler      *mongo.Database

	onceApi        sync.Once
	onceService    sync.Once
	onceRepository sync.Once
	onceClient     sync.Once
	onceHandler    sync.Once
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) InventoryV1Api(ctx context.Context) inventoryV1.InventoryServiceServer {
	c.onceApi.Do(func() {
		c.inventoryV1API = invV1Api.NewAPI(c.InventoryService(ctx))
	})
	return c.inventoryV1API
}

func (c *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	c.onceService.Do(func() {
		c.inventoryService = invService.NewService(c.InventoryRepository(ctx))
	})
	return c.inventoryService
}

func (c *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	c.onceRepository.Do(func() {
		repo, err := invRepository.NewMongoRepository(c.MongoDBHandler(ctx))
		if err != nil {
			panic(fmt.Sprintf("Error creating mongo repository: %s\n", err.Error()))
		}
		c.inventoryRepository = repo
	})
	return c.inventoryRepository
}

func (c *diContainer) MongoDBHandler(ctx context.Context) *mongo.Database {
	c.onceHandler.Do(func() {
		c.mongoDBHandler = c.MongoDBClient(ctx).Database(config.AppConfig().Mongo().Database())
	})
	return c.mongoDBHandler
}

func (c *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	c.onceClient.Do(func() {
		client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig().Mongo().URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping mongoDB server: %s\n", err.Error()))
		}
		closer.AddNamed("MongoDB Client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})
		c.mongoDBClient = client
	})
	return c.mongoDBClient
}
