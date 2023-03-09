package service

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
)

type serviceMock struct {
	msgs []*event.Event
}

func NewServiceMock(msgs []*event.Event) Service {
	return serviceMock{
		msgs: msgs,
	}
}

func (sm serviceMock) SetQueue(ctx context.Context, queue string, limit uint32) (err error) {
	switch queue {
	case "fail":
		err = ErrInternal
	}
	return
}

func (sm serviceMock) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	switch queue {
	case "missing":
		err = ErrQueueMissing
	case "fail":
		err = ErrInternal
	case "full":
		err = ErrQueueFull
	}
	return
}

func (sm serviceMock) Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error) {
	switch queue {
	case "fail":
		err = ErrInternal
	case "missing":
		err = ErrQueueMissing
	default:
		msgs = sm.msgs
	}
	return
}
