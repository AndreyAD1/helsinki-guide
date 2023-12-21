package specifications

import (
	"fmt"
	"strings"
)

type BuildingSpecificationByAlikeAddress struct {
	AddressPrefix string
	Limit         int
	Offset        int
}

func NewBuildingSpecificationByAlikeAddress(
	prefix string,
	limit,
	offset int,
) Specification {
	return &BuildingSpecificationByAlikeAddress{prefix, limit, offset}
}

func (b *BuildingSpecificationByAlikeAddress) ToSQL() (string, map[string]any) {
	queryTemplate := `SELECT *
	FROM (SELECT * FROM buildings WHERE deleted_at IS NULL) AS buildings
	JOIN addresses ON 
	buildings.address_id = addresses.id WHERE lower(street_address) 
	LIKE @search_pattern
	ORDER BY lower(street_address) LIMIT @limit OFFSET @offset;`
	queryArgs := map[string]any{
		"search_pattern": strings.ToLower(b.AddressPrefix) + "%",
		"limit":          b.Limit,
		"offset":         b.Offset,
	}
	return queryTemplate, queryArgs
}

type BuildingSpecificationByAddress struct {
	Address string
}

func NewBuildingSpecificationByAddress(address string) Specification {
	return &BuildingSpecificationByAddress{address}
}

func (b *BuildingSpecificationByAddress) ToSQL() (string, map[string]any) {
	queryTemplate := `SELECT * FROM 
	(SELECT * FROM buildings WHERE deleted_at IS NULL) AS buildings 
	JOIN addresses ON 
	buildings.address_id = addresses.id WHERE lower(street_address) LIKE @address
	ORDER BY name_fi, name_en, name_ru;`
	return queryTemplate, map[string]any{"address": strings.ToLower(b.Address)}
}

type BuildingSpecificationNearest struct {
	DistanceMeters int
	Latitude       string
	Longitude      string
	Limit          int
	Offset         int
}

func NewBuildingSpecificationNearest(
	distanceMeters int,
	latitude,
	longitude float64,
	limit,
	offset int,
) Specification {
	lat := fmt.Sprintf("%.2f", latitude)
	lon := fmt.Sprintf("%.2f", longitude)
	return &BuildingSpecificationNearest{distanceMeters, lat, lon, limit, offset}
}

func (b *BuildingSpecificationNearest) ToSQL() (string, map[string]any) {
	queryTemplate := `WITH nearest_buildings AS (
		SELECT * FROM 
		(SELECT * FROM buildings WHERE deleted_at IS NULL) AS buildings
		WHERE earth_box(ll_to_earth(@latitude, @longitude), @distance) @> 
		ll_to_earth(latitude_wgs84, longitude_wgs84) 
		AND
		earth_distance(
			ll_to_earth(@latitude, @longitude), 
			ll_to_earth(latitude_wgs84, longitude_wgs84)
		) <= @distance
		ORDER BY (
			earth_distance(
				ll_to_earth(@latitude, @longitude),
				ll_to_earth(latitude_wgs84, longitude_wgs84)
			)
		)
	)
	SELECT *
	FROM nearest_buildings
	JOIN addresses ON nearest_buildings.address_id = addresses.id 
	LIMIT @limit OFFSET @offset;`
	args := map[string]any{
		"distance":  b.DistanceMeters,
		"latitude":  b.Latitude,
		"longitude": b.Longitude,
		"limit":     b.Limit,
		"offset":    b.Offset,
	}
	return queryTemplate, args
}
