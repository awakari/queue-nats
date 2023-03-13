package grpc

import (
	"context"
	"errors"
	"github.com/awakari/queue-nats/service"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	serviceController struct {
		svc service.Service
	}
)

func NewServiceController(svc service.Service) ServiceServer {
	return serviceController{
		svc: svc,
	}
}

func (sc serviceController) SetQueue(ctx context.Context, req *SetQueueRequest) (resp *emptypb.Empty, err error) {
	resp = &emptypb.Empty{}
	err = sc.svc.SetQueue(ctx, req.Name, req.Limit)
	err = encodeError(err)
	return
}

func (sc serviceController) SubmitMessage(ctx context.Context, req *SubmitMessageRequest) (resp *emptypb.Empty, err error) {
	resp = &emptypb.Empty{}
	var msg *event.Event
	msg, err = format.FromProto(req.Msg)
	if err == nil {
		err = sc.svc.SubmitMessage(ctx, req.Queue, msg)
	}
	err = encodeError(err)
	return
}

func (sc serviceController) SubmitMessageBatch(ctx context.Context, req *SubmitMessageBatchRequest) (resp *BatchResponse, err error) {
	resp = &BatchResponse{}
	var msg *event.Event
	var msgs []*event.Event
	for _, msgProto := range req.Msgs {
		msg, err = format.FromProto(msgProto)
		if err != nil {
			break
		}
		msgs = append(msgs, msg)
	}
	if err == nil {
		resp.Count, err = sc.svc.SubmitMessageBatch(ctx, req.Queue, msgs)
	}
	if err != nil {
		resp.Err = err.Error()
		err = nil // to avoid the nil response when some messages from the batch have been accepted
	}
	return
}

func (sc serviceController) Poll(ctx context.Context, req *PollRequest) (resp *PollResponse, err error) {
	var srcMsgs []*event.Event
	var dstMsgs []*pb.CloudEvent
	srcMsgs, err = sc.svc.Poll(ctx, req.Queue, req.Limit)
	if err == nil {
		var dstMsg *pb.CloudEvent
		for _, srcMsg := range srcMsgs {
			dstMsg, err = format.ToProto(srcMsg)
			if err != nil {
				break
			}
			dstMsgs = append(dstMsgs, dstMsg)
		}
	}
	resp = &PollResponse{
		Msgs: dstMsgs,
	}
	err = encodeError(err)
	return
}

func encodeError(src error) (dst error) {
	switch {
	case src == nil:
		dst = nil
	case errors.Is(src, service.ErrQueueMissing):
		dst = status.Error(codes.NotFound, src.Error())
	case errors.Is(src, service.ErrQueueFull):
		dst = status.Error(codes.ResourceExhausted, src.Error())
	default:
		dst = status.Error(codes.Internal, src.Error())
	}
	return
}
