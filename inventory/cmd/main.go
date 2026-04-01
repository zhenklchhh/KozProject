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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	invPartApi "github.com/zhenklchhh/KozProject/inventory/internal/api/inventory/v1"
	"github.com/zhenklchhh/KozProject/inventory/internal/config"
	invPartRepo "github.com/zhenklchhh/KozProject/inventory/internal/repository/part"
	invPartService "github.com/zhenklchhh/KozProject/inventory/internal/service/part"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load("./deploy/env/.env")
	if err != nil {
		log.Printf("failed to load config: %v\n", err)
		return
	}

	lis, err := net.Listen("tcp", cfg.GRPC().Address())
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
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Mongo().URI()))
	if err != nil {
		log.Printf("failed to connect database: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	pingCtx, cancel := context.WithTimeout(ctx, cfg.HTTP().GetPingTimeout())
	defer cancel()
	err = client.Ping(pingCtx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}
	db := client.Database(cfg.Mongo().Database())
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
			cfg.GRPC().Address(),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}

		httpMux := http.NewServeMux()
		httpMux.Handle("/api/", mux)

		staticFS := http.FileServer(http.Dir(cfg.HTTP().StaticDir()))
		httpMux.Handle("/", staticFS)

		httpMux.HandleFunc("/api/inventory/v1/inventory.swagger.json", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("📄 Serving swagger file: %s", cfg.HTTP().GetSwaggerFile())
			w.Header().Set("Content-Type", "application/json")
			http.ServeFile(w, r, cfg.HTTP().GetSwaggerFile())
		})
		gwServer = &http.Server{
			Addr:              cfg.HTTP().Address(),
			Handler:           httpMux,
			ReadHeaderTimeout: cfg.HTTP().GetReadHeaderTimeout(),
		}

		log.Printf("📖 Swagger UI available at http://%s\n", cfg.GRPC().Address())
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
