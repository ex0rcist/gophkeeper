package grpchandlers

import (
	"context"
	"errors"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/service"
	"gophkeeper/pkg/constants"
	"gophkeeper/pkg/proto/keeper/grpcapi"
	"strconv"
	"time"

	pb "gophkeeper/pkg/proto/keeper/grpcapi"

	"go.uber.org/dig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const tokenLifetime = time.Hour * 24

type UsersServer struct {
	grpcapi.UnimplementedUsersServer

	config       *config.Config
	usersManager service.UsersManager
}

type UsersServerDependencies struct {
	dig.In

	Config       *config.Config
	UsersManager service.UsersManager
}

func NewUsersServer(deps UsersServerDependencies) *UsersServer {
	return &UsersServer{
		config:       deps.Config,
		usersManager: deps.UsersManager,
	}
}

func (s *UsersServer) RegisterV1(ctx context.Context, in *pb.RegisterRequestV1) (*pb.RegisterResponseV1, error) {
	var response pb.RegisterResponseV1

	// Register user
	user, err := s.usersManager.RegisterUser(ctx, in.Login, in.Password)

	// Check if user exists
	if errors.Is(err, entities.ErrUserAlreadyExists) {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	// Other errors
	if err != nil && !errors.Is(err, entities.ErrUserNotFound) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Generate access token
	token, err := s.authUser(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to auth: %v", err)
	}

	response.AccessToken = token

	return &response, nil
}

func (s *UsersServer) LoginV1(ctx context.Context, in *pb.LoginRequestV1) (*pb.LoginResponseV1, error) {
	var response pb.LoginResponseV1

	// Login user
	user, err := s.usersManager.LoginUser(ctx, in.Login, in.Password)

	// Check credentials
	if errors.Is(err, entities.ErrBadCredentials) {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Other errors
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	// Generate access token
	token, err := s.authUser(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to auth: %v", err)
	}

	response.AccessToken = token

	return &response, nil
}

func (s *UsersServer) authUser(userID int) (string, error) {
	token, err := auth.CreateToken(userID, time.Now().Add(tokenLifetime), []byte(s.config.SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func extractClientID(ctx context.Context) (int32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("failed to get metadata")
	}

	values := md.Get(constants.ClientIDHeader)
	if len(values) == 0 {
		return 0, errors.New("missing client id metadata")
	}

	v := values[0]

	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}
