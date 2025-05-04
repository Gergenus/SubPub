package app

import (
	"log/slog"
	grpcapp "subpub/internal/app/grpc"
	"subpub/subpub"
)

type App struct {
	GRPCrv *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, pubsub subpub.SubPub) *App {

	app := grpcapp.NewApp(log, grpcPort, pubsub)
	return &App{
		GRPCrv: app,
	}
}
