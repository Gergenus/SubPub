package app

import (
	"context"
	"log/slog"
	grpcapp "subpub/internal/app/grpc"
	"subpub/internal/services/eventbus"
)

type App struct {
	GRPCrv *grpcapp.App
}

func NewApp(ctx context.Context, log *slog.Logger, grpcPort int, pubsub *eventbus.Eventbus) *App {

	app := grpcapp.NewApp(ctx, log, grpcPort, pubsub)
	return &App{
		GRPCrv: app,
	}
}
