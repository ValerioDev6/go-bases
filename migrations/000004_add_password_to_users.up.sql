ALTER TABLE users ADD COLUMN password TEXT;

UPDATE users AS u
SET password = a.password_hash
FROM auth_credentials AS a
WHERE a.user_id = u.id;

DROP TABLE auth_credentials;
