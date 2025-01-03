// Common models for server and client
package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Secret struct {
	ID         uint64    `db:"id" json:"id"`
	UserID     int       `db:"user_id"`
	Title      string    `db:"title" json:"title"`
	Metadata   string    `db:"metadata" json:"metadata"`
	SecretType string    `db:"secret_type" json:"secret_type"`
	Payload    []byte    `db:"payload" json:"payload"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`

	Creds *Credentials `db:"-"`
	Text  *Text        `db:"-"`
	Blob  *Blob        `db:"-"`
	Card  *Card        `db:"-"`
}

type Secrets []*Secret

func NewSecret(t SecretType) *Secret {
	s := Secret{SecretType: string(t)}

	return &s
}

// Secret type
type SecretType string

const (
	CredSecret    SecretType = "credential"
	TextSecret    SecretType = "text"
	BlobSecret    SecretType = "blob"
	CardSecret    SecretType = "card"
	UnknownSecret SecretType = "unknown"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Text struct {
	Content string `json:"content"`
}

type Blob struct {
	FileName  string `json:"file_name"`
	FileBytes []byte `json:"file_bytes"`
	IsDone    bool   `json:"is_done"`
}

type Card struct {
	Number   string `json:"number"`
	ExpYear  uint32 `json:"exp_year"`
	ExpMonth uint32 `json:"exp_month"`
	CVV      uint32 `json:"cvv"`
}

func (s Secret) ToClipboard() string {
	var b bytes.Buffer

	switch SecretType(s.SecretType) {
	case CredSecret:
		b.WriteString(fmt.Sprintf("login: %s\n", s.Creds.Login))
		b.WriteString(fmt.Sprintf("password: %s", s.Creds.Password))
	case CardSecret:
		b.WriteString(fmt.Sprintf("Card Number: %s\n", s.Card.Number))
		b.WriteString(fmt.Sprintf("Exp: %d/%d", s.Card.ExpMonth, s.Card.ExpYear))
		b.WriteString(fmt.Sprintf("CVV: %d", s.Card.CVV))
	case TextSecret:
		b.WriteString(fmt.Sprintf("Text: %s\n", s.Text.Content))
	case BlobSecret:
		// do nothing, file should be saved
	}

	return b.String()
}

const timeFormat = "2006-01-02T15:04:05Z07:00"

// Serialize to JSON
func (s Secret) MarshalJSON() ([]byte, error) {
	var (
		payload []byte
		err     error
	)

	switch SecretType(s.SecretType) {
	case CredSecret:
		payload, err = json.Marshal(s.Creds)
	case CardSecret:
		payload, err = json.Marshal(s.Card)
	case TextSecret:
		payload, err = json.Marshal(s.Text)
	case BlobSecret:
		payload, err = json.Marshal(s.Blob)
	}

	if err != nil {
		return nil, fmt.Errorf("secret payload marshaling failed: %w", err)
	}

	jv, err := json.Marshal(map[string]string{
		"id":          strconv.FormatUint(s.ID, 10),
		"title":       s.Title,
		"secret_type": s.SecretType,
		"metadata":    s.Metadata,
		"payload":     string(payload),
		"created_at":  s.CreatedAt.Format(timeFormat),
		"updated_at":  s.UpdatedAt.Format(timeFormat),
	})

	if err != nil {
		return nil, fmt.Errorf("secret marshaling failed: %w", err)
	}

	return jv, nil
}

// Deserialize from JSON
func (s *Secret) UnmarshalJSON(src []byte) error {
	var (
		data map[string]string
		err  error
	)

	if err = json.Unmarshal(src, &data); err != nil {
		return fmt.Errorf("secret unmarshaling failed: %w", err)
	}

	s.ID, _ = strconv.ParseUint(data["id"], 10, 64)
	s.Title = data["title"]
	s.SecretType = data["secret_type"]
	s.Metadata = data["metadata"]
	s.CreatedAt, _ = time.Parse(timeFormat, data["created_at"])
	s.UpdatedAt, _ = time.Parse(timeFormat, data["updated_at"])

	switch SecretType(data["secret_type"]) {
	case CredSecret:
		s.Creds = &Credentials{}
		err = json.Unmarshal([]byte(data["payload"]), s.Creds)
	case CardSecret:
		s.Card = &Card{}
		err = json.Unmarshal([]byte(data["payload"]), s.Card)
	case TextSecret:
		s.Text = &Text{}
		err = json.Unmarshal([]byte(data["payload"]), s.Text)
	case BlobSecret:
		s.Blob = &Blob{}
		err = json.Unmarshal([]byte(data["payload"]), s.Blob)
	}

	if err != nil {
		return fmt.Errorf("secret payload unmarshaling failed: %w", err)
	}

	return nil
}

// Indicates status of secret for particular client: new, updated or read
// type SecretPreviewStatus string

// const (
// 	SecretPreviewNew     SecretPreviewStatus = "new"
// 	SecretPreviewUpdated SecretPreviewStatus = "updated"
// 	SecretPreviewRead    SecretPreviewStatus = "read"
// )

// Holds just preview secret information: metadata, date, id.
//
// Doesn't include any private user info
// type SecretPreview struct {
// 	ID        uint64
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	Metadata  string
// 	SType     string
// 	Status    SecretPreviewStatus
// }

// type SecretPreviews []*SecretPreview
