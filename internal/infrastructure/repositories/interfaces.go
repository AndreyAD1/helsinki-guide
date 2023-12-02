package repositories

import (
	"context"

	i "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories/types"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories/specifications"
)

type BuildingRepository interface {
	Add(context.Context, i.Building) (*i.Building, error)
	Remove(context.Context, i.Building) error
	Update(context.Context, i.Building) (*i.Building, error)
	Query(context.Context, s.Specification) ([]i.Building, error)
}

type NeighbourhoodRepository interface {
	Add(context.Context, i.Neighbourhood) (*i.Neighbourhood, error)
	Remove(context.Context, i.Neighbourhood) error
	Update(context.Context, i.Neighbourhood) (*i.Neighbourhood, error)
	Query(context.Context, s.Specification) ([]i.Neighbourhood, error)
}

type ActorRepository interface {
	Add(context.Context, i.Actor) (*i.Actor, error)
	Remove(context.Context, i.Actor) error
	Update(context.Context, i.Actor) (*i.Actor, error)
	Query(context.Context, s.Specification) ([]i.Actor, error)
}
