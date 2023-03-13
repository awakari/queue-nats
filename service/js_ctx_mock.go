package service

import (
	"errors"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/nats-io/nats.go"
)

type jsCtxMock struct {
}

func NewJetStreamContextMock() nats.JetStreamContext {
	return jsCtxMock{}
}

func (js jsCtxMock) Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) PublishMsg(m *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error) {
	var msg event.Event
	_ = format.Protobuf.Unmarshal(m.Data, &msg)
	if msg.ID() == "fail" {
		return nil, errors.New("failed")
	}
	if msg.ID() == "full" {
		return nil, errors.New("nats: maximum messages exceeded")
	}
	if msg.ID() == "missing" {
		return nil, nats.ErrNoStreamResponse
	}
	return &nats.PubAck{}, nil
}

func (js jsCtxMock) PublishAsync(subj string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) PublishMsgAsync(m *nats.Msg, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) PublishAsyncPending() int {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) PublishAsyncComplete() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) Subscribe(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) SubscribeSync(subj string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ChanSubscribe(subj string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ChanQueueSubscribe(subj, queue string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) QueueSubscribe(subj, queue string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) QueueSubscribeSync(subj, queue string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) PullSubscribe(subj, durable string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	if subj == "fail" {
		return nil, errors.New("failed")
	}
	if subj == "missing" {
		return nil, nats.ErrNoMatchingStream
	}
	return &nats.Subscription{
		Subject: subj,
		Queue:   durable,
	}, nil
}

func (js jsCtxMock) AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	if cfg.Name == "existing" {
		return nil, nats.ErrStreamNameAlreadyInUse
	}
	if cfg.Name == "fail" {
		return nil, errors.New("fail")
	}
	return &nats.StreamInfo{
		Config: *cfg,
	}, nil
}

func (js jsCtxMock) UpdateStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	return &nats.StreamInfo{
		Config: *cfg,
	}, nil
}

func (js jsCtxMock) DeleteStream(name string, opts ...nats.JSOpt) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) StreamInfo(stream string, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) PurgeStream(name string, opts ...nats.JSOpt) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) StreamsInfo(opts ...nats.JSOpt) <-chan *nats.StreamInfo {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) Streams(opts ...nats.JSOpt) <-chan *nats.StreamInfo {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) StreamNames(opts ...nats.JSOpt) <-chan string {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) GetMsg(name string, seq uint64, opts ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) GetLastMsg(name, subject string, opts ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) DeleteMsg(name string, seq uint64, opts ...nats.JSOpt) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) SecureDeleteMsg(name string, seq uint64, opts ...nats.JSOpt) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) AddConsumer(stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	return &nats.ConsumerInfo{
		Stream: stream,
		Name:   cfg.Name,
	}, nil
}

func (js jsCtxMock) UpdateConsumer(stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	return &nats.ConsumerInfo{
		Stream: stream,
		Name:   cfg.Name,
	}, nil
}

func (js jsCtxMock) DeleteConsumer(stream, consumer string, opts ...nats.JSOpt) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ConsumerInfo(stream, name string, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	if stream == "missing" {
		return nil, nats.ErrConsumerNotFound
	}
	return &nats.ConsumerInfo{
		Stream: stream,
		Name:   stream,
	}, nil
}

func (js jsCtxMock) ConsumersInfo(stream string, opts ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) Consumers(stream string, opts ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ConsumerNames(stream string, opts ...nats.JSOpt) <-chan string {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) AccountInfo(opts ...nats.JSOpt) (*nats.AccountInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) StreamNameBySubject(s string, opt ...nats.JSOpt) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) KeyValue(bucket string) (nats.KeyValue, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) CreateKeyValue(cfg *nats.KeyValueConfig) (nats.KeyValue, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) DeleteKeyValue(bucket string) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) KeyValueStoreNames() <-chan string {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) KeyValueStores() <-chan nats.KeyValueStatus {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ObjectStore(bucket string) (nats.ObjectStore, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) CreateObjectStore(cfg *nats.ObjectStoreConfig) (nats.ObjectStore, error) {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) DeleteObjectStore(bucket string) error {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ObjectStoreNames(opts ...nats.ObjectOpt) <-chan string {
	//TODO implement me
	panic("implement me")
}

func (js jsCtxMock) ObjectStores(opts ...nats.ObjectOpt) <-chan nats.ObjectStoreStatus {
	//TODO implement me
	panic("implement me")
}
