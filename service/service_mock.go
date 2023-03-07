package service

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"time"
)

type serviceMock struct {
	msgs []*event.Event
}

func NewServiceMock(msgs []*event.Event) Service {
	return serviceMock{
		msgs: msgs,
	}
}

func (sm serviceMock) Create(ctx context.Context, queue string, limit uint32) (err error) {
	switch queue {
	case "existing":
		err = ErrQueueAlreadyExists
	case "fail":
		err = ErrQueueCreate
	}
	return
}

func (sm serviceMock) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	if queue == "fail" {
		err = ErrSubmitMessage
	}
	return
}

func (sm serviceMock) PollMessages(ctx context.Context, queue string, limit uint32, timeout time.Duration) (msgs []*event.Event, err error) {
	if queue == "fail" {
		err = ErrPollMessages
	} else {
		msgs = sm.msgs
	}
	return
}
