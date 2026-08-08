-- +goose Up

CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,

                       email VARCHAR(255) NOT NULL UNIQUE,
                       username VARCHAR(100) NOT NULL,
                       password_hash TEXT NOT NULL,

                       role VARCHAR(20) NOT NULL DEFAULT 'user',
                       status VARCHAR(20) NOT NULL DEFAULT 'active',
                       email_verified BOOLEAN NOT NULL DEFAULT FALSE,

                       last_login_at TIMESTAMPTZ NULL,

                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);


-- +goose Down

DROP TABLE users;