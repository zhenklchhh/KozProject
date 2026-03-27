package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApi "github.com/zhenklchhh/KozProject/order/internal/api/order/v1"
	invClient "github.com/zhenklchhh/KozProject/order/internal/client/inventory/v1"
	payClient "github.com/zhenklchhh/KozProject/order/internal/client/payment/v1"
	"github.com/zhenklchhh/KozProject/order/internal/migrator"
	"github.com/zhenklchhh/KozProject/order/internal/repository/order"
	service "github.com/zhenklchhh/KozProject/order/internal/service/order"
	orderV1 "github.com/zhenklchhh/KozProject/shared/pkg/api/order/v1"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

const (
	httpPort           = "8082"
	readHeaderTimeout  = 5 * time.Second
	shutdownTimeout    = 10 * time.Second
	inventoryClientUri = "INVENTORY_CLIENT_URI"
	paymentClientUri   = "PAYMENT_CLIENT_URI"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}
	dbURI := os.Getenv("DB_URI")
	sqlDB, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer sqlDB.Close()
	migratorRunner := migrator.NewMigrator(sqlDB, os.Getenv("MIGRATIONS_DIR"))
	if err = migratorRunner.Up(); err != nil {
		log.Printf("Error applying migrations: %v\n", err)
		return
	}
	pool, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		log.Printf("Error connect to database: %v\n", err)
		return
	}
	repo := order.NewPostgresRepository(pool)
	connInv, err := grpc.NewClient(os.Getenv(inventoryClientUri), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create inventory client connection: %v", err)
		return
	}
	connPay, err := grpc.NewClient(os.Getenv(paymentClientUri), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create payment client connection: %v", err)
		return
	}
	inventoryClient := invClient.NewClient(inventoryV1.NewInventoryServiceClient(connInv))
	paymentClient := payClient.NewClient(paymentV1.NewPaymentServiceClient(connPay))
	svc := service.NewService(repo, paymentClient, inventoryClient)
	apiHandler := orderApi.NewApi(svc)
	apiServer, err := orderV1.NewServer(apiHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Recoverer)
	r.Mount("/", apiServer)
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}
}
