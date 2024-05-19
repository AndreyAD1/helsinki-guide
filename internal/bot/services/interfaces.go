package services

import "context"

type Buildings interface {
	GetBuildings(
		ctx context.Context,
		addressPrefix string,
		limit,
		offset int,
	) ([]BuildingDTO, error)
	GetBuildingsByAddress(c context.Context, address string) ([]BuildingDTO, error)
	GetNearestBuildings(
		ctx context.Context,
		distance int,
		latitude,
		longitude float64,
		limit,
		offset int,
	) ([]BuildingDTO, error)
	GetBuildingByID(c context.Context, ID int64) (*BuildingDTO, error)
}
type Users interface {
	GetPreferredLanguage(ctx context.Context, userID int64) (*Language, error)
	SetLanguage(ctx context.Context, userID int64, language Language) error
}
