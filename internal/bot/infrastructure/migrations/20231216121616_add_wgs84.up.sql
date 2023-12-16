BEGIN;

ALTER TABLE buildings 
    ADD COLUMN latitude_wgs84 double precision,
    ADD COLUMN longitude_wgs84 double precision;

COMMIT;