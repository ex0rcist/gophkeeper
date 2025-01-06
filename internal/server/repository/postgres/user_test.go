package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/storage/postgres"
	"gophkeeper/pkg/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsersRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUsersRepository(UsersRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
			WithArgs("testuser", "hashedpassword").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		id, err := repo.Create(context.Background(), models.User{
			Login:    "testuser",
			Password: "hashedpassword",
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
	})
}

func TestUsersRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUsersRepository(UsersRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		createdAt := time.Now()
		rows := sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
			AddRow(1, "testuser", createdAt, "hashedpassword")
		mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnRows(rows)

		user, err := repo.GetUserByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Login)
		assert.Equal(t, createdAt, user.CreatedAt)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE id = \$1`).
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, entities.ErrUserNotFound))
	})
}

func TestUsersRepository_GetUserByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewUsersRepository(UsersRepositoryDependencies{
		PostgresConn: &postgres.PostgresConn{DB: sqlxDB},
	})

	t.Run("Success", func(t *testing.T) {
		createdAt := time.Now()
		rows := sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
			AddRow(1, "testuser", createdAt, "hashedpassword")
		mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE login = \$1`).
			WithArgs("testuser").
			WillReturnRows(rows)

		user, err := repo.GetUserByLogin(context.Background(), "testuser")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Login)
		assert.Equal(t, createdAt, user.CreatedAt)
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE login = \$1`).
			WithArgs("testuser").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByLogin(context.Background(), "testuser")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.True(t, errors.Is(err, entities.ErrUserNotFound))
	})
}
