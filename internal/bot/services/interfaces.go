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
	GetNearestBuildingPreviews(
		ctx context.Context,
		distance int,
		latitude,
		longitude float64,
		limit,
		offset int,
	) ([]BuildingPreview, error)
}
type Users interface {
	SetLanguage(ctx context.Context, userID int64, language string) error
}
