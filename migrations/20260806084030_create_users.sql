-- +goose Up

CREATE TABLE users (

          id BIGSERIAL PRIMARY KEY,

          name VARCHAR(100) NOT NULL,

          email VARCHAR(255) NOT NULL UNIQUE,

          password_hash TEXT NOT NULL,

          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

          deleted_at TIMESTAMPTZ NULL

);

CREATE INDEX idx_users_email
    ON users(email);

CREATE INDEX idx_users_deleted_at
    ON users(deleted_at);

-- +goose Down

DROP TABLE users;