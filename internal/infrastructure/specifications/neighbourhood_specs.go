package specifications

import i "github.com/AndreyAD1/helsinki-guide/internal"

type NeighbourhoodSpecificationByName struct {
	neigbourhood i.Neighbourhood
}

func NewNeighbourhoodSpecificationByName(n i.Neighbourhood) Specification {
	return &NeighbourhoodSpecificationByName{n}
}

func (a *NeighbourhoodSpecificationByName) ToSQL() (string, map[string]any) {
	query := `SELECT id, name, municipality, created_at, updated_at, 
	deleted_at FROM neighbourhoods WHERE name = @name AND `
	params := map[string]any{
		"name": a.neigbourhood.Name,
	}
	if a.neigbourhood.Municipality == nil {
		query := query + "municipality is NULL;"
		return query, params
	}

	query = query + "municipality = @municipality;"
	params["municipality"] = a.neigbourhood.Municipality
	return query, params
}

type NeighbourhoodSpecificationAll struct {
	limit int
	offset int
}

func NewNeighbourhoodSpecificationAll(limit, offset int) Specification {
	return &NeighbourhoodSpecificationAll{limit, offset}
}

func (a *NeighbourhoodSpecificationAll) ToSQL() (string, map[string]any) {
	query := `SELECT id, name, municipality, created_at, updated_at, 
	deleted_at FROM neighbourhoods ORDER BY name LIMIT @limit OFFSET @offset;`
	return query, map[string]any{"limit": a.limit, "offset": a.offset}
}
