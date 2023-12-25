package repositories

import "context"

type BuildingRepository interface {
	Add(context.Context, Building) (*Building, error)
	Remove(context.Context, Building) error
	Update(context.Context, Building) (*Building, error)
	Query(context.Context, Specification) ([]Building, error)
}

type NeighbourhoodRepository interface {
	Add(context.Context, Neighbourhood) (*Neighbourhood, error)
	Remove(context.Context, Neighbourhood) error
	Update(context.Context, Neighbourhood) (*Neighbourhood, error)
	Query(context.Context, Specification) ([]Neighbourhood, error)
}

type ActorRepository interface {
	Add(context.Context, Actor) (*Actor, error)
	Remove(context.Context, Actor) error
	Update(context.Context, Actor) (*Actor, error)
	Query(context.Context, Specification) ([]Actor, error)
}

type Specification interface {
	ToSQL() (string, map[string]any)
}
