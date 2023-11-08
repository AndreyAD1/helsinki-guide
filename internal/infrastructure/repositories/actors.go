package repositories

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Actor struct {
	ID        int
	Name      string
	TitleFi   *string
	TitleEn   *string
	TitleRu   *string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type ActorRepository interface {
	GetActor(id int) (*Actor, error)
	SetActor(item Actor) (*Actor, error)
}

type actorStorage struct {
	dbPool *pgxpool.Pool
}

func NewActorRepo(dbPool *pgxpool.Pool) ActorRepository {
	return &actorStorage{dbPool}
}

func (a *actorStorage) GetActor(id int) (*Actor, error) {
	return &Actor{}, nil
}

func (a *actorStorage) SetActor(item Actor) (*Actor, error) {
	return &item, nil
}
