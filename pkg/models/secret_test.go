package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecret_MarshalJSON(t *testing.T) {
	createdAt := time.Date(2025, time.January, 6, 16, 25, 40, 721000000, time.UTC).Truncate(time.Millisecond)
	updatedAt := createdAt.Add(time.Hour).Truncate(time.Millisecond)
	credentials := &Credentials{
		Login:    "user1",
		Password: "pass123",
	}

	secret := Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		SecretType: string(CredSecret),
		Creds:      credentials,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	expected := map[string]string{
		"id":          "1",
		"title":       "Test Secret",
		"metadata":    "metadata",
		"secret_type": "credential",
		"payload":     `{"login":"user1","password":"pass123"}`,
		"created_at":  createdAt.Format(timeFormat),
		"updated_at":  updatedAt.Format(timeFormat),
	}

	result, err := json.Marshal(secret)

	assert.NoError(t, err)

	var resultMap map[string]string
	err = json.Unmarshal(result, &resultMap)
	assert.NoError(t, err)
	assert.Equal(t, expected, resultMap)
}

func TestSecret_UnmarshalJSON(t *testing.T) {
	createdAt := time.Date(2025, time.January, 6, 16, 25, 40, 0, time.UTC).Truncate(time.Millisecond)
	updatedAt := createdAt.Add(time.Hour).Truncate(time.Millisecond)
	jsonData := `{
		"id": "1",
		"title": "Test Secret",
		"metadata": "metadata",
		"secret_type": "credential",
		"payload": "{\"login\":\"user1\",\"password\":\"pass123\"}",
		"created_at": "` + createdAt.Format(timeFormat) + `",
		"updated_at": "` + updatedAt.Format(timeFormat) + `"
	}`

	var secret Secret
	err := json.Unmarshal([]byte(jsonData), &secret)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), secret.ID)
	assert.Equal(t, "Test Secret", secret.Title)
	assert.Equal(t, "metadata", secret.Metadata)
	assert.Equal(t, string(CredSecret), secret.SecretType)
	assert.NotNil(t, secret.Creds)
	assert.Equal(t, "user1", secret.Creds.Login)
	assert.Equal(t, "pass123", secret.Creds.Password)
	assert.Equal(t, createdAt, secret.CreatedAt)
	assert.Equal(t, updatedAt, secret.UpdatedAt)
}

func TestSecret_ToClipboard(t *testing.T) {
	credentials := &Credentials{
		Login:    "user1",
		Password: "pass123",
	}
	secret := Secret{
		Title:      "Test Secret",
		SecretType: string(CredSecret),
		Creds:      credentials,
	}

	expected := "login: user1\npassword: pass123"
	result := secret.ToClipboard()
	assert.Equal(t, expected, result)
}

func TestSecret_ToClipboard_Card(t *testing.T) {
	card := &Card{
		Number:   "1234567812345678",
		ExpYear:  2030,
		ExpMonth: 12,
		CVV:      123,
	}
	secret := Secret{
		Title:      "Test Card",
		SecretType: string(CardSecret),
		Card:       card,
	}

	expected := "Card Number: 1234567812345678\nExp: 12/2030CVV: 123"
	result := secret.ToClipboard()
	assert.Equal(t, expected, result)
}

func TestSecret_ToClipboard_Text(t *testing.T) {
	text := &Text{
		Content: "This is a secret note.",
	}
	secret := Secret{
		Title:      "Test Text",
		SecretType: string(TextSecret),
		Text:       text,
	}

	expected := "Text: This is a secret note.\n"
	result := secret.ToClipboard()
	assert.Equal(t, expected, result)
}

func TestSecret_MarshalUnmarshalJSON_RoundTrip(t *testing.T) {
	createdAt := time.Date(2025, time.January, 6, 16, 25, 40, 0, time.UTC).Truncate(time.Millisecond)
	updatedAt := createdAt.Add(time.Hour).Truncate(time.Millisecond)
	credentials := &Credentials{
		Login:    "user1",
		Password: "pass123",
	}

	secret := Secret{
		ID:         1,
		Title:      "Test Secret",
		Metadata:   "metadata",
		SecretType: string(CredSecret),
		Creds:      credentials,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	data, err := json.Marshal(secret)
	assert.NoError(t, err)

	var unmarshaledSecret Secret
	err = json.Unmarshal(data, &unmarshaledSecret)
	assert.NoError(t, err)

	secret.CreatedAt = secret.CreatedAt.Truncate(time.Millisecond)
	secret.UpdatedAt = secret.UpdatedAt.Truncate(time.Millisecond)
	unmarshaledSecret.CreatedAt = unmarshaledSecret.CreatedAt.Truncate(time.Millisecond)
	unmarshaledSecret.UpdatedAt = unmarshaledSecret.UpdatedAt.Truncate(time.Millisecond)

	assert.Equal(t, secret, unmarshaledSecret)
}
