package grpc

import (
	"context"
	"fmt"
	"github.com/awakari/queue-nats/service"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"os"
	"testing"
)

const port = 8080

var log = slog.Default()
var msg0 = cloudevents.NewEvent()
var msg1 = cloudevents.NewEvent()

func TestMain(m *testing.M) {
	msg0.SetID("3426d090-1b8a-4a09-ac9c-41f2de24d5ac")
	msg0.SetType("type0")
	msg0.SetSource("source0")
	msg0.SetSpecVersion("1.0")
	msg0.SetExtension("foo", "bar")
	msg0.SetData("text/plain", "yohoho")
	msg1.SetID("f7102c87-3ce4-4bb0-8527-b4644f685b13")
	msg1.SetType("type1")
	msg1.SetSource("source1")
	msg1.SetSpecVersion("1.0")
	msg1.SetExtension("bool", true)
	msg1.SetData("application/octet-stream", []byte{1, 2, 3})
	svc := service.NewServiceMock(
		[]*event.Event{
			&msg0,
			&msg1,
		},
	)
	svc = service.NewLoggingMiddleware(svc, log)
	go func() {
		err := Serve(svc, port)
		if err != nil {
			log.Error("", err)
		}
	}()
	code := m.Run()
	os.Exit(code)
}

func TestServiceController_SetQueue(t *testing.T) {
	//
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	client := NewServiceClient(conn)
	//
	cases := map[string]error{
		"fail": status.Error(codes.Internal, "failed to"),
		"ok":   nil,
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			_, err := client.SetQueue(context.TODO(), &SetQueueRequest{
				Name: k,
			})
			assert.ErrorIs(t, err, c)
		})
	}
}

func TestServiceController_SubmitMessage(t *testing.T) {
	//
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	client := NewServiceClient(conn)
	//
	cases := map[string]error{
		"fail":    status.Error(codes.Internal, "failed to"),
		"missing": status.Error(codes.NotFound, "missing queue"),
		"full":    status.Error(codes.ResourceExhausted, "queue is full"),
		"ok":      nil,
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			_, err := client.SubmitMessage(context.TODO(), &SubmitMessageRequest{
				Queue: "queue0",
				Msg: &pb.CloudEvent{
					Id: k,
				},
			})
			assert.ErrorIs(t, err, c)
		})
	}
}

func TestServiceController_SubmitMessageBatch(t *testing.T) {
	//
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	client := NewServiceClient(conn)
	//
	cases := map[string]struct {
		msgIds []string
		resp   *BatchResponse
	}{
		"ok": {
			msgIds: []string{
				"msg0",
				"msg1",
				"msg2",
			},
			resp: &BatchResponse{
				Count: 3,
			},
		},
		"fail on 2nd": {
			msgIds: []string{
				"msg0",
				"fail",
				"msg2",
			},
			resp: &BatchResponse{
				Count: 1,
				Err:   "failed to",
			},
		},
		"not enough space in the queue": {
			msgIds: []string{
				"msg0",
				"msg1",
				"full",
			},
			resp: &BatchResponse{
				Count: 2,
			},
		},
		"queue lost": {
			msgIds: []string{
				"missing",
				"msg1",
				"msg2",
			},
			resp: &BatchResponse{
				Count: 0,
				Err:   "missing queue",
			},
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			var msgs []*pb.CloudEvent
			for _, msgId := range c.msgIds {
				msg := &pb.CloudEvent{
					Id: msgId,
				}
				msgs = append(msgs, msg)
			}
			resp, err := client.SubmitMessageBatch(context.TODO(), &SubmitMessageBatchRequest{
				Queue: k,
				Msgs:  msgs,
			})
			assert.Nil(t, err)
			assert.Equal(t, c.resp.Count, resp.Count)
			assert.Equal(t, c.resp.Err, resp.Err)
		})
	}
}

func TestServiceController_Poll(t *testing.T) {
	//
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	client := NewServiceClient(conn)
	//
	cases := map[string]struct {
		req  *PollRequest
		msgs []*pb.CloudEvent
		err  error
	}{
		"fail": {
			req: &PollRequest{Queue: "fail"},
			err: status.Error(codes.Internal, "failed to"),
		},
		"missing": {
			req: &PollRequest{Queue: "missing"},
			err: status.Error(codes.NotFound, "missing queue"),
		},
		"ok": {
			req: &PollRequest{},
			msgs: []*pb.CloudEvent{
				{
					Id:          "3426d090-1b8a-4a09-ac9c-41f2de24d5ac",
					Source:      "source0",
					SpecVersion: "1.0",
					Type:        "type0",
					Attributes: map[string]*pb.CloudEventAttributeValue{
						"foo": {
							Attr: &pb.CloudEventAttributeValue_CeString{
								CeString: "bar",
							},
						},
						"datacontenttype": {
							Attr: &pb.CloudEventAttributeValue_CeString{
								CeString: "text/plain",
							},
						},
					},
					Data: &pb.CloudEvent_BinaryData{
						BinaryData: []byte("yohoho"),
					},
				},
				{
					Id:          "f7102c87-3ce4-4bb0-8527-b4644f685b13",
					Source:      "source1",
					SpecVersion: "1.0",
					Type:        "type1",
					Attributes: map[string]*pb.CloudEventAttributeValue{
						"bool": {
							Attr: &pb.CloudEventAttributeValue_CeBoolean{
								CeBoolean: true,
							},
						},
						"datacontenttype": {
							Attr: &pb.CloudEventAttributeValue_CeString{
								CeString: "application/octet-stream",
							},
						},
					},
					Data: &pb.CloudEvent_BinaryData{
						BinaryData: []byte{1, 2, 3},
					},
				},
			},
		},
	}
	//
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			resp, err := client.Poll(context.TODO(), c.req)
			assert.ErrorIs(t, err, c.err)
			if err == nil {
				msgs := resp.Msgs
				assert.Equal(t, len(c.msgs), len(msgs))
				for i, msg := range msgs {
					assert.Equal(t, c.msgs[i].Id, msg.Id)
					assert.Equal(t, c.msgs[i].Data, msg.Data)
					assert.Equal(t, c.msgs[i].Type, msg.Type)
					assert.Equal(t, c.msgs[i].Source, msg.Source)
					assert.Equal(t, c.msgs[i].SpecVersion, msg.SpecVersion)
					assert.Equal(t, len(c.msgs[i].Attributes), len(msg.Attributes))
					for attrK, attrV := range c.msgs[i].Attributes {
						assert.Equal(t, attrV, msg.Attributes[attrK])
					}
				}
			}
		})
	}
}
