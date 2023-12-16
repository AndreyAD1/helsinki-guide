BEGIN;

ALTER TABLE buildings
    DROP COLUMN latitude_wgs84,
    DROP COLUMN longitude_wgs84;

COMMIT;