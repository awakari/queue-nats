package service

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_SetQueue(t *testing.T) {
	svc := NewService(NewJetStreamContextMock(), 1)
	cases := map[string]error{
		"ok":       nil,
		"existing": nil,
		"fail":     ErrInternal,
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			err := svc.SetQueue(context.TODO(), k, 10)
			assert.ErrorIs(t, err, c)
		})
	}
}

func TestService_SubmitMessage(t *testing.T) {
	svc := NewService(NewJetStreamContextMock(), 1)
	msg0 := cloudevents.NewEvent()
	cases := map[string]struct {
		queue string
		msg   *event.Event
		err   error
	}{
		"ok": {
			queue: "subj0",
			msg:   &msg0,
		},
		"missing": {
			queue: "missing",
			msg:   &msg0,
			err:   ErrQueueMissing,
		},
		"fail": {
			queue: "fail",
			msg:   &msg0,
			err:   ErrInternal,
		},
		"full": {
			queue: "full",
			msg:   &msg0,
			err:   ErrQueueFull,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			err := svc.SubmitMessage(context.TODO(), c.queue, c.msg)
			assert.ErrorIs(t, err, c.err)
		})
	}
}

func TestService_Poll(t *testing.T) {
	svc := NewService(NewJetStreamContextMock(), 1)
	cases := map[string]struct {
		queue string
		limit uint32
		msgs  []*event.Event
		err   error
	}{
		"fail": {
			queue: "queue0",
			limit: 10,
			err:   ErrInternal,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			msgs, err := svc.Poll(context.TODO(), c.queue, c.limit)
			assert.ErrorIs(t, err, c.err)
			if err == nil {
				assert.Equal(t, c.msgs, msgs)
			}
		})
	}
}
