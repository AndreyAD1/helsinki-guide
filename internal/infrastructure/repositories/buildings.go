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

type BuildingStorage struct {
	dbPool *pgxpool.Pool
}

func NewBuildingRepo(dbPool *pgxpool.Pool) *BuildingStorage {
	return &BuildingStorage{dbPool}
}

func (b *BuildingStorage) Add(ctx context.Context, building i.Building) (*i.Building, error) {
	transaction, err := b.dbPool.Begin(ctx)
	if err != nil {
		logMsg := "can not begin a transaction"
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, err
	}
	defer func() {
		logMsg := fmt.Sprintf(
			"finish a transaction for a building '%v'",
			building.Address.StreetAddress,
		)
		slog.DebugContext(ctx, logMsg)
		err := transaction.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		}
	}()

	address, err := b.getAddress(ctx, transaction, building.Address)
	if err != nil {
		return nil, err
	}
	building.Address = address

	err = transaction.QueryRow(
		ctx,
		insertBuilding,
		building.Code,
		building.NameFi,
		building.NameEn,
		building.NameRu,
		address.ID,
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
	).Scan(&building.ID, &building.CreatedAt)
	if err != nil {
		itemName := fmt.Sprintf("building '%v'", building.Address.StreetAddress)
		return nil, processPostgresError(ctx, itemName, err)
	}

	for _, authorID := range building.AuthorIDs {
		res, err := transaction.Exec(
			ctx,
			insertBuildingAuthor,
			building.ID,
			authorID,
		)
		if err != nil {
			itemName := fmt.Sprintf("building author %v", authorID)
			return nil, processPostgresError(ctx, itemName, err)
		}
		if res.RowsAffected() != 1 {
			logMsg := fmt.Sprintf(
				"couldn't add a building author: %v - %v; affected rows: %v",
				building.ID,
				authorID,
				res.RowsAffected(),
			)
			slog.WarnContext(ctx, logMsg)
			return nil, ErrInsertFailed
		}
	}

	uses, err := b.setUses(
		ctx,
		transaction,
		insertInitialUses,
		building.ID,
		building.InitialUses,
	)
	if err != nil {
		return nil, err
	}
	building.InitialUses = uses
	uses, err = b.setUses(
		ctx,
		transaction,
		insertCurrentUses,
		building.ID,
		building.CurrentUses,
	)
	if err != nil {
		return nil, err
	}
	building.CurrentUses = uses

	if err := transaction.Commit(ctx); err != nil {
		logMsg := fmt.Sprintf(
			"can not close a transaction for the building %v - %v",
			building.NameEn,
			building.Address.StreetAddress,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
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
		var address i.Address
		if err := rows.Scan(
			&building.ID,
			&building.Code,
			&building.NameFi,
			&building.NameEn,
			&building.NameRu,
			&address.ID,
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
			&address.NeighbourhoodID,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
		); err != nil {
			return nil, err
		}
		building.Address = address

		authorIDs, err := b.getAuthorIds(ctx, building.ID)
		if err != nil {
			return nil, err
		}
		building.AuthorIDs = authorIDs

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

func (b *BuildingStorage) getAuthorIds(
	ctx context.Context,
	buildingID int64,
) ([]int64, error) {
	authorQuery := `SELECT actor_id FROM building_authors 
	WHERE building_id = $1;`
	rows, err := b.dbPool.Query(ctx, authorQuery, buildingID)
	if err != nil {
		logMsg := fmt.Sprintf(
			"a query error for a building %v: '%v'",
			buildingID,
			authorQuery,
		)
		slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	defer rows.Close()
	var authorIDs []int64
	for rows.Next() {
		var authorID int64
		if err := rows.Scan(&authorID); err != nil {
			msg := fmt.Sprintf(
				"can not scan an author ID for a building '%v'",
				buildingID,
			)
			slog.ErrorContext(ctx, msg, slog.Any(logger.ErrorKey, err))
			return nil, err
		}
		authorIDs = append(authorIDs, authorID)
	}
	return authorIDs, nil
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

func (b *BuildingStorage) getAddress(
	ctx context.Context,
	transaction pgx.Tx,
	address i.Address,
) (i.Address, error) {
	err := transaction.QueryRow(
		ctx,
		getAddress,
		address.StreetAddress,
	).Scan(
		&address.ID,
		&address.StreetAddress,
		&address.NeighbourhoodID,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logMsg := fmt.Sprintf("can not get an address: %v", address)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return i.Address{}, err
	}
	if errors.Is(err, pgx.ErrNoRows) {
		if err := transaction.QueryRow(
			ctx,
			insertAddress,
			address.StreetAddress,
			address.NeighbourhoodID,
		).Scan(
			&address.ID,
			&address.StreetAddress,
			&address.NeighbourhoodID,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
		); err != nil {
			itemName := fmt.Sprintf(
				"address: '%v-%v'",
				address.StreetAddress,
				*address.NeighbourhoodID,
			)
			return i.Address{}, processPostgresError(ctx, itemName, err)
		}
	}
	return address, nil
}

func (b *BuildingStorage) setUses(
	ctx context.Context,
	transaction pgx.Tx,
	insertQuery string,
	buildingID int64,
	uses []i.UseType,
) ([]i.UseType, error) {
	storedUseTypes := []i.UseType{}
	for _, useType := range uses {
		err := transaction.QueryRow(ctx, getUseType, useType.NameEn).Scan(
			&useType.ID,
			&useType.NameFi,
			&useType.NameEn,
			&useType.NameRu,
			&useType.CreatedAt,
			&useType.UpdatedAt,
			&useType.DeletedAt,
		)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			logMsg := fmt.Sprintf("can not get a use type: %v", useType.NameEn)
			slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			return storedUseTypes, err
		}
		if errors.Is(err, pgx.ErrNoRows) {
			if err := transaction.QueryRow(
				ctx,
				insertUseType,
				useType.NameFi,
				useType.NameEn,
				useType.NameRu,
			).Scan(
				&useType.ID,
				&useType.NameFi,
				&useType.NameEn,
				&useType.NameRu,
				&useType.CreatedAt,
				&useType.UpdatedAt,
				&useType.DeletedAt,
			); err != nil {
				itemName := fmt.Sprintf("use type: %v", useType.NameEn)
				return storedUseTypes, processPostgresError(ctx, itemName, err)
			}
		}

		res, err := transaction.Exec(
			ctx,
			insertQuery,
			buildingID,
			useType.ID,
		)
		if err != nil {
			itemName := fmt.Sprintf("building use %v - %v", buildingID, useType.ID)
			return storedUseTypes, processPostgresError(ctx, itemName, err)
		}
		if res.RowsAffected() != 1 {
			logMsg := fmt.Sprintf(
				"couldn't add a building use: %v - %v; affected rows: %v: %v",
				buildingID,
				useType.ID,
				res.RowsAffected(),
				insertQuery,
			)
			slog.WarnContext(ctx, logMsg)
			return storedUseTypes, ErrInsertFailed
		}
		storedUseTypes = append(storedUseTypes, useType)
	}
	return storedUseTypes, nil
}

func processPostgresError(ctx context.Context, itemName string, err error) error {
	var pgxError *pgconn.PgError
	if errors.As(err, &pgxError) {
		switch pgxError.Code {
		case pgerrcode.UniqueViolation:
			logMsg := fmt.Sprintf("the '%v' is not unique", itemName)
			slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			return ErrDuplicate
		case pgerrcode.ForeignKeyViolation:
			logMsg := fmt.Sprintf("the missed '%v' foreign key", itemName)
			slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
			return ErrNoDependency
		}
	}
	logMsg := fmt.Sprintf("the unexpected DB error for '%v'", itemName)
	slog.WarnContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
	return err
}
