package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zhenklchhh/KozProject/order/internal/config"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
	"github.com/zhenklchhh/KozProject/platform/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	handler     http.Handler
	listener    net.Listener
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
	err := a.runHttpServer(ctx)
	if err != nil {
		logger.Error(ctx, "error running order http server\n", zap.Error(err))
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	deps := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initHandler,
		a.initHttpServer,
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
	a.diContainer = NewDiContainer()
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
	lis, err := net.Listen("tcp", config.AppConfig().HTTP().Address())
	if err != nil {
		return err
	}
	closer.AddNamed("Order TCP listener", func(ctx context.Context) error {
		lerr := lis.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}
		return nil
	})
	a.listener = lis
	return nil
}

func (a *App) initHandler(ctx context.Context) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Recoverer)
	r.Mount("/", &a.diContainer.orderV1Api)
	a.handler = r
	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().HTTP().Address(),
		Handler:           a.handler,
		ReadHeaderTimeout: config.AppConfig().HTTP().GetReadHeaderTimeout(),
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	closer.AddNamed("Order Http Server", func(ctx context.Context) error {
		cancel()
		return a.httpServer.Shutdown(shutdownCtx)
	})
	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP-сервер запущен по адресу %s\n", config.AppConfig().HTTP().Address()))
	err := a.httpServer.Serve(a.listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❌ Ошибка запуска сервера: %v\n", zap.Error(err))
	}
	return err
}
