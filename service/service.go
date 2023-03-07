package service

import (
	"context"
	"errors"
	"fmt"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
	"time"
)

type Service interface {
	Create(ctx context.Context, queue string, limit uint32) (err error)
	SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error)
	PollMessages(ctx context.Context, queue string, limit uint32, timeout time.Duration) (msgs []*event.Event, err error)
}

type service struct {
	js nats.JetStreamContext
}

var ErrQueueAlreadyExists = errors.New("queue already exists")

var ErrQueueCreate = errors.New("failed to create a queue")

var ErrSubmitMessage = errors.New("failed to submit a message")

var ErrPollMessages = errors.New("failed to poll messages")

func NewService(js nats.JetStreamContext) Service {
	return service{
		js: js,
	}
}

func (svc service) Create(ctx context.Context, queue string, limit uint32) (err error) {
	l := int(limit)
	err = svc.addStream(queue, l)
	if err == nil {
		err = svc.addConsumer(queue, l)
	}
	if err != nil {
		if errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
			err = fmt.Errorf("%w: %s", ErrQueueAlreadyExists, queue)
		} else {
			err = fmt.Errorf("%w \"%s\": %s", ErrQueueCreate, queue, err)
		}
	}
	return
}

func (svc service) addStream(name string, limit int) (err error) {
	streamConfig := nats.StreamConfig{
		Name: name,
		Subjects: []string{
			name,
		},
		MaxMsgs:   int64(limit),
		Discard:   nats.DiscardNew,
		Retention: nats.WorkQueuePolicy,
	}
	_, err = svc.js.AddStream(&streamConfig)
	return
}

func (svc service) addConsumer(name string, limit int) (err error) {
	consumerConfig := nats.ConsumerConfig{
		Durable:         name,
		FilterSubject:   name,
		AckPolicy:       nats.AckAllPolicy,
		MaxRequestBatch: limit,
	}
	_, err = svc.js.AddConsumer(name, &consumerConfig)
	return
}

func (svc service) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	var natsMsgData []byte
	natsMsgData, err = format.Protobuf.Marshal(msg)
	if err == nil {
		natsMsg := nats.Msg{
			Subject: queue,
			Data:    natsMsgData,
		}
		_, err = svc.js.PublishMsg(&natsMsg)
	}
	if err != nil {
		err = fmt.Errorf("%w, id: %s, queue: %s, err: %s", ErrSubmitMessage, msg.ID(), queue, err)
	}
	return
}

func (svc service) PollMessages(ctx context.Context, queue string, limit uint32, timeout time.Duration) (msgs []*event.Event, err error) {
	var sub *nats.Subscription
	sub, err = svc.js.PullSubscribe(queue, queue)
	if err == nil {
		l := int(limit)
		err = sub.AutoUnsubscribe(l)
		if err == nil {
			msgs, err = fetch(sub, l, timeout)
		}
	}
	if err != nil {
		err = fmt.Errorf("%w, queue: %s, limit: %d, error: %s", ErrPollMessages, queue, limit, err)
	}
	return
}

func fetch(sub *nats.Subscription, limit int, timeout time.Duration) (msgs []*event.Event, err error) {
	var natsMsgs []*nats.Msg
	natsMsgs, err = sub.Fetch(limit, nats.MaxWait(timeout))
	if errors.Is(err, nats.ErrTimeout) {
		err = nil
	}
	if err == nil {
		for _, natsMsg := range natsMsgs {
			var msg event.Event
			err = format.Protobuf.Unmarshal(natsMsg.Data, &msg)
			if err == nil {
				msgs = append(msgs, &msg)
				err = natsMsg.Ack()
			}
			if err != nil {
				break
			}
		}
	}
	return
}
