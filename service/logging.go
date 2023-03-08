package service

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/exp/slog"
)

type logging struct {
	svc Service
	log *slog.Logger
}

func NewLogging(svc Service, log *slog.Logger) Service {
	return logging{
		svc: svc,
		log: log,
	}
}

func (l logging) SetQueue(ctx context.Context, name string, limit uint32) (err error) {
	defer func() {
		l.log.Debug(fmt.Sprintf("SetQueue(name=%s,  limit=%d): %s", name, limit, err))
	}()
	return l.svc.SetQueue(ctx, name, limit)
}

func (l logging) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	defer func() {
		l.log.Debug(fmt.Sprintf("SubmitMessage(queue=%s, msg.Id=%s): %s", queue, msg.ID(), err))
	}()
	return l.svc.SubmitMessage(ctx, queue, msg)
}

func (l logging) Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error) {
	defer func() {
		l.log.Debug(fmt.Sprintf("Poll(queue=%s, limit=%d): %d, %s", queue, limit, len(msgs), err))
	}()
	return l.svc.Poll(ctx, queue, limit)
}
