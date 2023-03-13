package service

import (
	"context"
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
	cases := map[string]error{
		"ok":      nil,
		"missing": ErrQueueMissing,
		"fail":    ErrInternal,
		"full":    ErrQueueFull,
	}
	for k, expectedErr := range cases {
		t.Run(k, func(t *testing.T) {
			msg := event.New()
			msg.SetID(k)
			err := svc.SubmitMessage(context.TODO(), "queue0", &msg)
			assert.ErrorIs(t, err, expectedErr)
		})
	}
}

func TestService_SubmitMessageBatch(t *testing.T) {
	svc := NewService(NewJetStreamContextMock(), 1)
	cases := map[string]struct {
		msgIds []string
		count  uint32
		err    error
	}{
		"ok": {
			msgIds: []string{
				"msg0",
				"msg1",
				"msg2",
			},
			count: 3,
		},
		"fail on 2nd": {
			msgIds: []string{
				"msg0",
				"fail",
				"msg2",
			},
			count: 1,
			err:   ErrInternal,
		},
		"not enough space in the queue": {
			msgIds: []string{
				"msg0",
				"msg1",
				"full",
			},
			count: 2,
		},
		"queue lost": {
			msgIds: []string{
				"missing",
				"msg1",
				"msg2",
			},
			count: 0,
			err:   ErrQueueMissing,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			var msgs []*event.Event
			for _, msgId := range c.msgIds {
				msg := event.New()
				msg.SetID(msgId)
				msgs = append(msgs, &msg)
			}
			count, err := svc.SubmitMessageBatch(context.TODO(), "queue0", msgs)
			assert.Equal(t, c.count, count)
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
