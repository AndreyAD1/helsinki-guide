package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type settingsStorage struct {
	dbPool *pgxpool.Pool
}

func NewSettingRepo(dbPool *pgxpool.Pool) *settingsStorage {
	return &settingsStorage{dbPool}
}

func (s *settingsStorage) Add(ctx context.Context, settings Settings) (*Settings, error) {
	return nil, ErrNotImplemented
}

func (s *settingsStorage) Remove(ctx context.Context, settings Settings) error {
	return ErrNotImplemented
}

func (s *settingsStorage) Update(ctx context.Context, settings Settings) (*Settings, error) {
	return nil, ErrNotImplemented
}

func (s *settingsStorage) Query(ctx context.Context, spec Specification) ([]Settings, error) {
	return nil, ErrNotImplemented
}
