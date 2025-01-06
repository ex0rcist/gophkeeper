package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/storage/postgres"
	"gophkeeper/pkg/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretsRepository_GetSecret(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSecretsRepository(SecretsRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "metadata", "secret_type", "payload"}).AddRow(1, 1, "Test Title", "{}", "type", "payload")
		mock.ExpectQuery(`SELECT \* FROM secrets WHERE id = \$1 AND user_id = \$2`).WithArgs(1, 1).WillReturnRows(rows)

		secret, err := repo.GetSecret(context.Background(), 1, 1)
		assert.NoError(t, err)
		assert.NotNil(t, secret)
		assert.Equal(t, "Test Title", secret.Title)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM secrets WHERE id = \$1 AND user_id = \$2`).WithArgs(1, 1).WillReturnError(sql.ErrNoRows)

		secret, err := repo.GetSecret(context.Background(), 1, 1)
		assert.Error(t, err)
		assert.Nil(t, secret)
		assert.True(t, errors.Is(err, entities.ErrUserNotFound))
	})
}

func TestSecretsRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSecretsRepository(SecretsRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO secrets \(user_id, title, metadata, secret_type, payload\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
			WithArgs(1, "Test Title", "{}", "credential", []byte("payload")).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repo.Create(context.Background(), &models.Secret{
			UserID:     1,
			Title:      "Test Title",
			Metadata:   "{}",
			SecretType: "credential",
			Payload:    []byte("payload"),
		})

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), id)
	})
}

func TestSecretsRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSecretsRepository(SecretsRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT 1 FROM secrets WHERE id = \$1 FOR UPDATE`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
		mock.ExpectExec(`UPDATE secrets SET updated_at = \$1, title = \$2, metadata = \$3, secret_type = \$4, payload = \$5 WHERE id = \$6`).
			WithArgs(sqlmock.AnyArg(), "Updated Title", "{}", "credential", []byte("new_payload"), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), &models.Secret{
			ID:         1,
			Title:      "Updated Title",
			Metadata:   "{}",
			SecretType: "credential",
			Payload:    []byte("new_payload"),
		})

		assert.NoError(t, err)
	})
}

func TestSecretsRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSecretsRepository(SecretsRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM secrets WHERE id = \$1 AND user_id = \$2`).WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(context.Background(), 1, 1)
		assert.NoError(t, err)
	})
}
