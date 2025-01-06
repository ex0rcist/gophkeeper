package postgres

import (
	"context"
	"database/sql"
	"errors"

	"gophkeeper/internal/server/entities"
	"gophkeeper/internal/server/repository"
	strg "gophkeeper/internal/server/storage/postgres"
	"gophkeeper/pkg/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
)

var _ repository.UsersRepository = UsersRepository{}

// User repository using PostgreSQL
type UsersRepository struct {
	db *sqlx.DB
}

type UsersRepositoryDependencies struct {
	dig.In
	PostgresConn *strg.PostgresConn
}

// Create new postgresql user repository
func NewUsersRepository(deps UsersRepositoryDependencies) *UsersRepository {
	return &UsersRepository{
		db: deps.PostgresConn.DB,
	}
}

// Create new User
func (r UsersRepository) Create(ctx context.Context, user models.User) (int, error) {
	var newUserID int

	result := r.db.QueryRowContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		user.Login,
		user.Password,
	)

	err := result.Scan(&newUserID)
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

// Get User by ID
func (r UsersRepository) GetUserByID(ctx context.Context, ID int) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowxContext(ctx, "SELECT id, login, created_at, password FROM users WHERE id = $1", ID).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrUserNotFound
	}

	return &user, err
}

// Get User by login
func (r UsersRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowxContext(ctx, "SELECT id, login, created_at, password FROM users WHERE login = $1", login).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, entities.ErrUserNotFound
	}

	return &user, err
}
