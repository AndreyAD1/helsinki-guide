package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	s "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	i "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	l "github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
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
	deleted_at FROM neighbourhoods WHERE name = $1 AND `
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

	municipality := ""
	if neighbourhood.Municipality != nil {
		municipality = *neighbourhood.Municipality
	}
	logMsg := fmt.Sprintf(
		"neigbourhood duplication: '%v-%v'",
		neighbourhood.Name,
		municipality,
	)
	slog.DebugContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
	var row pgx.Row
	if neighbourhood.Municipality == nil {
		row = n.dbPool.QueryRow(
			ctx,
			selectQuery+"municipality is NULL;",
			neighbourhood.Name,
		)
	} else {
		row = n.dbPool.QueryRow(
			ctx,
			selectQuery+"municipality = $2;",
			neighbourhood.Name,
			neighbourhood.Municipality,
		)
	}

	err = row.Scan(
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

func (n *neighbourhoodStorage) Query(ctx context.Context, spec s.Specification) ([]i.Neighbourhood, error) {
	query, queryArgs := spec.ToSQL()
	slog.DebugContext(ctx, fmt.Sprintf("send the query %v: %v", query, queryArgs))
	rows, err := n.dbPool.Query(ctx, query, pgx.NamedArgs(queryArgs))
	if err != nil {
		logMsg := fmt.Sprintf("a query error: '%v'", query)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	defer rows.Close()
	var neigbourhoods []i.Neighbourhood
	for rows.Next() {
		var neigbourhood i.Neighbourhood
		if err := rows.Scan(
			&neigbourhood.ID,
			&neigbourhood.Name,
			&neigbourhood.Municipality,
			&neigbourhood.CreatedAt,
			&neigbourhood.UpdatedAt,
			&neigbourhood.DeletedAt,
		); err != nil {
			msg := fmt.Sprintf(
				"can not scan an actor from a query result: %v: %v",
				query,
				queryArgs,
			)
			slog.ErrorContext(ctx, msg, slog.Any(logger.ErrorKey, err))
			return nil, err
		}
		neigbourhoods = append(neigbourhoods, neigbourhood)
	}
	return neigbourhoods, nil
}
