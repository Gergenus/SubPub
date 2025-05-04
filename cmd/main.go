package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"subpub/internal/app"
	"subpub/subpub"
	"syscall"
	"time"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	sub := subpub.NewSubPub()

	defer sub.Close(context.Background())
	application := app.NewApp(log, 1488, sub)

	go application.GRPCrv.Run()
	// Graceful shutdown
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application.GRPCrv.Stop()
}
