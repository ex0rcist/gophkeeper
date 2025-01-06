package grpchandlers

import (
	"errors"
	"slices"
	"sync"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gophkeeper/internal/server/entities"
	pb "gophkeeper/pkg/proto/keeper/grpcapi"
)

type sub struct {
	stream   pb.Notification_SubscribeV1Server
	id       uint64
	finished chan<- bool
}

// HealthServer verifies current health status of the service.
type NotificationServer struct {
	pb.UnimplementedNotificationServer

	logger      *zap.SugaredLogger
	subscribers sync.Map
}

type NotificationServerDependencies struct {
	dig.In
	Logger *zap.SugaredLogger
}

func NewNotificationServer(deps NotificationServerDependencies) *NotificationServer {
	return &NotificationServer{logger: deps.Logger}
}

func (s *NotificationServer) SubscribeV1(in *pb.SubscribeV1Request, stream pb.Notification_SubscribeV1Server) error {
	var subs []sub

	ctx := stream.Context()

	userID, err := extractUserID(ctx)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("received subscribe from client #", in.Id, "user ID", userID)

	fin := make(chan bool)

	v, ok := s.subscribers.Load(userID)
	if ok {
		// Try to cast value into sub slice
		subs, ok = v.([]sub)
		if !ok {
			return status.Error(codes.Internal, "failed to cast subscribers")
		}
	}

	// Append to subscribers slice
	subs = append(
		subs,
		sub{
			stream:   stream,
			id:       in.Id,
			finished: fin,
		},
	)

	// Store subs in map
	s.subscribers.Store(userID, subs)

	for {
		select {
		case <-fin:
			s.logger.Infof("closing stream for client #%d", in.Id)
			return nil
		case <-ctx.Done():
			s.logger.Infof("client #%d has disconnected", in.Id)
			return nil
		}
	}
}

func (s *NotificationServer) notifyClients(userID uint64, clientID uint64, ID uint64, updated bool) error {
	v, ok := s.subscribers.Load(userID)
	if !ok {
		return entities.ErrNoSubscribers
	}

	subs, ok := v.([]sub)
	if !ok {
		return errors.New("failed to cast to subs")
	}

	var unsubs []int

	for i, sub := range subs {
		if sub.id == clientID {
			// Skip originating client
			continue
		}

		resp := &pb.SubscribeResponseV1{
			Id:      ID,
			Updated: updated,
		}

		if err := sub.stream.Send(resp); err != nil {
			s.logger.Error("failed to send notification to client: %v", err)
			select {
			case sub.finished <- true:
				s.logger.Info("client unsubscribed: %v", clientID)
			default:
			}

			// Mark as unsub
			unsubs = append(unsubs, i)
		}
	}

	// Delete unsubs from slice
	for _, unsub := range unsubs {
		subs = slices.Delete(subs, unsub, unsub)
	}

	if len(subs) > 0 {
		s.subscribers.Store(userID, subs)
	} else {
		// All clients unsubscribed
		s.subscribers.Delete(userID)
	}

	return nil
}
