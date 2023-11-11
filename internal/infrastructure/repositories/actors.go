package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type actorStorage struct {
	dbPool *pgxpool.Pool
}

func NewActorRepo(dbPool *pgxpool.Pool) ActorRepository {
	return &actorStorage{dbPool}
}

func (a *actorStorage) Add(ctx context.Context, actor i.Actor) (*i.Actor, error) {
	query := `INSERT INTO actors (name, title_fi, title_en, title_ru)
	VALUES ($1, $2, $3, $4) RETURNING id;`
	var id int64
    err := a.dbPool.QueryRow(
		ctx, 
		query, 
		actor.Name, 
		actor.TitleFi,
		actor.TitleEn,
		actor.TitleRu, 
	).Scan(&id)
    if err != nil {
        var pgxError *pgconn.PgError
        if errors.As(err, &pgxError) {
            if pgxError.Code == pgerrcode.UniqueViolation {
				logMsg := fmt.Sprintf("actor duplication: %v", actor)
				slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
                return nil, ErrDuplicate
            }
        }
		logMsg := fmt.Sprintf("unexpected DB error for an actor %v", actor.Name)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
        return nil, err
    }
    actor.ID = id

	return &actor, nil
}

func (a *actorStorage) Remove(ctx context.Context, actor i.Actor) error {
	return ErrNotImplemented
}

func (a *actorStorage) Update(ctx context.Context, actor i.Actor) (*i.Actor, error) {
	return nil, ErrNotImplemented
}

func (a *actorStorage) Query(ctx context.Context, spec s.Specification) ([]i.Actor, error) {
	query, queryArgs := spec.ToSQL()
	slog.DebugContext(ctx, fmt.Sprintf("send the query %v: %v", query, queryArgs))
	rows, err := a.dbPool.Query(ctx, query, pgx.NamedArgs(queryArgs))
	if err != nil {
		logMsg := fmt.Sprintf("a query error: '%v'", query)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	defer rows.Close()
	var actors []i.Actor
	for rows.Next() {
		var actor i.Actor
		if err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.TitleFi,
			&actor.TitleEn,
			&actor.TitleRu,
			&actor.CreatedAt,
			&actor.UpdatedAt,
			&actor.DeletedAt,
		); err != nil {
			msg := fmt.Sprintf(
				"can not scan an actor from a query result: %v: %v",
				query,
				queryArgs,
			)
			slog.ErrorContext(ctx, msg, slog.Any(logger.ErrorKey, err))
			return nil, err
		}
		actors = append(actors, actor)
	}
	slog.DebugContext(ctx, fmt.Sprintf("received actors: %v", actors))
	return actors, nil
}
