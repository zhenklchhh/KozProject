package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhenklchhh/KozProject/inventory/internal/app"
	"github.com/zhenklchhh/KozProject/inventory/internal/config"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
	"github.com/zhenklchhh/KozProject/platform/pkg/logger"
	"go.uber.org/zap"
)

const (
	path = "./deploy/compose/inventory/.env"
	shutdownTimeout = 5 * time.Second
)

func main() {
	err := config.Load(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v\n", err))
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(context.Background())
	if err != nil {
		logger.Error(appCtx, "failed to create application", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "error running app", zap.Error(err))
		return	
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "error to close work of application", zap.Error(err))
	}
}
