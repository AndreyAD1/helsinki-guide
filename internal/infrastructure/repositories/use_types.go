package repositories

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UseType struct {
	ID int
	NameFi string
	NameEn string
	NameRu string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type UseTypeRepository interface {
	GetUseType (id int) (UseType, error)
	SetUseType (item UseType) (UseType, error)
}

type useTypeStorage struct {
	dbPool *pgxpool.Pool
}

func NewUseTypeRepo(dbPool *pgxpool.Pool) UseTypeRepository {
	return &useTypeStorage{dbPool}
}

func (u *useTypeStorage) GetUseType(id int) (UseType, error) {
	return UseType{}, nil
}

func (u *useTypeStorage) SetUseType(item UseType) (UseType, error) {
	return item, nil
}
