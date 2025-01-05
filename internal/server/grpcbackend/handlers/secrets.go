package grpchandlers

import (
	"context"
	"errors"
	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/service"
	"gophkeeper/pkg/constants"
	"gophkeeper/pkg/convert"
	"gophkeeper/pkg/proto/keeper/grpcapi"

	"go.uber.org/dig"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "gophkeeper/pkg/proto/keeper/grpcapi"
)

type SecretsServer struct {
	grpcapi.UnimplementedSecretsServer

	secretsManager service.SecretsManager
}

type SecretsServerDependencies struct {
	dig.In

	SecretsManager service.SecretsManager
}

func NewSecretsServer(deps SecretsServerDependencies) *SecretsServer {
	return &SecretsServer{
		secretsManager: deps.SecretsManager,
	}
}

// Saves new secret or updates existing one
func (s *SecretsServer) SaveUserSecretV1(ctx context.Context, in *pb.SaveUserSecretRequestV1) (*emptypb.Empty, error) {
	var err error

	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	secret := convert.ProtoToSecret(in.Secret)
	secret.UserID = int(userID)

	// Save secret
	if secret.ID > 0 {
		_, err = s.secretsManager.UpdateSecret(ctx, secret)
	} else {
		_, err = s.secretsManager.CreateSecret(ctx, secret)
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// // Send notifications
	// clientID, err := extractClientID(ctx)
	// if err == nil {
	// 	if secret.ID > 0 {
	// 		// Update notification
	// 		err = s.notifyClients(userID, clientID, secret.ID, true)
	// 	} else {
	// 		// New secret notification
	// 		err = s.notifyClients(userID, clientID, secretID, false)
	// 	}

	// 	if err != nil {
	// 		s.log.Error("failed to notify clients: ", err)
	// 	}
	// }

	return &emptypb.Empty{}, nil
}

func (s *SecretsServer) GetUserSecretV1(ctx context.Context, in *pb.GetUserSecretRequestV1) (*pb.GetUserSecretResponseV1, error) {
	var response pb.GetUserSecretResponseV1

	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Acquire secret
	secret, err := s.secretsManager.GetSecret(ctx, in.Id, userID)
	if errors.Is(err, entities.ErrSecretNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Other errors
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Secret = convert.SecretToProto(secret)

	return &response, nil
}

func (s *SecretsServer) GetUserSecretsV1(ctx context.Context, in *emptypb.Empty) (*pb.GetUserSecretsResponseV1, error) {
	var response pb.GetUserSecretsResponseV1

	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Acquire secrets
	secrets, err := s.secretsManager.GetUserSecrets(ctx, userID)

	// Other errors
	if err != nil && !errors.Is(err, entities.ErrNoSecrets) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response.Secrets = convert.SecretsToProto(secrets)

	return &response, nil
}

func (s *SecretsServer) DeleteUserSecretV1(ctx context.Context, in *pb.DeleteUserSecretRequestV1) (*emptypb.Empty, error) {

	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Delete
	err = s.secretsManager.DeleteSecret(ctx, in.Id, userID)
	if errors.Is(err, entities.ErrSecretNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Other errors
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func extractUserID(ctx context.Context) (uint64, error) {
	uid := ctx.Value(constants.CtxUserIDKey)

	userID, ok := uid.(uint64)
	if !ok {
		return 0, errors.New("failed to extract user id from context")
	}

	return userID, nil
}
