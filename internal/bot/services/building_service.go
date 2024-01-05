package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
)

type BuildingService struct {
	buildingCollection repositories.BuildingRepository
	actorCollection    repositories.ActorRepository
}

func NewBuildingService(
	buildingCollection repositories.BuildingRepository,
	actorCollection repositories.ActorRepository,
) BuildingService {
	return BuildingService{buildingCollection, actorCollection}
}

func NewBuildingDTO(b r.Building, authors []r.Actor) BuildingDTO {
	var authorNames []string
	for _, author := range authors {
		authorNames = append(authorNames, author.Name)
	}
	var authorPtr *[]string
	if authorNames != nil {
		authorPtr = &authorNames
	}
	return BuildingDTO{
		NameFi:            b.NameFi,
		NameEn:            b.NameEn,
		NameRu:            b.NameRu,
		Address:           b.Address.StreetAddress,
		DescriptionFi:     b.FloorDescriptionFi,
		DescriptionEn:     b.FloorDescriptionEn,
		DescriptionRu:     b.FloorDescriptionRu,
		CompletionYear:    b.CompletionYear,
		Authors:           authorPtr,
		HistoryFi:         b.HistoryFi,
		HistoryEn:         b.HistoryEn,
		HistoryRu:         b.HistoryRu,
		NotableFeaturesFi: b.ReasoningFi,
		NotableFeaturesEn: b.ReasoningEn,
		NotableFeaturesRu: b.ReasoningRu,
		FacadesFi:         b.FacadesFi,
		FacadesEn:         b.FacadesEn,
		FacadesRu:         b.FacadesRu,
		DetailsFi:         b.SpecialFeaturesFi,
		DetailsEn:         b.SpecialFeaturesEn,
		DetailsRu:         b.SpecialFeaturesRu,
		SurroundingsFi:    b.SurroundingsFi,
		SurroundingsEn:    b.SurroundingsEn,
		SurroundingsRu:    b.SpecialFeaturesRu,
	}
}

func (bs BuildingService) GetBuildingPreviews(
	ctx context.Context,
	addressPrefix string,
	limit,
	offset int,
) ([]BuildingPreview, error) {
	addressPrefix = strings.TrimLeft(addressPrefix, " ")
	spec := r.NewBuildingSpecificationByAlikeAddress(addressPrefix, limit, offset)
	buildings, err := bs.buildingCollection.Query(ctx, spec)
	if err != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf("can not get building for '%v'", addressPrefix),
			slog.Any(logger.ErrorKey, err),
		)
		return nil, err
	}

	previews := make([]BuildingPreview, len(buildings))
	for i, building := range buildings {
		name := ""
		if building.NameFi != nil {
			name = *building.NameFi
		}
		previews[i] = BuildingPreview{building.ID, building.Address.StreetAddress, name}
	}
	return previews, nil
}

func (bs BuildingService) GetBuildingsByAddress(
	ctx context.Context,
	address string,
) ([]BuildingDTO, error) {
	address = strings.TrimSpace(address)
	spec := r.NewBuildingSpecificationByAddress(address)
	buildings, err := bs.buildingCollection.Query(ctx, spec)
	if err != nil {
		return nil, err
	}
	buildingsDto := make([]BuildingDTO, len(buildings))
	for i, building := range buildings {
		spec := r.NewAuthorSpecificationByBuilding(building.ID)
		authors, err := bs.actorCollection.Query(ctx, spec)
		if err != nil {
			return nil, err
		}
		buildingsDto[i] = NewBuildingDTO(building, authors)
	}

	return buildingsDto, nil
}

func (bs BuildingService) GetNearestBuildingPreviews(
	ctx context.Context,
	distanceMeters int,
	latitude,
	longitude float64,
	limit,
	offset int,
) ([]BuildingPreview, error) {
	spec := r.NewBuildingSpecificationNearest(
		distanceMeters,
		latitude,
		longitude,
		limit,
		offset,
	)
	buildings, err := bs.buildingCollection.Query(ctx, spec)
	if err != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf(
				"can not get nearest buildings for '%.2f-%.2f'",
				latitude,
				longitude,
			),
			slog.Any(logger.ErrorKey, err),
		)
		return nil, err
	}

	previews := make([]BuildingPreview, len(buildings))
	for i, building := range buildings {
		name := ""
		if building.NameFi != nil {
			name = *building.NameFi
		}
		previews[i] = BuildingPreview{
			building.ID,
			building.Address.StreetAddress,
			name,
		}
	}
	return previews, nil
}
