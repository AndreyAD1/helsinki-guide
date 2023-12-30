ALTER TABLE users
    DROP column telegram_id,
    ADD column username varchar NOT NULL,
    ADD column role varchar;