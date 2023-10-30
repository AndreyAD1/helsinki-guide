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

func (bs BuildingService) GetBuildingPreviews(ctx context.Context) ([]BuildingPreview, error) {
	buildings, err := bs.storage.GetBuildingsWithAddress(ctx, 500, 0)
	if err != nil {
		return nil, err
	}

	previews := make([]BuildingPreview, len(buildings))
	for i, building := range buildings {
		previews[i] = BuildingPreview{building.StreetAddress, building.NameEn}
	}
	return previews, nil
}

func NewBuildingService(storage repositories.BuildingRepository) BuildingService {
	return BuildingService{storage}
}
