package services

import (
	"context"

	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
)

type BuildingService struct {
	storage repositories.BuildingRepository
}

type BuildingPreview struct {
	Address string
	Name    string
}

func NewBuildingService(storage repositories.BuildingRepository) BuildingService {
	return BuildingService{storage}
}

func (bs BuildingService) GetBuildingPreviews(ctx context.Context) ([]BuildingPreview, error) {
	buildings, err := bs.storage.GetBuildingsWithAddress(ctx, 500, 0)
	if err != nil {
		return nil, err
	}

	previews := make([]BuildingPreview, len(buildings))
	for i, building := range buildings {
		previews[i] = BuildingPreview{building.StreetAddress, building.NameFi}
	}
	return previews, nil
}

func (bs BuildingService) GetBuildingsByAddress(
	ctx context.Context, 
	address string,
) ([]BuildingDTO, error) {
	buildings, err := bs.storage.GetBuildingsByAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	buildingsDto := make([]BuildingDTO, len(buildings))
	for i, building := range buildings {
		buildingsDto[i] = NewBuildingDTO(building, "en")
	}
	return buildingsDto, nil
}

