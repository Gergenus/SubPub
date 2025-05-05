package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"subpub/subpub"
)

type Eventbus struct {
	log *slog.Logger
	sp  subpub.SubPub
}

type loggedSubscription struct {
	sub     subpub.Subscription
	log     *slog.Logger
	subject string
}

func (l *loggedSubscription) Unsubscribe() {
	const op = "eventbus.Unsubscribe"
	log := l.log.With(slog.String("op", op))
	log.Info("unsubscribing", slog.String("subject", l.subject))
	l.sub.Unsubscribe()
	log.Info("unsubscribed successfully", slog.String("subject", l.subject))
}

func (e *Eventbus) Subscribe(subject string, cb subpub.MessageHandler) (subpub.Subscription, error) {
	const op = "eventbus.Subscribe"
	log := e.log.With(slog.String("op", op))
	sub, err := e.sp.Subscribe(subject, cb)
	log.Info("subscribe created", slog.String("subject", subject))
	if err != nil {
		log.Error("failed to subscribe", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &loggedSubscription{
		sub:     sub,
		log:     e.log,
		subject: subject,
	}, nil
}

func (e *Eventbus) Publish(subject string, msg interface{}) error {
	const op = "eventbus.Publish"
	log := e.log.With(slog.String("op", op))
	err := e.sp.Publish(subject, msg)
	if err != nil {
		log.Error("publish failed", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("publishing has been comlpleted", slog.String("subject", subject))
	return nil
}

func (e *Eventbus) Close(ctx context.Context) error {
	const op = "eventbus.Close"
	log := e.log.With(slog.String("op", op))
	log.Info("closing the app")
	err := e.sp.Close(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func NewEventBus(log *slog.Logger, sp subpub.SubPub) *Eventbus {
	return &Eventbus{
		log: log,
		sp:  sp,
	}
}
