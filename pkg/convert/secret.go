// Provides functions for converting between protobuf models and regular models
package convert

import (
	"gophkeeper/pkg/models"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gophkeeper/pkg/proto/keeper/v1"
)

// Returns corresponding secret type
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

// Returns corresponding protobuf secret type
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
		Id:        secret.ID,
		CreatedAt: timestamppb.New(secret.CreatedAt),
		UpdatedAt: timestamppb.New(secret.UpdatedAt),
		Metadata:  secret.Metadata,
		Type:      pb.SecretType_SECRET_TYPE_UNSPECIFIED,
	}

	pbSecret.Type = TypeToProto(secret.SecretType)

	switch secret.SecretType {
	case string(models.CredSecret):
		pbSecret.Content = &pb.Secret_Credential{
			Credential: &pb.Credential{
				Login:    secret.Creds.Login,
				Password: secret.Creds.Password,
			},
		}
	case string(models.TextSecret):
		pbSecret.Content = &pb.Secret_Text{
			Text: &pb.Text{
				Text: secret.Text.Content,
			},
		}
	case string(models.BlobSecret):
		pbSecret.Content = &pb.Secret_Blob{
			Blob: &pb.Blob{
				FileName: secret.Blob.FileName,
			},
		}
	case string(models.CardSecret):
		pbSecret.Content = &pb.Secret_Card{
			Card: &pb.Card{
				Number:   secret.Card.Number,
				ExpYear:  secret.Card.ExpYear,
				ExpMonth: secret.Card.ExpMonth,
				Cvv:      secret.Card.CVV,
			},
		}
	}

	return pbSecret
}

// Converts protobuf models to regular models
func ProtoToSecret(pbSecret *pb.Secret) *models.Secret {
	secret := &models.Secret{
		ID:         pbSecret.Id,
		CreatedAt:  pbSecret.CreatedAt.AsTime(),
		UpdatedAt:  pbSecret.UpdatedAt.AsTime(),
		Metadata:   pbSecret.Metadata,
		SecretType: string(models.UnknownSecret),
	}

	secret.SecretType = string(ProtoToType(pbSecret.Type))

	switch secret.SecretType {
	case string(models.CredSecret):
		pbCred := pbSecret.Content.(*pb.Secret_Credential)
		secret.Creds = &models.Credentials{
			Login:    pbCred.Credential.GetLogin(),
			Password: pbCred.Credential.GetPassword(),
		}
	case string(models.TextSecret):
		pbText := pbSecret.Content.(*pb.Secret_Text)
		secret.Text = &models.Text{
			Content: pbText.Text.GetText(),
		}
	case string(models.BlobSecret):
		pbBlob := pbSecret.Content.(*pb.Secret_Blob)
		secret.Blob = &models.Blob{
			FileName: pbBlob.Blob.GetFileName(),
		}
	case string(models.CardSecret):
		pbCard := pbSecret.Content.(*pb.Secret_Card)
		secret.Card = &models.Card{
			Number:   pbCard.Card.GetNumber(),
			ExpYear:  pbCard.Card.GetExpYear(),
			ExpMonth: pbCard.Card.GetExpMonth(),
			CVV:      pbCard.Card.GetCvv(),
		}
	}

	return secret
}
