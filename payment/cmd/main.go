package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	paymentApi "github.com/zhenklchhh/KozProject/payment/internal/api/payment/v1"
	"github.com/zhenklchhh/KozProject/payment/internal/config"
	paymentRepo "github.com/zhenklchhh/KozProject/payment/internal/repository/payment"
	paymentService "github.com/zhenklchhh/KozProject/payment/internal/service/payment"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
)

func main() {
	cfg, err := config.Load("./deploy/env/.env")
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", cfg.GRPC().Address())
	if err != nil {
		log.Printf("failed listening: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()
	repo := paymentRepo.NewRepository()
	service := paymentService.NewService(repo)
	api := paymentApi.NewApi(service)
	paymentV1.RegisterPaymentServiceServer(s, api)
	reflection.Register(s)
	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", cfg.GRPC().Address())
		err = s.Serve(lis)
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

		err = paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			mux,
			cfg.GRPC().Address(),
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
		staticDir := projectRoot + "/payment/static"
		swaggerFile := projectRoot + "/api/payment/v1/payment.swagger.json"

		httpMux := http.NewServeMux()
		httpMux.Handle("/api/", mux)

		staticFS := http.FileServer(http.Dir(staticDir))
		httpMux.Handle("/", staticFS)

		httpMux.HandleFunc("/api/payment/v1/payment.swagger.json", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📄 Serving swagger file: %s", swaggerFile)
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, swaggerFile)
		})
		gwServer = &http.Server{
			Addr:              cfg.HTTP().Address(),
			Handler:           httpMux,
			ReadHeaderTimeout: cfg.HTTP().GetReadHeaderTimeout(),
		}

		log.Printf("📖 Swagger UI available at http://%s/\n", cfg.HTTP().Address())
		err = gwServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to listen server: %v\n", err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if gwServer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP().GetReadHeaderTimeout())
		defer cancel()
		if err = gwServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}
	s.GracefulStop()
}
