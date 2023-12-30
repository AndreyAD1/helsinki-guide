package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userStorage struct {
	dbPool *pgxpool.Pool
}

func NewUserRepo(dbPool *pgxpool.Pool) *userStorage {
	return &userStorage{dbPool}
}

func (s *userStorage) Add(ctx context.Context, user User) (*User, error) {
	return nil, ErrNotImplemented
}

func (s *userStorage) AddOrUpdate(ctx context.Context, user User) (*User, error) {
	insertQuery := `INSERT INTO users (telegram_id, language)
	VALUES ($1, $2) ON CONFLICT (telegram_id) DO UPDATE 
	SET language = $2, updated_at = now()
	RETURNING id, created_at, updated_at;`
	err := s.dbPool.QueryRow(
		ctx,
		insertQuery,
		user.TelegramID,
		user.PreferredLanguage,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logMsg := fmt.Sprintf(
			"can not add or update a user %v: %v",
			user.TelegramID,
			user.PreferredLanguage,
		)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, err
	}
	return &user, nil
}

func (s *userStorage) Remove(ctx context.Context, user User) error {
	return ErrNotImplemented
}

func (s *userStorage) Update(ctx context.Context, user User) (*User, error) {
	return nil, ErrNotImplemented
}

func (s *userStorage) Query(ctx context.Context, spec Specification) ([]User, error) {
	return nil, ErrNotImplemented
}
