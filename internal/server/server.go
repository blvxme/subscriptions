package server

import (
	"context"
	"fmt"
	"github.com/blvxme/subpub"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"subscriptions/internal/proto"
	"syscall"
)

type Server struct {
	grpcServer *grpc.Server
	handler    proto.PubSubServer
	logger     *logrus.Logger
	sp         *subpub.SubPub
}

func NewServer(handler proto.PubSubServer, logger *logrus.Logger, sp *subpub.SubPub) *Server {
	entry := logrus.NewEntry(logger)

	return &Server{
		grpcServer: grpc.NewServer(
			grpc.UnaryInterceptor(
				grpcmiddleware.ChainUnaryServer(
					grpclogrus.UnaryServerInterceptor(entry),
					grpcrecovery.UnaryServerInterceptor(),
				),
			),
			grpc.StreamInterceptor(
				grpcmiddleware.ChainStreamServer(
					grpclogrus.StreamServerInterceptor(entry),
					grpcrecovery.StreamServerInterceptor(),
				),
			),
		),
		handler: handler,
		logger:  logger,
		sp:      sp,
	}
}

func (s *Server) Start(port string) error {
	s.logger.Infof("Starting server on port %s\n", port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	proto.RegisterPubSubServer(s.grpcServer, s.handler)

	go s.gracefulShutdown()

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start serving: %w", err)
	}

	return nil
}

func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")

	if err := s.sp.Close(context.Background()); err != nil {
		s.logger.Errorf("Failed to shut down server: %v", err)
		return
	}

	s.grpcServer.GracefulStop()
}
