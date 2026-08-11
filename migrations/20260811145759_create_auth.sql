-- +goose Up

CREATE TABLE auth_sessions (
                               id BIGSERIAL PRIMARY KEY,

                               user_id BIGINT NOT NULL,

                               refresh_token_hash VARCHAR(255) NOT NULL UNIQUE,

                               user_agent VARCHAR(500),
                               ip_address VARCHAR(45),

                               expires_at TIMESTAMPTZ NOT NULL,
                               revoked_at TIMESTAMPTZ,

                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                               CONSTRAINT fk_auth_sessions_user
                                   FOREIGN KEY (user_id)
                                       REFERENCES users(id)
                                       ON DELETE CASCADE
);

CREATE INDEX idx_auth_sessions_user_id
    ON auth_sessions(user_id);

CREATE INDEX idx_auth_sessions_expires_at
    ON auth_sessions(expires_at);


CREATE TABLE oauth_accounts (
                                id BIGSERIAL PRIMARY KEY,

                                user_id BIGINT NOT NULL,

                                provider VARCHAR(50) NOT NULL,
                                provider_user_id VARCHAR(255) NOT NULL,

                                email VARCHAR(255),

                                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                CONSTRAINT fk_oauth_accounts_user
                                    FOREIGN KEY (user_id)
                                        REFERENCES users(id)
                                        ON DELETE CASCADE,

                                CONSTRAINT uq_oauth_provider_user
                                    UNIQUE (provider, provider_user_id)
);

CREATE INDEX idx_oauth_accounts_user_id
    ON oauth_accounts(user_id);


-- +goose Down

DROP TABLE IF EXISTS oauth_accounts;
DROP TABLE IF EXISTS auth_sessions;