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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	invPartRepo "github.com/zhenklchhh/KozProject/inventory/internal/repository/part"
	invPartService "github.com/zhenklchhh/KozProject/inventory/internal/service/part"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort          = 50051
	httpPort          = 8082
	readHeaderTimeout = 10 * time.Second
)

func main() {
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
	repo := invPartRepo.NewRepository()
	service := invPartService.NewService(repo)
	inventoryV1.RegisterInventoryServiceServer(s, service)
	reflection.Register(s)
	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve server: %v\n", err)
			return
		}
	}()
	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mux := runtime.NewServeMux()

		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
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
		if err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Printf("failed to listen server: %v\n", err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.GracefulStop()
}
