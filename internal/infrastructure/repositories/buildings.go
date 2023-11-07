package repositories

import (
	"context"
	"fmt"
	"log/slog"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BuildingRepository interface {
	Add(context.Context, i.Building) (*i.Building, error)
	Remove(context.Context, i.Building) error
	Update(context.Context, i.Building) (*i.Building, error)
	Query(context.Context, s.BuildingSpecification) ([]i.Building, error)
}

type BuildingStorage struct {
	dbPool *pgxpool.Pool
}

func NewBuildingRepo(dbPool *pgxpool.Pool) *BuildingStorage {
	return &BuildingStorage{dbPool}
}

func (b *BuildingStorage) Add(ctx context.Context, building i.Building) (*i.Building, error) {
	// queryTemplate := `INSERT INTO building
	// (code, name_fi, name_en, name_ru, address_id, construction_start_year,
	// completion_year, complex_fi, complex_en, complex_ru, history_fi,
	// history_en, history_ru, reasoning_fi, reasoning_en, reasoning_ru,
	// protection_status_fi, protection_status_en, protection_status_ru,
	// info_source_fi,...
	// )
	// `
	return nil, ErrNotImplemented
}

func (b *BuildingStorage) Remove(ctx context.Context, building i.Building) error {
	return ErrNotImplemented
}

func (b *BuildingStorage) Update(ctx context.Context, building i.Building) (*i.Building, error) {
	return nil, ErrNotImplemented
}

func (b *BuildingStorage) Query(
	ctx context.Context, 
	spec s.BuildingSpecification,
) ([]i.Building, error) {
	query := spec.ToSQL()
	rows, err := b.dbPool.Query(ctx, query)
	if err != nil {
		logMsg := fmt.Sprintf("can not query buildings: '%v'", query)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	var buildings []i.Building
	for rows.Next() {
		var building i.Building
		if err := rows.Scan(
			&building.ID,
			&building.Code,
			&building.NameFi,
			&building.NameEn,
			&building.NameRu,
			&building.Address,
			&building.CompletionYear,
			&building.HistoryFi,
			&building.HistoryEn,
			&building.HistoryRu,
		); err != nil {
			return nil, err
		}
		buildings = append(buildings, building)
	}
	return buildings, nil
}
