package specifications

import (
	"fmt"
	"strings"
)

type BuildingSpecificationAll struct {
	limit int
	offset int
}

func NewBuildingSpecificationAll(limit, offset int) Specification {
	return &BuildingSpecificationAll{limit, offset}
}

func (b *BuildingSpecificationAll) ToSQL() (string, map[string]any) {
	queryTemplate := `SELECT *
	FROM (SELECT * FROM buildings WHERE deleted_at IS NULL) AS buildings
	JOIN addresses ON buildings.address_id = addresses.id 
	ORDER BY buildings.id LIMIT @limit OFFSET @offset;`
	queryArgs := map[string]any{
		"limit":          b.limit,
		"offset":         b.offset,
	}
	return queryTemplate, queryArgs
}

type BuildingSpecificationByAlikeAddress struct {
	addressPrefix string
	limit         int
	offset        int
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
		"search_pattern": strings.ToLower(b.addressPrefix) + "%",
		"limit":          b.limit,
		"offset":         b.offset,
	}
	return queryTemplate, queryArgs
}

func AlikeAddressSpecIsEqual(addressPrefix string, limit, offset int) func(s *BuildingSpecificationByAlikeAddress) bool {
	return func(s *BuildingSpecificationByAlikeAddress) bool {
		addressMatch := s.addressPrefix == addressPrefix
		limitMatch := s.limit == limit
		offsetMatch := s.offset == offset
		return addressMatch && limitMatch && offsetMatch
	}
}

type BuildingSpecificationByAddress struct {
	address string
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
	return queryTemplate, map[string]any{"address": strings.ToLower(b.address)}
}

func BuildingByAddressIsEqual(address string) func(s *BuildingSpecificationByAddress) bool {
	return func(s *BuildingSpecificationByAddress) bool {
		return address == s.address
	}
}

type BuildingSpecificationNearest struct {
	distanceMeters int
	latitude       string
	longitude      string
	limit          int
	offset         int
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
		"distance":  b.distanceMeters,
		"latitude":  b.latitude,
		"longitude": b.longitude,
		"limit":     b.limit,
		"offset":    b.offset,
	}
	return queryTemplate, args
}

func NearestSpecIsEqual(
	distanceMeters int,
	latitude,
	longitude float64,
	limit,
	offset int,
) func(s *BuildingSpecificationNearest) bool {
	return func(s *BuildingSpecificationNearest) bool {
		distanceMatch := s.distanceMeters == distanceMeters
		latMatch := s.latitude == fmt.Sprintf("%.2f", latitude)
		lonMatch := s.longitude == fmt.Sprintf("%.2f", longitude)
		limitMatch := s.limit == limit
		offsetMatch := s.offset == offset
		return distanceMatch && latMatch && lonMatch && limitMatch && offsetMatch
	}
}
