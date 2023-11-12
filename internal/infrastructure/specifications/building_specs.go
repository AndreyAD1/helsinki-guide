package specifications

type Specification interface {
	ToSQL() (string, map[string]any)
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
) *BuildingSpecificationByAlikeAddress {
	return &BuildingSpecificationByAlikeAddress{prefix, limit, offset}
}

func (b *BuildingSpecificationByAlikeAddress) ToSQL() (string, map[string]any) {
	queryTemplate := `SELECT *
	FROM buildings JOIN addresses ON 
	buildings.address_id = addresses.id WHERE street_address 
	ILIKE @search_pattern
	ORDER BY street_address LIMIT @limit OFFSET @offset;`
	queryArgs := map[string]any{
		"search_pattern": b.addressPrefix + "%",
		"limit":          b.limit,
		"offset":         b.offset,
	}
	return queryTemplate, queryArgs
}

type BuildingSpecificationByAddress struct {
	address string
}

func NewBuildingSpecificationByAddress(address string) *BuildingSpecificationByAddress {
	return &BuildingSpecificationByAddress{address}
}

func (b *BuildingSpecificationByAddress) ToSQL() (string, map[string]any) {
	queryTemplate := `SELECT * FROM buildings JOIN addresses ON 
	buildings.address_id = addresses.id WHERE street_address = @address
	ORDER BY name_fi, name_en, name_ru;`
	return queryTemplate, map[string]any{"address": b.address}
}
