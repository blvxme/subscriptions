package handler

import (
	"context"
	"github.com/blvxme/subpub"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"subscriptions/internal/proto"
)

type Handler struct {
	proto.UnimplementedPubSubServer

	sp     *subpub.SubPub
	logger *logrus.Logger
}

func NewHandler(sp *subpub.SubPub, logger *logrus.Logger) *Handler {
	return &Handler{
		sp:     sp,
		logger: logger,
	}
}

func (h *Handler) Subscribe(req *proto.SubscribeRequest, g grpc.ServerStreamingServer[proto.Event]) error {
	cb := func(msg interface{}) {
		data, ok := msg.(string)
		if ok {
			if err := g.Send(&proto.Event{Data: data}); err != nil {
				h.logger.Errorf("Failed to send message: %+v", err)
			}
		} else {
			h.logger.Errorf("Invalid message: %+v", msg)
		}
	}

	subscription, err := h.sp.Subscribe(req.GetKey(), cb)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to subscribe to the given subject %s: %+v", req.GetKey(), err)
	}
	defer subscription.Unsubscribe()

	subscription.HandleMessagesWithContext(g.Context())

	return status.Error(codes.OK, "OK")
}

func (h *Handler) Publish(ctx context.Context, req *proto.PublishRequest) (*emptypb.Empty, error) {
	if err := h.sp.Publish(req.GetKey(), req.GetData()); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to publish message %+v: %+v:", req, err)
	}
	return &emptypb.Empty{}, status.Error(codes.OK, "OK")
}

func (h *Handler) mustEmbedUnimplementedPubSubServer() {}
