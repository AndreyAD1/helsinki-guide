DROP TABLE users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    telegram_id bigint UNIQUE NOT NULL,
    language varchar,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
)