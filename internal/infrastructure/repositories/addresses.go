package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	i "github.com/AndreyAD1/helsinki-guide/internal"
)

type AddressRepository interface {
	GetAddresses(ctx context.Context, limit, offset int) ([]i.Address, error)
	GetAlikeAddresses(context.Context, string) ([]i.Address, error)
}

type AddressStorage struct {
	dbPool *pgxpool.Pool
}

func NewAddressRepo(dbPool *pgxpool.Pool) *AddressStorage {
	return &AddressStorage{dbPool}
}

func (as *AddressStorage) GetAddresses(ctx context.Context, limit, offset int) ([]i.Address, error) {
	query := fmt.Sprintf(
		"SELECT * FROM addresses ORDER BY street_address LIMIT %v OFFSET %v",
		limit,
		offset,
	)
	rows, err := as.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var addresses []i.Address
	for rows.Next() {
		var address i.Address
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

func (as *AddressStorage) GetAlikeAddresses(ctx context.Context, like string) ([]i.Address, error) {
	return nil, ErrNotImplemented
}
