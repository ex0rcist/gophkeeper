package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"gophkeeper/pkg/models"
	pb "gophkeeper/pkg/proto/keeper/grpcapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockUsersClient is a mock implementation of pb.UsersClient.
type MockUsersClient struct {
	mock.Mock
}

func (m *MockUsersClient) LoginV1(ctx context.Context, req *pb.LoginRequestV1, opts ...grpc.CallOption) (*pb.LoginResponseV1, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.LoginResponseV1), args.Error(1)
}

func (m *MockUsersClient) RegisterV1(ctx context.Context, req *pb.RegisterRequestV1, opts ...grpc.CallOption) (*pb.RegisterResponseV1, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.RegisterResponseV1), args.Error(1)
}

// MockSecretsClient is a mock implementation of pb.SecretsClient.
type MockSecretsClient struct {
	mock.Mock
}

func (m *MockSecretsClient) GetUserSecretsV1(ctx context.Context, req *emptypb.Empty, opts ...grpc.CallOption) (*pb.GetUserSecretsResponseV1, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetUserSecretsResponseV1), args.Error(1)
}

func (m *MockSecretsClient) GetUserSecretV1(ctx context.Context, req *pb.GetUserSecretRequestV1, opts ...grpc.CallOption) (*pb.GetUserSecretResponseV1, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetUserSecretResponseV1), args.Error(1)
}

func (m *MockSecretsClient) SaveUserSecretV1(ctx context.Context, req *pb.SaveUserSecretRequestV1, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockSecretsClient) DeleteUserSecretV1(ctx context.Context, req *pb.DeleteUserSecretRequestV1, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func TestGRPCClient_Login(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockUsersClient := new(MockUsersClient)
		client := &GRPCClient{usersClient: mockUsersClient}

		mockUsersClient.On("LoginV1", mock.Anything, mock.Anything).Return(&pb.LoginResponseV1{AccessToken: "test-token"}, nil)

		token, err := client.Login(context.Background(), "testuser", "testpass")

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("Error", func(t *testing.T) {
		mockUsersClient := new(MockUsersClient)
		client := &GRPCClient{usersClient: mockUsersClient}

		mockUsersClient.On("LoginV1", mock.Anything, mock.Anything).Return(nil, errors.New("login error"))

		token, err := client.Login(context.Background(), "testuser", "testpass")

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestGRPCClient_Register(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockUsersClient := new(MockUsersClient)
		client := &GRPCClient{usersClient: mockUsersClient}

		mockUsersClient.On("RegisterV1", mock.Anything, mock.Anything).Return(&pb.RegisterResponseV1{AccessToken: "test-token"}, nil)

		token, err := client.Register(context.Background(), "testuser", "testpass")

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
	})

	t.Run("Error", func(t *testing.T) {
		mockUsersClient := new(MockUsersClient)
		client := &GRPCClient{usersClient: mockUsersClient}

		mockUsersClient.On("RegisterV1", mock.Anything, mock.Anything).Return(nil, errors.New("register error"))

		token, err := client.Register(context.Background(), "testuser", "testpass")

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestGRPCClient_LoadSecrets(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockSecretsClient := new(MockSecretsClient)
		client := &GRPCClient{secretsClient: mockSecretsClient}

		mockSecretsClient.On("GetUserSecretsV1", mock.Anything, mock.Anything).Return(&pb.GetUserSecretsResponseV1{Secrets: []*pb.Secret{{Id: 1, Title: "Test Secret"}}}, nil)

		secrets, err := client.LoadSecrets(context.Background())

		assert.NoError(t, err)
		assert.Len(t, secrets, 1)
		assert.Equal(t, uint64(1), secrets[0].ID)
		assert.Equal(t, "Test Secret", secrets[0].Title)
	})

	t.Run("Error", func(t *testing.T) {
		mockSecretsClient := new(MockSecretsClient)
		client := &GRPCClient{secretsClient: mockSecretsClient}

		mockSecretsClient.On("GetUserSecretsV1", mock.Anything, mock.Anything).Return(nil, errors.New("load error"))

		secrets, err := client.LoadSecrets(context.Background())

		assert.Error(t, err)
		assert.Nil(t, secrets)
	})
}

func TestGRPCClient_SaveSecret(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockSecretsClient := new(MockSecretsClient)
		client := &GRPCClient{secretsClient: mockSecretsClient}

		mockSecretsClient.On("SaveUserSecretV1", mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

		secret := &models.Secret{ID: 1, Title: "Test Secret", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err := client.SaveSecret(context.Background(), secret)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockSecretsClient := new(MockSecretsClient)
		client := &GRPCClient{secretsClient: mockSecretsClient}

		mockSecretsClient.On("SaveUserSecretV1", mock.Anything, mock.Anything).Return(nil, errors.New("save error"))

		secret := &models.Secret{ID: 1, Title: "Test Secret"}
		err := client.SaveSecret(context.Background(), secret)

		assert.Error(t, err)
	})
}

func TestGRPCClient_DeleteSecret(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSecretsClient := new(MockSecretsClient)
		client := &GRPCClient{secretsClient: mockSecretsClient}

		mockSecretsClient.On("DeleteUserSecretV1", mock.Anything, mock.Anything).Return(&emptypb.Empty{}, nil)

		err := client.DeleteSecret(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockSecretsClient := new(MockSecretsClient)
		client := &GRPCClient{secretsClient: mockSecretsClient}

		mockSecretsClient.On("DeleteUserSecretV1", mock.Anything, mock.Anything).Return(nil, errors.New("delete error"))

		err := client.DeleteSecret(context.Background(), 1)

		assert.Error(t, err)
	})
}
