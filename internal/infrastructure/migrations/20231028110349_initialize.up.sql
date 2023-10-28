BEGIN;

CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "username" varchar NOT NULL,
  "role" varchar,
  "language" varchar
);

CREATE TABLE "buildings" (
  "id" integer PRIMARY KEY,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp,
  "code" varchar,
  "name_fi" varchar,
  "name_en" varchar,
  "name_ru" varchar,
  "address_id" integer,
  "construction_start_year" integer,
  "completion_year" integer,
  "complex_fi" varchar,
  "complex_en" varchar,
  "complex_ru" varchar,
  "history_fi" varchar,
  "history_en" varchar,
  "history_ru" varchar,
  "reasoning_fi" varchar,
  "reasoning_en" varchar,
  "reasoning_ru" varchar,
  "protection_status_fi" varchar,
  "protection_status_en" varchar,
  "protection_status_ru" varchar,
  "info_source_fi" varchar,
  "info_source_en" varchar,
  "info_source_ru" varchar,
  "surroundings_fi" varchar,
  "surroundings_en" varchar,
  "surroundings_ru" varchar,
  "foundation_fi" varchar,
  "foundation_en" varchar,
  "foundation_ru" varchar,
  "frame_fi" varchar,
  "frame_en" varchar,
  "frame_ru" varchar,
  "floor_description_fi" varchar,
  "floor_description_en" varchar,
  "floor_description_ru" varchar,
  "facades_fi" varchar,
  "facades_en" varchar,
  "facades_ru" varchar,
  "special_features_fi" varchar,
  "special_features_en" varchar,
  "special_features_ru" varchar,
  "latitude_ETRSGK25" real,
  "longitude_ETRSGK25" real
);

CREATE TABLE "addresses" (
  "id" integer PRIMARY KEY,
  "street_address" varchar,
  "neighbourhood_id" integer,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "neighbourhoods" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "municipality" varchar,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "actors" (
  "id" integer PRIMARY KEY,
  "name" varchar NOT NULL,
  "title_fi" varchar,
  "title_en" varchar,
  "title_ru" varchar,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "building_authors" (
  "building_id" integer,
  "actor_id" integer
);

CREATE TABLE "building_designers" (
  "building_id" integer,
  "actor_id" integer
);

CREATE TABLE "building_builders" (
  "building_id" integer,
  "actor_id" integer
);

CREATE TABLE "building_contractors" (
  "building_id" integer,
  "actor_id" integer
);

CREATE TABLE "use_types" (
  "id" integer PRIMARY KEY,
  "name_fi" varchar UNIQUE,
  "name_en" varchar UNIQUE,
  "name_ru" varchar UNIQUE,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE "initial_uses" (
  "building_id" integer,
  "use_type_id" integer
);

CREATE TABLE "current_uses" (
  "building_id" integer,
  "use_type_id" integer
);

ALTER TABLE "addresses" ADD CONSTRAINT address_building 
FOREIGN KEY ("id") REFERENCES "buildings" ("address_id");

ALTER TABLE "addresses" ADD CONSTRAINT address_neighbourhood 
FOREIGN KEY ("neighbourhood_id") REFERENCES "neighbourhoods" ("id");

ALTER TABLE "building_authors" ADD CONSTRAINT building_author 
FOREIGN KEY ("building_id") REFERENCES "buildings" ("id");

ALTER TABLE "building_authors" ADD CONSTRAINT actor_author 
FOREIGN KEY ("actor_id") REFERENCES "actors" ("id");

ALTER TABLE "building_designers" ADD CONSTRAINT building_designer 
FOREIGN KEY ("building_id") REFERENCES "buildings" ("id");

ALTER TABLE "building_designers" ADD CONSTRAINT actor_designer 
FOREIGN KEY ("actor_id") REFERENCES "actors" ("id");

ALTER TABLE "building_builders" ADD CONSTRAINT building_builder 
FOREIGN KEY ("building_id") REFERENCES "buildings" ("id");

ALTER TABLE "building_builders" ADD CONSTRAINT actor_builder 
FOREIGN KEY ("actor_id") REFERENCES "actors" ("id");

ALTER TABLE "building_contractors" ADD CONSTRAINT building_contractor 
FOREIGN KEY ("building_id") REFERENCES "buildings" ("id");

ALTER TABLE "building_contractors" ADD CONSTRAINT actor_contractor 
FOREIGN KEY ("actor_id") REFERENCES "actors" ("id");

ALTER TABLE "initial_uses" ADD CONSTRAINT building_initial_use 
FOREIGN KEY ("building_id") REFERENCES "buildings" ("id");

ALTER TABLE "initial_uses" ADD CONSTRAINT use_type_initial_use 
FOREIGN KEY ("use_type_id") REFERENCES "use_types" ("id");

ALTER TABLE "current_uses" ADD CONSTRAINT building_current_use 
FOREIGN KEY ("building_id") REFERENCES "buildings" ("id");

ALTER TABLE "current_uses" ADD CONSTRAINT use_type_current_use 
FOREIGN KEY ("use_type_id") REFERENCES "use_types" ("id");

COMMIT;