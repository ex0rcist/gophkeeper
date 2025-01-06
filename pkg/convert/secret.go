// Provides functions for converting between protobuf models and regular models
package convert

import (
	"gophkeeper/pkg/models"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gophkeeper/pkg/proto/keeper/grpcapi"
)

// Returns  secret type
func ProtoToType(pbType pb.SecretType) models.SecretType {
	switch pbType {
	case pb.SecretType_SECRET_TYPE_CREDENTIAL:
		return models.CredSecret
	case pb.SecretType_SECRET_TYPE_TEXT:
		return models.TextSecret
	case pb.SecretType_SECRET_TYPE_BLOB:
		return models.BlobSecret
	case pb.SecretType_SECRET_TYPE_CARD:
		return models.CardSecret
	default:
		return models.UnknownSecret
	}
}

// Returns protobuf secret type
func TypeToProto(sType string) pb.SecretType {
	switch sType {
	case string(models.CredSecret):
		return pb.SecretType_SECRET_TYPE_CREDENTIAL
	case string(models.TextSecret):
		return pb.SecretType_SECRET_TYPE_TEXT
	case string(models.BlobSecret):
		return pb.SecretType_SECRET_TYPE_BLOB
	case string(models.CardSecret):
		return pb.SecretType_SECRET_TYPE_CARD
	default:
		return pb.SecretType_SECRET_TYPE_UNSPECIFIED
	}
}

// Converts secret models to protobuf counterpart
func SecretToProto(secret *models.Secret) *pb.Secret {
	pbSecret := &pb.Secret{
		Id:         secret.ID,
		Title:      secret.Title,
		Metadata:   secret.Metadata,
		Payload:    secret.Payload,
		SecretType: TypeToProto(secret.SecretType),
		CreatedAt:  timestamppb.New(secret.CreatedAt),
		UpdatedAt:  timestamppb.New(secret.UpdatedAt),
	}

	return pbSecret
}

// Converts protobuf models to regular models
func ProtoToSecret(pbSecret *pb.Secret) *models.Secret {
	secret := &models.Secret{
		ID:         pbSecret.Id,
		Title:      pbSecret.Title,
		Metadata:   pbSecret.Metadata,
		SecretType: string(ProtoToType(pbSecret.SecretType)),
		Payload:    pbSecret.Payload,
		CreatedAt:  pbSecret.CreatedAt.AsTime(),
		UpdatedAt:  pbSecret.UpdatedAt.AsTime(),
	}

	return secret
}

// Converts protobuf models to regular models
func ProtoToSecrets(pbSecrets []*pb.Secret) []*models.Secret {
	secr := []*models.Secret{}

	for _, s := range pbSecrets {
		secr = append(secr, ProtoToSecret(s))
	}

	return secr
}

// Converts secret models to protobuf counterpart
func SecretsToProto(secrets []*models.Secret) []*pb.Secret {
	secr := []*pb.Secret{}

	for _, s := range secrets {
		secr = append(secr, SecretToProto(s))
	}

	return secr
}
