BEGIN;

ALTER TABLE "current_uses" DROP CONSTRAINT use_type_current_use;

ALTER TABLE "current_uses" DROP CONSTRAINT building_current_use;

ALTER TABLE "initial_uses" DROP CONSTRAINT use_type_initial_use;

ALTER TABLE "initial_uses" DROP CONSTRAINT building_initial_use;

ALTER TABLE "building_contractors" DROP CONSTRAINT actor_contractor;

ALTER TABLE "building_contractors" DROP CONSTRAINT building_contractor;

ALTER TABLE "building_builders" DROP CONSTRAINT actor_builder;

ALTER TABLE "building_builders" DROP CONSTRAINT building_builder;

ALTER TABLE "building_designers" DROP CONSTRAINT actor_designer;

ALTER TABLE "building_designers" DROP CONSTRAINT building_designer;

ALTER TABLE "building_authors" DROP CONSTRAINT actor_author;

ALTER TABLE "building_authors" DROP CONSTRAINT building_author;

ALTER TABLE "addresses" DROP CONSTRAINT address_neighbourhood;

ALTER TABLE "buildings" DROP CONSTRAINT address_building;

DROP TABLE current_uses, initial_uses, use_types, building_contractors,
building_builders, building_designers, building_authors, actors,
neighbourhoods, addresses, buildings, users;

COMMIT;