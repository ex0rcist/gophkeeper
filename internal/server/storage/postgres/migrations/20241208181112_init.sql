-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    login varchar(100) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);
CREATE INDEX login_idx ON users (login);

CREATE TYPE secret_type AS ENUM ('credential', 'text', 'blob', 'card');
CREATE TABLE IF NOT EXISTS secrets (
    id serial PRIMARY KEY,
    user_id integer,
    title varchar(255) NOT NULL,
    metadata TEXT,
    secret_type secret_type NOT NULL,
    payload bytea NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE secrets;
-- +goose StatementEnd
