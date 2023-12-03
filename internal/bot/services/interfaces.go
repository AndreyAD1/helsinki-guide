package services

import "context"

type Buildings interface {
	GetBuildingPreviews(
		ctx context.Context,
		addressPrefix string,
		limit,
		offset int,
	) ([]BuildingPreview, error)
	GetBuildingsByAddress(c context.Context, address string) ([]BuildingDTO, error)
}