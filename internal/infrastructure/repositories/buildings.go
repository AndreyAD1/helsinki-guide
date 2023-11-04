package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AndreyAD1/helsinki-guide/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Building struct {
	ID                    int64
	Code                  *string
	NameFi                *string
	NameEn                *string
	NameRu                *string
	Address               int64
	ConstructionStartYear *int
	CompletionYear        *int
	ComplexFi             *string
	ComplexEn             *string
	ComplexRu             *string
	HistoryFi             *string
	HistoryEn             *string
	HistoryRu             *string
	ReasoningFi           *string
	ReasoningEn           *string
	ReasoningRu           *string
	ProtectionStatusFi    *string
	ProtectionStatusEn    *string
	ProtectionStatusRu    *string
	InfoSourceFi          *string
	InfoSourceEn          *string
	InfoSourceRu          *string
	SurroundingsFi        *string
	SurroundingsEn        *string
	SurroundingsRu        *string
	FoundationFi          *string
	FoundationEn          *string
	FoundationRu          *string
	FrameFi               *string
	FrameEn               *string
	FrameRu               *string
	FloorDescriptionFi    *string
	FloorDescriptionEn    *string
	FloorDescriptionRu    *string
	FacadesFi             *string
	FacadesEn             *string
	FacadesRu             *string
	SpeciaFeaturesFi      *string
	SpeciaFeaturesEn      *string
	SpeciaFeaturesRu      *string
	latitude_ETRSGK25     *float32
	longitude_ERRSGK25    *float32
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	DeletedAt             *time.Time
}

type BuildingWithAddress struct {
	ID            int64
	Code          string
	NameFi        string
	StreetAddress string
}

type BuildingRepository interface {
	GetAllBuildingsAndAddresses(
		ctx context.Context,
		addressPrefix string,
		limit,
		offset int) ([]BuildingWithAddress, error)
	GetBuildingsByAddress(context.Context, string) ([]Building, error)
}

type BuildingStorage struct {
	dbPool *pgxpool.Pool
}

func NewBuildingRepo(dbPool *pgxpool.Pool) *BuildingStorage {
	return &BuildingStorage{dbPool}
}

func (bs *BuildingStorage) GetAllBuildingsAndAddresses(
	ctx context.Context,
	addressPrefix string,
	limit,
	offset int,
) ([]BuildingWithAddress, error) {
	queryTemplate := `SELECT buildings.ID, buildings.code, name_fi,
	street_address FROM buildings JOIN addresses ON 
	buildings.address_id = addresses.id WHERE street_address ILIKE '%v%%'
	ORDER BY street_address LIMIT %v OFFSET %v;`

	query := fmt.Sprintf(
		queryTemplate,
		addressPrefix,
		limit,
		offset,
	)
	rows, err := bs.dbPool.Query(ctx, query)
	if err != nil {
		logTemplate := "can not get addresses: address prefix '%v', limit %v, offset %v, query '%v'"
		logMsg := fmt.Sprintf(
			logTemplate,
			addressPrefix,
			limit,
			offset,
			query,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	var buildings []BuildingWithAddress
	for rows.Next() {
		var building BuildingWithAddress
		if err := rows.Scan(
			&building.ID,
			&building.Code,
			&building.NameFi,
			&building.StreetAddress,
		); err != nil {
			return nil, err
		}
		buildings = append(buildings, building)
	}
	return buildings, nil
}

func (bs *BuildingStorage) GetBuildingsByAddress(
	ctx context.Context,
	address string,
) ([]Building, error) {
	queryTemplate := `SELECT buildings.ID, buildings.code, name_fi, name_en,
	completion_year, history_fi, history_en, history_ru 
	FROM buildings 
	WHERE address_id = (SELECT id FROM addresses WHERE street_address = '%s');`

	query := fmt.Sprintf(queryTemplate, address)
	rows, err := bs.dbPool.Query(ctx, query)
	if err != nil {
		logMsg := fmt.Sprintf("can not get a building: address prefix '%v'", address)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return nil, fmt.Errorf("%v: %w", logMsg, err)
	}
	var buildings []Building
	for rows.Next() {
		var building Building
		if err := rows.Scan(
			&building.ID,
			&building.Code,
			&building.NameFi,
			&building.NameEn,
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
