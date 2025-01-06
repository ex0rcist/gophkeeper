package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/repository"
	strg "gophkeeper/internal/server/storage/postgres"
	"gophkeeper/pkg/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
)

var _ repository.SecretsRepository = SecretsRepository{}

type SecretsRepositoryDependencies struct {
	dig.In
	PostgresConn *strg.PostgresConn
}

// Secrets repository using PostgreSQL
type SecretsRepository struct {
	db *sqlx.DB
}

// Create new postgresql secret repository
func NewSecretsRepository(deps SecretsRepositoryDependencies) *SecretsRepository {
	return &SecretsRepository{
		db: deps.PostgresConn.DB,
	}
}

// Find secret by id
func (r SecretsRepository) GetSecret(ctx context.Context, secretID uint64, userID uint64) (*models.Secret, error) {
	var secret models.Secret

	query := `SELECT * FROM secrets WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRowxContext(ctx, query, secretID, userID).StructScan(&secret)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrUserNotFound
	}

	return &secret, err
}

// Find user's secrets
func (r SecretsRepository) GetUserSecrets(ctx context.Context, userID uint64) (models.Secrets, error) {
	var secrets models.Secrets

	query := "SELECT * FROM secrets WHERE user_id = $1 ORDER BY updated_at DESC"
	err := r.db.SelectContext(ctx, &secrets, query, userID)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

// Create new secret
func (r SecretsRepository) Create(ctx context.Context, secret *models.Secret) (uint64, error) {
	var newSecretID uint64

	query := `INSERT INTO secrets (user_id, title, metadata, secret_type, payload)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	result := r.db.QueryRowxContext(ctx, query, secret.UserID, secret.Title, secret.Metadata, secret.SecretType, secret.Payload)
	err := result.Scan(&newSecretID)
	if err != nil {
		return 0, err
	}

	return newSecretID, nil
}

// Ensure secret exists and update secret (in one transaction)
func (r SecretsRepository) Update(ctx context.Context, secret *models.Secret) error {
	return runInTx(r.db, func(tx *sqlx.Tx) error {
		err := tx.QueryRowxContext(ctx, "SELECT 1 FROM secrets WHERE id = $1 FOR UPDATE", secret.ID).Scan(new(int))
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("secret with ID %d not found: %w", secret.ID, err)
			}
			return err
		}

		sql := `UPDATE secrets SET updated_at = $1, title = $2, metadata = $3, secret_type = $4, payload = $5 WHERE id = $6;`
		_, err = tx.ExecContext(ctx, sql,
			secret.UpdatedAt,
			secret.Title,
			secret.Metadata,
			secret.SecretType,
			secret.Payload,
			secret.ID,
		)

		return err
	})
}

func (r SecretsRepository) Delete(ctx context.Context, secretID uint64, userID uint64) error {
	query := `DELETE FROM secrets WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, secretID, userID)

	return err
}

func (r SecretsRepository) Pong() {
	fmt.Println("alive")
}

func runInTx(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
