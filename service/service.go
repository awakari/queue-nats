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
	SetQueue(ctx context.Context, name string, limit uint32) (err error)
	SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error)
	SubmitMessageBatch(ctx context.Context, queue string, msgs []*event.Event) (count uint32, err error)
	Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error)
}

type service struct {
	js          nats.JetStreamContext
	pollTimeout time.Duration
}

var ErrQueueMissing = errors.New("missing queue")

var ErrQueueFull = errors.New("queue is full")

var ErrInternal = errors.New("failed to")

func NewService(js nats.JetStreamContext, pollTimeoutMillis uint32) Service {
	return service{
		js:          js,
		pollTimeout: time.Millisecond * time.Duration(pollTimeoutMillis),
	}
}

func (svc service) SetQueue(ctx context.Context, name string, limit uint32) (err error) {
	l := int(limit)
	err = svc.addStream(name, name, l)
	if err == nil {
		err = svc.addConsumer(name, name, l)
	}
	if err != nil {
		err = fmt.Errorf("%w create queue \"%s\": %s", ErrInternal, name, err)
	}
	return
}

func (svc service) addStream(queue, subject string, limit int) (err error) {
	streamConfig := nats.StreamConfig{
		Name: queue,
		Subjects: []string{
			subject,
		},
		MaxMsgs:   int64(limit),
		Discard:   nats.DiscardNew,
		Retention: nats.WorkQueuePolicy,
	}
	_, err = svc.js.AddStream(&streamConfig)
	if errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		_, err = svc.js.UpdateStream(&streamConfig)
	}
	return
}

func (svc service) addConsumer(queue, subject string, limit int) (err error) {
	consumerConfig := nats.ConsumerConfig{
		Name:            queue,
		Durable:         queue,
		FilterSubject:   subject,
		AckPolicy:       nats.AckExplicitPolicy,
		MaxRequestBatch: limit,
	}
	_, err = svc.js.ConsumerInfo(queue, queue)
	if errors.Is(err, nats.ErrConsumerNotFound) || errors.Is(err, nats.ErrStreamNotFound) {
		_, err = svc.js.AddConsumer(queue, &consumerConfig)
	} else {
		_, err = svc.js.UpdateConsumer(queue, &consumerConfig)
	}
	return
}

func (svc service) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	var natsMsg *nats.Msg
	natsMsg, err = convertToNatsMessage(queue, msg)
	if err == nil {
		_, err = svc.js.PublishMsg(natsMsg)
	}
	if err != nil {
		switch {
		case errors.Is(err, nats.ErrNoStreamResponse):
			err = fmt.Errorf("%w \"%s\": failed to submit the message with id \"%s\"", ErrQueueMissing, queue, msg.ID())
		case err.Error() == "nats: maximum messages exceeded":
			err = fmt.Errorf("%w: %s, message id: %s", ErrQueueFull, queue, msg.ID())
		default:
			err = fmt.Errorf("%w publish message id \"%s\" to the queue \"%s\": %s", ErrInternal, msg.ID(), queue, err)
		}
	}
	return
}

func (svc service) SubmitMessageBatch(ctx context.Context, queue string, msgs []*event.Event) (count uint32, err error) {
	var natsMsg *nats.Msg
	var natsMsgs []*nats.Msg
	for _, msg := range msgs {
		natsMsg, err = convertToNatsMessage(queue, msg)
		if err != nil {
			break
		}
		natsMsgs = append(natsMsgs, natsMsg)
	}
	if err == nil {
		for i, natsMsg := range natsMsgs {
			_, err = svc.js.PublishMsg(natsMsg)
			if err != nil {
				switch {
				case errors.Is(err, nats.ErrNoStreamResponse):
					err = fmt.Errorf("%w \"%s\": failed to submit the message with id \"%s\"", ErrQueueMissing, queue, msgs[i].ID())
				case err.Error() == "nats: maximum messages exceeded":
					err = nil // just exit the loop and return the current count
				default:
					err = fmt.Errorf("%w publish message id \"%s\" to the queue \"%s\": %s", ErrInternal, msgs[i].ID(), queue, err)
				}
				break
			}
			count++
		}
	}
	return
}

func convertToNatsMessage(queue string, src *event.Event) (dst *nats.Msg, err error) {
	var natsMsgData []byte
	natsMsgData, err = format.Protobuf.Marshal(src)
	if err == nil {
		dst = &nats.Msg{
			Subject: queue,
			Data:    natsMsgData,
		}
	}
	return
}

func (svc service) Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error) {
	var sub *nats.Subscription
	sub, err = svc.js.PullSubscribe(queue, queue)
	if err == nil {
		defer sub.Unsubscribe()
		msgs, err = svc.fetch(sub, limit)
	}
	if err != nil {
		if errors.Is(err, nats.ErrNoMatchingStream) {
			err = fmt.Errorf("%w \"%s\": failed to poll messages", ErrQueueMissing, queue)
		} else {
			err = fmt.Errorf("%w poll up to %d messages from the queue \"%s\": %s", ErrInternal, limit, queue, err)
		}
	}
	return
}

func (svc service) fetch(sub *nats.Subscription, limit uint32) (msgs []*event.Event, err error) {
	var natsMsgs []*nats.Msg
	natsMsgs, err = sub.Fetch(int(limit), nats.MaxWait(svc.pollTimeout))
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
