package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BuildingRepository interface {
	Add(context.Context, i.Building) (*i.Building, error)
	Remove(context.Context, i.Building) error
	Update(context.Context, i.Building) (*i.Building, error)
	Query(context.Context, s.Specification) ([]i.Building, error)
}

type BuildingStorage struct {
	dbPool *pgxpool.Pool
}

func NewBuildingRepo(dbPool *pgxpool.Pool) *BuildingStorage {
	return &BuildingStorage{dbPool}
}

func (b *BuildingStorage) Add(ctx context.Context, building i.Building) (*i.Building, error) {
	created_at := time.Now().Format(time.RFC3339)
	
	batch := &pgx.Batch{}
	var addressID int64
	batch.Queue(
		insertAddress, 
		building.Address.StreetAddress,
		building.Address.NeighbourhoodID,
		created_at,
	).QueryRow(func(row pgx.Row) error {
		err := row.Scan(&addressID)
		if err != nil {
			return processPostgresError(ctx, "address", err)
		}
		return nil
	})

	for _, authorID := range building.AuthorIds {
		batch.Queue(
			insertBuildingAuthor, 
			building.ID,
			authorID,
			created_at,
		).QueryRow(func(row pgx.Row) error {
			err := row.Scan()
			if err != nil {
				return processPostgresError(ctx, "building_author", err)
			}
			return nil
		})
	}

	var buildingID int64
	batch.Queue(
		insertBuilding,
		building.Code,
		building.NameFi,
		building.NameEn,
		building.NameRu,
		addressID,
		building.ConstructionStartYear,
		building.CompletionYear,
		building.ComplexFi,
		building.ComplexEn,
		building.ComplexRu,
		building.HistoryFi,
		building.HistoryEn,
		building.HistoryRu,
		building.ReasoningFi,
		building.ReasoningEn,
		building.ReasoningRu,
		building.ProtectionStatusFi,
		building.ProtectionStatusEn,
		building.ProtectionStatusRu,
		building.InfoSourceFi,
		building.InfoSourceEn,
		building.InfoSourceRu,
		building.SurroundingsFi,
		building.SurroundingsEn,
		building.SurroundingsRu,
		building.FoundationFi,
		building.FoundationEn,
		building.FoundationRu,
		building.FrameFi,
		building.FrameEn,
		building.FrameRu,
		building.FloorDescriptionFi,
		building.FloorDescriptionEn,
		building.FloorDescriptionRu,
		building.FacadesFi,
		building.FacadesEn,
		building.FacadesRu,
		building.SpecialFeaturesFi,
		building.SpecialFeaturesEn,
		building.SpecialFeaturesRu,
		building.Latitude_ETRSGK25,
		building.Longitude_ERRSGK25,
		created_at,
	).QueryRow(func(row pgx.Row) error {
		err := row.Scan(&buildingID)
		if err != nil {
			return processPostgresError(ctx, "building", err)
		}
		building.ID = buildingID

		return nil
	})
	err := b.dbPool.SendBatch(ctx, batch).Close()
	if err != nil {
		return nil, err
	}
	return &building, nil
}

func (b *BuildingStorage) Remove(ctx context.Context, building i.Building) error {
	return ErrNotImplemented
}

func (b *BuildingStorage) Update(ctx context.Context, building i.Building) (*i.Building, error) {
	return nil, ErrNotImplemented
}

func (b *BuildingStorage) Query(
	ctx context.Context,
	spec s.Specification,
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
			&building.SpecialFeaturesFi,
			&building.SpecialFeaturesEn,
			&building.SpecialFeaturesRu,
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

func processPostgresError(ctx context.Context, itemName string, err error) error {
	var pgxError *pgconn.PgError
	if errors.As(err, &pgxError) {
		switch pgxError.Code {
		case pgerrcode.UniqueViolation:
			logMsg := fmt.Sprintf("the %v is not unique", itemName)
			slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			return ErrDuplicate
		case pgerrcode.ForeignKeyViolation:
			logMsg := fmt.Sprintf("the missed %v foreign key", itemName)
			slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			return ErrNoDependency
		}
	}
	logMsg := fmt.Sprintf("unexpected DB error for %v", itemName)
	slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
	return err
}
