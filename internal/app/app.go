package app

import (
	"context"
	"log/slog"
	grpcapp "subpub/internal/app/grpc"
	"subpub/subpub"
)

type App struct {
	GRPCrv *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, pubsub subpub.SubPub, ctx context.Context) *App {

	app := grpcapp.NewApp(log, grpcPort, pubsub, ctx)
	return &App{
		GRPCrv: app,
	}
}
