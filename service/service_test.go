package service

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Create(t *testing.T) {
	svc := NewService(NewJetStreamContextMock())
	cases := map[string]error{
		"ok":       nil,
		"fail":     ErrQueueCreate,
		"existing": ErrQueueAlreadyExists,
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			err := svc.Create(context.TODO(), k, 10)
			assert.ErrorIs(t, err, c)
		})
	}
}

func TestService_SubmitMessage(t *testing.T) {
	svc := NewService(NewJetStreamContextMock())
	msg0 := cloudevents.NewEvent()
	cases := map[string]struct {
		queue string
		msg   *event.Event
		err   error
	}{
		"ok": {
			queue: "queue0",
			msg:   &msg0,
		},
		"fail": {
			queue: "fail",
			msg:   &msg0,
			err:   ErrSubmitMessage,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			err := svc.SubmitMessage(context.TODO(), c.queue, c.msg)
			assert.ErrorIs(t, err, c.err)
		})
	}
}

func TestService_PollMessages(t *testing.T) {
	svc := NewService(NewJetStreamContextMock())
	cases := map[string]struct {
		queue string
		limit uint32
		msgs  []*event.Event
		err   error
	}{
		"fail": {
			queue: "q0",
			limit: 10,
			err:   ErrPollMessages,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			msgs, err := svc.PollMessages(context.TODO(), c.queue, c.limit, 1000)
			assert.ErrorIs(t, err, c.err)
			if err == nil {
				assert.Equal(t, c.msgs, msgs)
			}
		})
	}
}
