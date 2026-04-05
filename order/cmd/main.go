package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/zhenklchhh/KozProject/order/internal/app"
	"github.com/zhenklchhh/KozProject/order/internal/config"
	"github.com/zhenklchhh/KozProject/platform/pkg/closer"
)

const (
	path            = "./deploy/compose/order/.env"
	shutdownTimeout = 5 * time.Second
)

func main() {
	err := config.Load(path)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	app, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatalf("Error creating application: %v\n", err)
	}

	err = app.Run(appCtx)
	if err != nil {
		log.Fatalf("Error to run application: %v\n", err)
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := closer.CloseAll(ctx); err != nil {
		log.Fatalf("Error to close work of application: %v\n", err)
	}
}
