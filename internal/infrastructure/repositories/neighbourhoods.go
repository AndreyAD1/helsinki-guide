package repositories

import (
	"context"
	"errors"
	"log/slog"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	l "github.com/AndreyAD1/helsinki-guide/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type neighbourhoodStorage struct {
	dbPool *pgxpool.Pool
}

func NewNeighbourhoodRepo(dbPool *pgxpool.Pool) NeighbourhoodRepository {
	return &neighbourhoodStorage{dbPool}
}

func (n *neighbourhoodStorage) Add(
	ctx context.Context, 
	neighbourhood i.Neighbourhood,
) (*i.Neighbourhood, error) {
	query := `INSERT INTO neighbourhoods (name, municipality)
	VALUES ($1, $2) RETURNING id;`
	var id int64
    err := n.dbPool.QueryRow(
		ctx, 
		query, 
		neighbourhood.Name, 
		neighbourhood.Municipality, 
	).Scan(&id)
    if err != nil {
        var pgxError *pgconn.PgError
        if errors.As(err, &pgxError) {
            if pgxError.Code == pgerrcode.UniqueViolation {
                return nil, ErrDuplicate
            }
        }
		slog.WarnContext(ctx, "unexpected DB error", slog.Any(l.ErrorKey, err))
        return nil, err
    }
    neighbourhood.ID = id

	return &neighbourhood, nil
}

func (n *neighbourhoodStorage) Remove(ctx context.Context, neighbourhood i.Neighbourhood) error {
	return ErrNotImplemented
}

func (n *neighbourhoodStorage) Update(ctx context.Context, neighbourhood i.Neighbourhood) (*i.Neighbourhood, error) {
	return nil, ErrNotImplemented
}

func (n *neighbourhoodStorage) Query(ctx context.Context, neighbourhood s.Specification) ([]i.Neighbourhood, error) {
	return nil, ErrNotImplemented
}
