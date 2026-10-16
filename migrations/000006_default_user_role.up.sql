DO $$
DECLARE
    user_role_id BIGINT;
BEGIN
    SELECT id INTO user_role_id
    FROM roles
    WHERE LOWER(name) = 'user' AND deleted_at IS NULL
    LIMIT 1;

    EXECUTE format(
        'ALTER TABLE users ALTER COLUMN role_id SET DEFAULT %s',
        user_role_id
    );
END $$;
