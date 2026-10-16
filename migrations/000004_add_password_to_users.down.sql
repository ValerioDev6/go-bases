CREATE TABLE auth_credentials (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO auth_credentials (user_id, password_hash)
SELECT id, password
FROM users
WHERE password IS NOT NULL;

ALTER TABLE users DROP COLUMN password;
