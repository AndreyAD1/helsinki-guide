package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userStorage struct {
	dbPool *pgxpool.Pool
}

func NewUserRepo(dbPool *pgxpool.Pool) *userStorage {
	return &userStorage{dbPool}
}

func (s *userStorage) Add(ctx context.Context, user User) (*User, error) {
	return nil, ErrNotImplemented
}

func (s *userStorage) Remove(ctx context.Context, user User) error {
	return ErrNotImplemented
}

func (s *userStorage) Update(ctx context.Context, user User) (*User, error) {
	return nil, ErrNotImplemented
}

func (s *userStorage) Query(ctx context.Context, spec Specification) ([]User, error) {
	return nil, ErrNotImplemented
}
