package specifications

import "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"

type ActorSpecificationByBuilding struct {
	BuildingID int64
}

func NewActorSpecificationByBuilding(buildingID int64) *ActorSpecificationByBuilding {
	return &ActorSpecificationByBuilding{buildingID}
}

func (a *ActorSpecificationByBuilding) ToSQL() (string, map[string]any) {
	query := `SELECT id, name, title_fi, title_en, title_ru, created_at,
	updated_at, deleted_at FROM actors JOIN building_authors ON id = actor_id
	WHERE building_id = @building_id;`
	return query, map[string]any{"building_id": a.BuildingID}
}

type ActorSpecificationByName struct {
	actor types.Actor
}

func NewActorSpecificationByName(a types.Actor) Specification {
	return &ActorSpecificationByName{a}
}

func (a *ActorSpecificationByName) ToSQL() (string, map[string]any) {
	query := `SELECT id, name, title_fi, title_en, title_ru, created_at,
	updated_at, deleted_at FROM actors WHERE name = @name;`
	return query, map[string]any{"name": a.actor.Name}
}

type ActorSpecificationAll struct {
	limit  int
	offset int
}

func NewActorSpecificationAll(limit, offset int) Specification {
	return &ActorSpecificationAll{limit, offset}
}

func (a *ActorSpecificationAll) ToSQL() (string, map[string]any) {
	query := `SELECT id, name, title_fi, title_en, title_ru, created_at,
	updated_at, deleted_at FROM actors ORDER BY name 
	LIMIT @limit OFFSET @offset`
	return query, map[string]any{"limit": a.limit, "offset": a.offset}
}
