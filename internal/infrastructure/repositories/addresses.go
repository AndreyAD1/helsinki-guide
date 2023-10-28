package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Address struct {
	ID              int64
	StreetAddress   string
	NeighbourhoodID int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type AddressRepository interface {
	GetAdresses(ctx context.Context, limit, offset int) ([]Address, error)
	GetAlikeAddresses(context.Context, string) ([]Address, error)
}

type AddressStorage struct {
	dbPool *pgxpool.Pool
}

func NewAddressRepo(dbPool *pgxpool.Pool) *AddressStorage {
	return &AddressStorage{dbPool}
}

func (ms *AddressStorage) GetAdresses(ctx context.Context, limit, offset int) ([]Address, error) {
	return []Address{}, nil
}

func (ms *AddressStorage) GetAlikeAddresses(ctx context.Context, like string) ([]Address, error) {
	return nil, ErrNotImplemented
}
