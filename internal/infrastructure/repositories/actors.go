package repositories

import (
	"context"
	"fmt"
	"log/slog"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActorRepository interface {
	Add(context.Context, i.Actor) (*i.Actor, error)
	Remove(context.Context, i.Actor) error
	Update(context.Context, i.Actor) (*i.Actor, error)
	Query(context.Context, s.Specification) ([]i.Actor, error)
}

type actorStorage struct {
	dbPool *pgxpool.Pool
}

func NewActorRepo(dbPool *pgxpool.Pool) ActorRepository {
	return &actorStorage{dbPool}
}

func (a *actorStorage) Add(ctx context.Context, actor i.Actor) (*i.Actor, error) {
	return nil, ErrNotImplemented
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
