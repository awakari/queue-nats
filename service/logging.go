package service

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/exp/slog"
)

type loggingMiddleware struct {
	svc Service
	log *slog.Logger
}

func NewLoggingMiddleware(svc Service, log *slog.Logger) Service {
	return loggingMiddleware{
		svc: svc,
		log: log,
	}
}

func (lm loggingMiddleware) SetQueue(ctx context.Context, name string, limit uint32) (err error) {
	defer func() {
		lm.log.Debug(fmt.Sprintf("SetQueue(name=%s,  limit=%d): %s", name, limit, err))
	}()
	return lm.svc.SetQueue(ctx, name, limit)
}

func (lm loggingMiddleware) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	defer func() {
		lm.log.Debug(fmt.Sprintf("SubmitMessage(queue=%s, msg.Id=%s): %s", queue, msg.ID(), err))
	}()
	return lm.svc.SubmitMessage(ctx, queue, msg)
}

func (lm loggingMiddleware) Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error) {
	defer func() {
		lm.log.Debug(fmt.Sprintf("Poll(queue=%s, limit=%d): %d, %s", queue, limit, len(msgs), err))
	}()
	return lm.svc.Poll(ctx, queue, limit)
}
