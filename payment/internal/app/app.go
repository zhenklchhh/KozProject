package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhenklchhh/KozProject/payment/internal/config"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
	"github.com/zhenklchhh/KozProject/platform/pkg/logger"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	diContainer       *diContainer
	grpcServer        *grpc.Server
	httpGatewayServer *http.Server
	listener          net.Listener
}

func NewApp(ctx context.Context) (*App, error) {
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
	deps := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initHttpGatewayServer,
	}
	for _, f := range deps {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(ctx context.Context) error {
	a.diContainer = NewDIContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	logger.Init(config.AppConfig().Logger().Level(), config.AppConfig().Logger().AsJson())
	return nil
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initListener(ctx context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC().Address())
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Error listen address: %v", config.AppConfig().GRPC().Address()), zap.Error(err))
		return err
	}
	closer.AddNamed("Payment TCP listener", func(ctx context.Context) error {
		listener.Close()
		return nil
	})
	a.listener = listener
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	s := grpc.NewServer()
	paymentV1.RegisterPaymentServiceServer(s, a.diContainer.PaymentV1Api(ctx))
	reflection.Register(s)
	closer.AddNamed("Payment gRPC server", func(ctx context.Context) error {
		s.GracefulStop()
		return nil
	})
	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		logger.Error(ctx, "Error to start payment grpc server", zap.Error(err))
		return err
	}
	return nil
}

func (a *App) initHttpGatewayServer(ctx context.Context) error {
	gwCtx, gwCancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
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

	a.httpGatewayServer = &http.Server{
		Addr:              config.AppConfig().HTTP().Address(),
		Handler:           httpMux,
		ReadHeaderTimeout: config.AppConfig().HTTP().GetReadHeaderTimeout(),
	}

	closer.AddNamed("Payment Http Gateway", func(ctx context.Context) error {
		gwCancel()
		return a.httpGatewayServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHttpGatewayServer(ctx context.Context) error {
	err := a.httpGatewayServer.Serve(a.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen server: %v\n", err)
	}
	return nil
}
