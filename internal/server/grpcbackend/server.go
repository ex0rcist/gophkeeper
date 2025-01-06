package grpcbackend

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// HTTP-server wrapper
type GRPCServer struct {
	address GRPCServerAddress
	server  *grpc.Server
	log     *zap.SugaredLogger
	notify  chan error
}

type GRPCServerDependencies struct {
	dig.In

	Address GRPCServerAddress
	Backend *Backend
	Logger  *zap.SugaredLogger
}

// Constructor
func NewGRPCServer(deps GRPCServerDependencies) *GRPCServer {
	return &GRPCServer{
		address: deps.Address,
		server:  deps.Backend.server,
		log:     deps.Logger,
		notify:  make(chan error, 1),
	}
}

// Run server in a goroutine
func (s *GRPCServer) Start() {
	go func() {
		s.log.Infof("starting GRPC-server on %s", s.address)

		listen, err := net.Listen("tcp", string(s.address))
		if err != nil {
			s.notify <- err
			return
		}

		s.notify <- s.server.Serve(listen)
		close(s.notify)
	}()
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.server.GracefulStop()
	return nil
}

// Return channel to handle errors
func (s *GRPCServer) Notify() <-chan error {
	return s.notify
}

// Describe itself
func (s GRPCServer) String() string {
	return fmt.Sprintf("GRPCserver [addr=%s]", s.address)
}
