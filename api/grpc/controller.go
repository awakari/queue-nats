package grpc

import (
	"context"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"output-log/service"
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

func (sc serviceController) Submit(ctx context.Context, req *pb.CloudEvent) (resp *emptypb.Empty, err error) {
	var msg *event.Event
	msg, err = format.FromProto(req)
	if err == nil {
		err = sc.svc.Submit(ctx, msg)
	}
	resp = &emptypb.Empty{}
	err = encodeError(err)
	return
}

func encodeError(svcErr error) (err error) {
	switch {
	case svcErr == nil:
		err = nil
	default:
		err = status.Error(codes.Internal, svcErr.Error())
	}
	return
}
