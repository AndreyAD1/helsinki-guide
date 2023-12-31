BEGIN;
DROP TABLE users;

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "created_at" timestamp with time zone NOT NULL DEFAULT now(),
  "updated_at" timestamp with time zone,
  "deleted_at" timestamp with time zone,
  "username" varchar NOT NULL,
  "role" varchar,
  "language" varchar
);
COMMIT;