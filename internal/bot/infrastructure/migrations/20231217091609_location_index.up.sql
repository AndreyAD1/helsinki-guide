CREATE INDEX CONCURRENTLY geolocation_index ON buildings
USING gist 
(ll_to_earth(latitude_wgs84, longitude_wgs84));
