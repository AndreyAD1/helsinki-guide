package specifications

import (
	"fmt"
)

type BuildingSpecification interface {
	ToSQL() string
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

func (b *BuildingSpecificationByAlikeAddress) ToSQL() string {
	queryTemplate := `SELECT *
	FROM buildings JOIN addresses ON 
	buildings.address_id = addresses.id WHERE street_address ILIKE '%v%%'
	ORDER BY street_address LIMIT %v OFFSET %v;`

	query := fmt.Sprintf(
		queryTemplate,
		b.addressPrefix,
		b.limit,
		b.offset,
	)
	return query
}

type BuildingSpecificationByAddress struct {
	address string
}

func NewBuildingSpecificationByAddress(address string) *BuildingSpecificationByAddress {
	return &BuildingSpecificationByAddress{address}
}

func (b *BuildingSpecificationByAddress) ToSQL() string {
	queryTemplate := `SELECT *
	FROM buildings JOIN addresses ON 
	buildings.address_id = addresses.id WHERE street_address = '%s'`
	query := fmt.Sprintf(queryTemplate, b.address)
	return query
}
