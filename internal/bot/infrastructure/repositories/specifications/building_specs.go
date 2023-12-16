package specifications

import "strings"

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
	FROM buildings JOIN addresses ON 
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
	queryTemplate := `SELECT * FROM buildings JOIN addresses ON 
	buildings.address_id = addresses.id WHERE lower(street_address) LIKE @address
	ORDER BY name_fi, name_en, name_ru;`
	return queryTemplate, map[string]any{"address": strings.ToLower(b.Address)}
}
