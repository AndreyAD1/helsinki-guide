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
	// 	completion_year, complex_fi, complex_en, complex_ru, history_fi,
	// 	history_en, history_ru, reasoning_fi, reasoning_en, reasoning_ru,
	// 	protection_status_fi, protection_status_en, protection_status_ru,
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
	query, queryArgs := spec.ToSQL()
	slog.DebugContext(ctx, fmt.Sprintf("send the query %v: %v", query, queryArgs))
	rows, err := b.dbPool.Query(ctx, query, pgx.NamedArgs(queryArgs))
	if err != nil {
		logMsg := fmt.Sprintf("a query error: '%v'", query)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	defer rows.Close()
	var buildings []i.Building
	for rows.Next() {
		var building i.Building
		var addressID int64
		var neighbourhoodID int64
		var address i.Address
		if err := rows.Scan(
			&building.ID,
			&building.Code,
			&building.NameFi,
			&building.NameEn,
			&building.NameRu,
			&addressID,
			&building.ConstructionStartYear,
			&building.CompletionYear,
			&building.ComplexFi,
			&building.ComplexEn,
			&building.ComplexRu,
			&building.HistoryFi,
			&building.HistoryEn,
			&building.HistoryRu,
			&building.ReasoningFi,
			&building.ReasoningEn,
			&building.ReasoningRu,
			&building.ProtectionStatusFi,
			&building.ProtectionStatusEn,
			&building.ProtectionStatusRu,
			&building.InfoSourceFi,
			&building.InfoSourceEn,
			&building.InfoSourceRu,
			&building.SurroundingsFi,
			&building.SurroundingsEn,
			&building.SurroundingsRu,
			&building.FoundationFi,
			&building.FoundationEn,
			&building.FoundationRu,
			&building.FrameFi,
			&building.FrameEn,
			&building.FrameRu,
			&building.FloorDescriptionFi,
			&building.FloorDescriptionEn,
			&building.FloorDescriptionRu,
			&building.FacadesFi,
			&building.FacadesEn,
			&building.FacadesRu,
			&building.SpeciaFeaturesFi,
			&building.SpeciaFeaturesEn,
			&building.SpeciaFeaturesRu,
			&building.Latitude_ETRSGK25,
			&building.Longitude_ERRSGK25,
			&building.CreatedAt,
			&building.UpdatedAt,
			&building.DeletedAt,
			&address.ID,
			&address.StreetAddress,
			&neighbourhoodID,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
		); err != nil {
			return nil, err
		}
		building.Address = address
		uses, err := b.getUses(ctx, initialUsesTable, building.ID)
		if err != nil {
			logMsg := fmt.Sprintf(
				"can not get uses for a building '%v'", 
				building.ID,
			)
			slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			buildings = append(buildings, building)
			continue
		}
		building.InitialUses = uses
		
		uses, err = b.getUses(ctx, currentUsesTable, building.ID)
		if err != nil {
			logMsg := fmt.Sprintf(
				"can not get uses for a building '%v'", 
				building.ID,
			)
			slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			buildings = append(buildings, building)
			continue
		}
		building.CurrentUses = uses
		buildings = append(buildings, building)
	}
	return buildings, nil
}

type UseTableNames string

const (
	initialUsesTable = UseTableNames("initial_uses")
	currentUsesTable = UseTableNames("current_uses")
)

func (b *BuildingStorage) getUses(
	ctx context.Context,
	table_name UseTableNames,
	buildingID int64,
) ([]i.UseType, error) {
	query := fmt.Sprintf(`SELECT id, name_fi, name_en, name_ru, 
	created_at,updated_at, deleted_at FROM use_types JOIN %v
	ON id = use_type_id WHERE building_id = $1;`, table_name)
	rows, err := b.dbPool.Query(ctx, query, buildingID)
	if err != nil {
		logMsg := fmt.Sprintf("a query error: '%v'", query)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	defer rows.Close()
	var uses []i.UseType
	for rows.Next() {
		var use i.UseType
		if err := rows.Scan(
			&use.ID, 
			&use.NameFi, 
			&use.NameEn, 
			&use.NameRu, 
			&use.CreatedAt,
			&use.UpdatedAt,
			&use.DeletedAt,
		); err != nil {
			msg := fmt.Sprintf(
				"can not scan %v for a building '%v'", 
				table_name, 
				buildingID,
			)
			slog.ErrorContext(ctx, msg, slog.Any(logger.ErrorKey, err))
			return nil, err
		}
		uses = append(uses, use)
	}
	return uses, nil
}