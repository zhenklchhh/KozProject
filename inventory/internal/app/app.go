package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhenklchhh/KozProject/inventory/internal/config"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
	"github.com/zhenklchhh/KozProject/platform/pkg/grpc/health"
	"github.com/zhenklchhh/KozProject/platform/pkg/logger"
	inventoryv1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer       *diContainer
	grpcServer        *grpc.Server
	gatewayHttpServer *http.Server
	listener          net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	numRunFuncs := 2
	errCh := make(chan error, numRunFuncs)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.runGRPCServer(ctx)
		if err != nil {
			errCh <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.runHttpGatewayServer(ctx)
		if err != nil {
			errCh <- err
		}
	}()
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()
	select {
	case err := <-errCh:
		return err
	case <-doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initHttpGatewayServer,
	}
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(config.AppConfig().Logger().Level(), config.AppConfig().Logger().AsJson())
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC().Address())
	if err != nil {
		return err
	}
	closer.AddNamed("Inventory TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}
		return nil
	})
	a.listener = listener
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)
	health.RegisterServer(a.grpcServer)
	inventoryv1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1Api(ctx))
	return nil
}

func (a *App) initHttpGatewayServer(ctx context.Context) error {
	gwCtx, gwCancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := inventoryv1.RegisterInventoryServiceHandlerFromEndpoint(
		gwCtx,
		mux,
		config.AppConfig().GRPC().Address(),
		opts,
	)
	if err != nil {
		return fmt.Errorf("Failed to register gateway: %w\n", err)
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/api/", mux)

	staticFS := http.FileServer(http.Dir(config.AppConfig().HTTP().StaticDir()))
	httpMux.Handle("/", staticFS)

	httpMux.HandleFunc("/api/inventory/v1/inventory.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📄 Serving swagger file: %s", config.AppConfig().HTTP().GetSwaggerFile())
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, config.AppConfig().HTTP().GetSwaggerFile())
	})
	a.gatewayHttpServer = &http.Server{
		Addr:              config.AppConfig().HTTP().Address(),
		Handler:           httpMux,
		ReadHeaderTimeout: config.AppConfig().HTTP().GetReadHeaderTimeout(),
	}

	closer.AddNamed("Inventory Http Gateway", func(ctx context.Context) error {
		gwCancel()
		return a.gatewayHttpServer.Shutdown(ctx)
	})

	logger.Info(ctx, fmt.Sprintf("📖 Swagger UI available at http://%s\n", config.AppConfig().GRPC().Address()))
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server listening on %s", config.AppConfig().GRPC().Address()))
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runHttpGatewayServer(ctx context.Context) error {
	err := a.gatewayHttpServer.Serve(a.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen server: %v\n", err)
	}
	return nil
}
