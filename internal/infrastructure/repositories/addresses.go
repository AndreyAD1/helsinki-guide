package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository interface {
	GetAllAdresses(context.Context) ([]string, error)
}

type AddressStorage struct {
	dbPool *pgxpool.Pool
}

func NewAddressRepo(dbPool *pgxpool.Pool) AddressRepository {
	return AddressStorage{dbPool}
}

func (ms AddressStorage) GetAllAdresses(ctx context.Context) ([]string, error) {
	return []string{}, nil
}
