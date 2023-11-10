package repositories

import (
	"context"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NeighbourhoodRepository interface {
	Add(context.Context, i.Neighbourhood) (*i.Neighbourhood, error)
	Remove(context.Context, i.Neighbourhood) error
	Update(context.Context, i.Neighbourhood) (*i.Neighbourhood, error)
	Query(context.Context, s.Specification) ([]i.Neighbourhood, error)
}

type neighbourhoodStorage struct {
	dbPool *pgxpool.Pool
}

func NewNeighbourhoodRepo(dbPool *pgxpool.Pool) NeighbourhoodRepository {
	return &neighbourhoodStorage{dbPool}
}

func (n *neighbourhoodStorage) Add(ctx context.Context, actor i.Neighbourhood) (*i.Neighbourhood, error) {
	return nil, ErrNotImplemented
}

func (n *neighbourhoodStorage) Remove(ctx context.Context, actor i.Neighbourhood) error {
	return ErrNotImplemented
}

func (n *neighbourhoodStorage) Update(ctx context.Context, actor i.Neighbourhood) (*i.Neighbourhood, error) {
	return nil, ErrNotImplemented
}

func (n *neighbourhoodStorage) Query(ctx context.Context, spec s.Specification) ([]i.Neighbourhood, error) {
	return nil, ErrNotImplemented
}
