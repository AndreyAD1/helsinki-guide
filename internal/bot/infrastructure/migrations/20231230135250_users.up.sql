BEGIN;
DROP TABLE users;
CREATE TYPE language AS ENUM('fi', 'en', 'ru');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id bigint UNIQUE NOT NULL,
    language language,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

COMMIT;