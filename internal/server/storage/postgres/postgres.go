package postgres

import (
	"context"
	"fmt"
	"gophkeeper/internal/server/storage"
	"gophkeeper/internal/server/storage/postgres/migrations"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/dig"
)

const migrationTimeout = 5 * time.Second

var _ storage.ServerStorage = PostgresStorage{}

type PostgresStorageDependencies struct {
	dig.In
	PostgresConn *PostgresConn
}

type PostgresStorage struct {
	db  *sqlx.DB
	dsn PostgresDSN
}

// type dbQueryTracer struct {
// 	logger *zerolog.Logger
// }

// // Logs query
// func (tracer *dbQueryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
// 	tracer.logger.Debug().Msg(fmt.Sprintf("Executing command \"%s\" with args %v", data.SQL, data.Args))
// 	return ctx
// }

// // Does nothing, required by interface
// func (tracer *dbQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
// 	// empty
// }

// PostgresStorage constructor
func NewPostgresStorage(deps PostgresStorageDependencies) (*PostgresStorage, error) {
	fmt.Println("new storage")

	conn := deps.PostgresConn
	if conn.Err != nil {
		return nil, conn.Err
	}

	storage := &PostgresStorage{db: conn.DB}

	// run migrations
	if err := storage.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return storage, nil
}

// Ping storage
func (s PostgresStorage) Ping(ctx context.Context) error {
	return s.db.Ping()
}

// Close storage
func (s PostgresStorage) Close(ctx context.Context) error {
	return s.db.Close()
}

// Stringer for logging
func (s PostgresStorage) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\t\tDSN: %s\n", s.dsn))

	return sb.String()
}

func (s PostgresStorage) GetDB() *sqlx.DB {
	fmt.Println("getDB")
	return s.db
}

// Performs DB migrations
func (s PostgresStorage) migrate() error {
	goose.SetBaseFS(migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), migrationTimeout)
	defer cancel()

	err := goose.RunContext(ctx, "up", s.db.DB, ".")
	if err != nil {
		return err
	}

	return nil
}
