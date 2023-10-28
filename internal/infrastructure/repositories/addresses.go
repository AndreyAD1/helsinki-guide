package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Address struct {
	ID              int64
	StreetAddress   string
	NeighbourhoodID *int64
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
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

func (as *AddressStorage) GetAdresses(ctx context.Context, limit, offset int) ([]Address, error) {
	query := fmt.Sprintf(
		"SELECT * FROM addresses ORDER BY street_address LIMIT %v OFFSET %v",
		limit,
		offset,
	)
	rows, err := as.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var addresses []Address
	for rows.Next() {
		var address Address
		if err := rows.Scan(
			&address.ID, 
			&address.StreetAddress, 
			&address.NeighbourhoodID, 
			&address.CreatedAt, 
			&address.UpdatedAt, 
			&address.DeletedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (as *AddressStorage) GetAlikeAddresses(ctx context.Context, like string) ([]Address, error) {
	return nil, ErrNotImplemented
}
