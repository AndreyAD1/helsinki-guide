package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/logger"
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
	selectQuery := `SELECT id, name, municipality, created_at, updated_at,
	deleted_at FROM neighbourhoods WHERE name = $1 and municipality = $2;`
	insertQuery := `INSERT INTO neighbourhoods (name, municipality)
	VALUES ($1, $2) RETURNING id, name, municipality, created_at, updated_at,
	deleted_at;`
	var saved i.Neighbourhood
	err := n.dbPool.QueryRow(
		ctx,
		insertQuery,
		neighbourhood.Name,
		neighbourhood.Municipality,
	).Scan(
		&saved.ID, 
		&saved.Name, 
		&saved.Municipality,
		&saved.CreatedAt,
		&saved.UpdatedAt,
		&saved.DeletedAt,
	)
	if err == nil {
		return &saved, nil
	}
	var pgxError *pgconn.PgError
	if !errors.As(err, &pgxError) {
		slog.WarnContext(ctx, "unexpected DB error", slog.Any(l.ErrorKey, err))
		return nil, err
	}
	unexpectedMsg := fmt.Sprintf(
		"unexpected DB error for a neighbourhood '%v-%v'", 
		neighbourhood.Name,
		neighbourhood.Municipality,
	)
	if pgxError.Code != pgerrcode.UniqueViolation {
		slog.WarnContext(ctx, unexpectedMsg, slog.Any(logger.ErrorKey, err))
		return nil, err
	}

	logMsg := fmt.Sprintf(
		"neigbourhood duplication: '%v-%v'", 
		neighbourhood.Name,
		neighbourhood.Municipality,
	)
	slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
	err = n.dbPool.QueryRow(
		ctx,
		selectQuery,
		neighbourhood.Name,
		neighbourhood.Municipality,
	).Scan(
		&saved.ID,
		&saved.Name,
		&saved.Municipality,
		&saved.CreatedAt,
		&saved.UpdatedAt,
		&saved.DeletedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, unexpectedMsg, slog.Any(logger.ErrorKey, err))
		return nil, err
	}
	return &saved, ErrDuplicate
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
