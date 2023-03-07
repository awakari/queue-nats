package service

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/exp/slog"
	"time"
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

func (l logging) Create(ctx context.Context, queue string, limit uint32) (err error) {
	defer func() {
		l.log.Debug(fmt.Sprintf("Create(queue=%s, limit=%d): %s", queue, limit, err))
	}()
	return l.svc.Create(ctx, queue, limit)
}

func (l logging) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	defer func() {
		l.log.Debug(fmt.Sprintf("SubmitMessage(queue=%s, msg.Id=%s): %s", queue, msg.ID(), err))
	}()
	return l.svc.SubmitMessage(ctx, queue, msg)
}

func (l logging) PollMessages(ctx context.Context, queue string, limit uint32, timeout time.Duration) (msgs []*event.Event, err error) {
	defer func() {
		l.log.Debug(fmt.Sprintf("PollMessages(queue=%s, limit=%d, timeout=%s): %d, %s", queue, limit, timeout, len(msgs), err))
	}()
	return l.svc.PollMessages(ctx, queue, limit, timeout)
}
