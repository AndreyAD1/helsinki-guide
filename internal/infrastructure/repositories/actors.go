package repositories

import (
	"context"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActorRepository interface {
	Add(context.Context, i.Actor) (*i.Actor, error)
	Remove(context.Context, i.Actor) error
	Update(context.Context, i.Actor) (*i.Actor, error)
	Query(context.Context, s.Specification) ([]i.Actor, error)
}

type actorStorage struct {
	dbPool *pgxpool.Pool
}

func NewActorRepo(dbPool *pgxpool.Pool) ActorRepository {
	return &actorStorage{dbPool}
}

func (a *actorStorage) Add(ctx context.Context, actor i.Actor) (*i.Actor, error) {
	return nil, ErrNotImplemented
}

func (b *actorStorage) Remove(ctx context.Context, actor i.Actor) error {
	return ErrNotImplemented
}

func (b *actorStorage) Update(ctx context.Context, actor i.Actor) (*i.Actor, error) {
	return nil, ErrNotImplemented
}

func (b *actorStorage) Query(ctx context.Context, spec s.Specification) ([]i.Actor, error) {
	return nil, ErrNotImplemented
}
