package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	i "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/stretchr/testify/require"
)


func testGetNearestBuildings(t *testing.T) {
	storageN := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := i.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	storage := repositories.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	streetAddress := "test street"
	_ = i.Building{
		NameEn: &nameEn,
		Address: i.Address{
			StreetAddress:   streetAddress,
			NeighbourhoodID: &savedNeighbour.ID,
		},
	}
	tests := []struct {
		name    string
		storedBuildings []i.Building
		distance int
		latitude    float64
		longitude float64
		limit int
		offset int
		expectedBuildings []i.Building
	}{
		{
			"no buildings",
			[]i.Building{},
			10000,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			for _, b := range tt.storedBuildings {
				_, err := storage.Add(ctx, b)
				require.NoError(t, err)
			}
			specification := s.NewBuildingSpecificationNearest(
				tt.distance,
				tt.latitude,
				tt.longitude,
				tt.limit,
				tt.offset,
			)
			got, err := storage.Query(context.Background(), specification)
			require.NoError(t, err)
			require.Equal(t, tt.expectedBuildings, got)
		})
	}
}