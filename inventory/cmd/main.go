package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	invPartApi "github.com/zhenklchhh/KozProject/inventory/internal/api/inventory/v1"
	invPartRepo "github.com/zhenklchhh/KozProject/inventory/internal/repository/part"
	invPartService "github.com/zhenklchhh/KozProject/inventory/internal/service/part"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort          = 50051
	httpPort          = 8082
	readHeaderTimeout = 10 * time.Second
	pingTimeout       = 5 * time.Second
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load env file: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()
	s := grpc.NewServer()

	dbURI := os.Getenv("MONGO_DB_URI")
	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect database: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	err = client.Ping(pingCtx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}
	db := client.Database("inventory-db")
	repo, err := invPartRepo.NewMongoRepository(db)
	if err != nil {
		log.Printf("failed to create inventory repository: %v\n", err)
		return
	}
	service := invPartService.NewService(repo)
	apiHandler := invPartApi.NewAPI(service)
	inventoryV1.RegisterInventoryServiceServer(s, apiHandler)
	reflection.Register(s)
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve server: %v\n", err)
			return
		}
	}()
	var gwServer *http.Server
	gwCtx, gwCancel := context.WithCancel(ctx)
	go func() {
		defer gwCancel()
		mux := runtime.NewServeMux()

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			gwCtx,
			mux,
			fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		projectRoot := os.Getenv("PROJECT_ROOT")
		if projectRoot == "" {
			projectRoot = "../.."
		}
		staticDir := projectRoot + "/inventory/static"
		swaggerFile := projectRoot + "/api/inventory/v1/inventory.swagger.json"

		httpMux := http.NewServeMux()
		httpMux.Handle("/api/", mux)

		staticFS := http.FileServer(http.Dir(staticDir))
		httpMux.Handle("/", staticFS)

		httpMux.HandleFunc("/api/inventory/v1/inventory.swagger.json", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📄 Serving swagger file: %s", swaggerFile)
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, swaggerFile)
		})
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: readHeaderTimeout,
		}

		log.Printf("📖 Swagger UI available at http://localhost:%d/\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to listen server: %v\n", err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.GracefulStop()
}
