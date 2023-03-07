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
	"time"
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

func (sc serviceController) Create(ctx context.Context, req *CreateRequest) (resp *emptypb.Empty, err error) {
	resp = &emptypb.Empty{}
	err = sc.svc.Create(ctx, req.Queue, req.Limit)
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

func (sc serviceController) PollMessages(ctx context.Context, req *PollMessagesRequest) (resp *PollResponse, err error) {
	var srcMsgs []*event.Event
	var dstMsgs []*pb.CloudEvent
	srcMsgs, err = sc.svc.PollMessages(ctx, req.Queue, req.Limit, time.Millisecond*time.Duration(req.TimeoutMillis))
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
	case errors.Is(src, service.ErrQueueAlreadyExists):
		dst = status.Error(codes.AlreadyExists, src.Error())
	default:
		dst = status.Error(codes.Internal, src.Error())
	}
	return
}