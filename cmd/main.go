package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"subpub/internal/app"
	"subpub/internal/config"
	"subpub/internal/services/eventbus"
	"subpub/subpub"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	sub := subpub.NewSubPub()
	subService := eventbus.NewEventBus(log, sub)
	log.Info("starting application", slog.String("env", cfg.Env), slog.Int("port", cfg.GRPC.Port), slog.Any("cfg", cfg))

	upper_ctx, upper_cancel := context.WithCancel(context.Background())
	application := app.NewApp(upper_ctx, log, cfg.GRPC.Port, subService)

	go application.GRPCrv.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	upper_cancel()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	defer cancel()

	application.GRPCrv.Stop()
	sub.Close(ctx)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
