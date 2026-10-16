INSERT INTO roles (name)
SELECT 'user'
WHERE NOT EXISTS (
    SELECT 1
    FROM roles
    WHERE LOWER(name) = 'user' AND deleted_at IS NULL
);

CREATE TABLE auth_credentials (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE auth_tokens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX auth_tokens_user_id_index ON auth_tokens (user_id);
CREATE INDEX auth_tokens_expires_at_index ON auth_tokens (expires_at);
