package postgres

import (
	"gophkeeper/internal/server/config"
	"gophkeeper/internal/server/entities"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
)

/* PostgresDSN */
type PostgresDSN entities.SecretConnURI

type PostgresDSNDependencies struct {
	config.Dependency
}

func NewPostgresDSN(deps PostgresDSNDependencies) PostgresDSN {
	return PostgresDSN(deps.Config.PostgresDSN)
}

/* PostgresConn */
type PostgresConn struct {
	DB  *sqlx.DB
	Err error
	DSN PostgresDSN
}

type PostgresConnDependencies struct {
	dig.In
	DSN PostgresDSN
}

func NewPostgresConn(deps PostgresConnDependencies) *PostgresConn {
	conn := &PostgresConn{}
	conn.DB, conn.Err = sqlx.Open("pgx", string(deps.DSN))
	conn.DSN = deps.DSN

	if conn.Err != nil {
		return conn
	}

	if err := conn.DB.Ping(); err != nil {
		conn.Err = err
	}

	return conn
}
