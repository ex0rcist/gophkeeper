package grpc

import (
	"context"
	"fmt"
	"gophkeeper/internal/keeper/api"
	"gophkeeper/internal/keeper/api/grpc/interceptor"
	"gophkeeper/internal/keeper/config"
	"gophkeeper/internal/keeper/entities"
	"gophkeeper/pkg/constants"
	"gophkeeper/pkg/convert"
	"gophkeeper/pkg/models"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gophkeeper/pkg/proto/keeper/grpcapi"
)

const (
	DefaultClientTimeout = time.Second * 5
)

type GRPCClient struct {
	config        *config.Config
	usersClient   pb.UsersClient
	secretsClient pb.SecretsClient
	notifyClient  pb.NotificationClient
	accessToken   string
	password      string // passw to encrypt payload
	chunkSize     int
	clientID      int32 // Unique ID to distinguish between multiple running clients for same user
	previews      sync.Map
}

var _ api.IApiClient = &GRPCClient{}

func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	var opts []grpc.DialOption

	newClient := GRPCClient{
		config:    cfg,
		chunkSize: constants.ChunkSize,
		clientID:  int32(rand.IntN(math.MaxInt32)),
	}

	// Unary interceptors
	opts = append(
		opts,
		grpc.WithChainUnaryInterceptor(
			interceptor.Timeout(DefaultClientTimeout),
			interceptor.AddAuth(&newClient.accessToken, newClient.clientID),
		),
	)

	// Stream interceptor
	// opts = append(
	// 	opts,
	// 	grpc.WithStreamInterceptor(interceptor.AddAuthStream(&newClient.accessToken, newClient.clientID)),
	// )

	// // TLS
	// if cfg.EnableTLS {
	// 	tlsCredential, err := loadTLSConfig("ca-cert.pem", "client-cert.pem", "client-key.pem")
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to load TLS config: %w", err)
	// 	}

	// 	opts = append(
	// 		opts,
	// 		grpc.WithTransportCredentials(
	// 			tlsCredential,
	// 		),
	// 	)
	// } else {
	// 	opts = append(
	// 		opts,
	// 		grpc.WithTransportCredentials(
	// 			insecure.NewCredentials(),
	// 		),
	// 	)
	// }

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// create gRPC client
	c, err := grpc.NewClient(
		string(cfg.ServerAddress),
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// register services
	newClient.usersClient = pb.NewUsersClient(c)
	newClient.secretsClient = pb.NewSecretsClient(c)
	newClient.notifyClient = pb.NewNotificationClient(c)

	return &newClient, nil
}

// func loadTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
// 	// Read CA cert
// 	caPem, err := cert.Cert.ReadFile(caCertFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read CA cert: %w", err)
// 	}

// 	// Read client cert
// 	clientCertPEM, err := cert.Cert.ReadFile(clientCertFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read client cert: %w", err)
// 	}

// 	// Read client key
// 	clientKeyPEM, err := cert.Cert.ReadFile(clientKeyFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read client key: %w", err)
// 	}

// 	// Create key pair
// 	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
// 	}

// 	// Create cert pool and append CA's cert
// 	certPool := x509.NewCertPool()
// 	if !certPool.AppendCertsFromPEM(caPem) {
// 		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
// 	}

// 	// Create config
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{clientCert},
// 		RootCAs:      certPool,
// 	}

// 	return credentials.NewTLS(config), nil
// }

func (c *GRPCClient) Login(ctx context.Context, login string, password string) (string, error) {
	req := &pb.LoginRequestV1{
		Login:    login,
		Password: password,
	}

	response, err := c.usersClient.LoginV1(ctx, req)
	if err != nil {
		return "", parseError(err)
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

func (c *GRPCClient) Register(ctx context.Context, login string, password string) (string, error) {
	req := &pb.RegisterRequestV1{
		Login:    login,
		Password: password,
	}

	response, err := c.usersClient.RegisterV1(ctx, req)
	if err != nil {
		return "", parseError(err)
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

func (c *GRPCClient) LoadSecrets(ctx context.Context) ([]*models.Secret, error) {
	// form gRPC request
	request := emptypb.Empty{}

	// performing gRPC call
	response, err := c.secretsClient.GetUserSecretsV1(ctx, &request)
	if err != nil {
		return nil, parseError(err)
	}

	secrets := convert.ProtoToSecrets(response.Secrets)
	return secrets, nil
}

func (c *GRPCClient) LoadSecret(ctx context.Context, ID uint64) (*models.Secret, error) {
	// form gRPC request
	request := &pb.GetUserSecretRequestV1{
		Id: ID,
	}

	// performing gRPC call
	response, err := c.secretsClient.GetUserSecretV1(context.Background(), request)
	if err != nil {
		return nil, parseError(err)
	}

	secret := convert.ProtoToSecret(response.Secret)

	return secret, nil
}

func (c *GRPCClient) SaveSecret(ctx context.Context, secret *models.Secret) error {
	sec := &pb.Secret{
		Title:      secret.Title,
		Metadata:   secret.Metadata,
		SecretType: convert.TypeToProto(secret.SecretType),
		Payload:    secret.Payload,
		CreatedAt:  timestamppb.New(secret.CreatedAt),
		UpdatedAt:  timestamppb.New(secret.UpdatedAt),
	}

	if secret.ID > 0 {
		sec.Id = secret.ID
	}

	request := &pb.SaveUserSecretRequestV1{Secret: sec}
	_, err := c.secretsClient.SaveUserSecretV1(ctx, request)

	return parseError(err)
}

func (c *GRPCClient) DeleteSecret(ctx context.Context, id uint64) error {
	request := &pb.DeleteUserSecretRequestV1{Id: id}
	_, err := c.secretsClient.DeleteUserSecretV1(ctx, request)

	return parseError(err)
}

func (c *GRPCClient) SetToken(token string) {
	c.accessToken = token
}

func (c *GRPCClient) GetToken() string {
	return c.accessToken
}

func (c *GRPCClient) SetPassword(password string) {
	c.password = password
}

func (c *GRPCClient) GetPassword() string {
	return c.password
}

func parseError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.Unavailable:
		return entities.ErrServerUnavailable
	case codes.Unauthenticated:
		return entities.ErrUnauthenticated
	case codes.AlreadyExists:
		return entities.ErrAlreadyExist
	default:
		return err
	}
}
