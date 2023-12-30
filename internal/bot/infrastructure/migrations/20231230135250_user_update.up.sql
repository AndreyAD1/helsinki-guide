ALTER TABLE users
    DROP COLUMN role,
    DROP COLUMN username,
    ADD COLUMN telegram_id bigint UNIQUE NOT NULL;