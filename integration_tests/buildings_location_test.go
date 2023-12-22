package integrationtests

import (
	"context"
	"log"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	i "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func testGetNearestBuildings(t *testing.T) {
	storageN := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := i.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	storage := repositories.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	address := i.Address{
		StreetAddress:   "test street",
		NeighbourhoodID: &savedNeighbour.ID,
	}
	type buildingInfo struct {
		storedBuilding i.Building
		isExpected     bool
	}
	tests := []struct {
		name              string
		storedBuildings   []i.Building
		distance          int
		latitude          float64
		longitude         float64
		limit             int
		offset            int
		expectedBuildings []i.Building
	}{
		{
			"no buildings",
			[]i.Building{},
			100,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{},
		},
		{
			"one building without coordinates",
			[]i.Building{
				{
					NameEn:  &nameEn,
					Address: address,
				},
			},
			100,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{},
		},
		{
			"one too far building",
			[]i.Building{
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(0.0)),
					Longitude_WGS84: utils.GetPointer(float64(0.0)),
				},
			},
			100,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{},
		},
		{
			"one very close building",
			[]i.Building{
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.36)),
					Longitude_WGS84: utils.GetPointer(float64(24.75)),
				},
			},
			100,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.36)),
					Longitude_WGS84: utils.GetPointer(float64(24.75)),
				},
			},
		},
		{
			"one close and two far buildings",
			[]i.Building{
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(0.0)),
					Longitude_WGS84: utils.GetPointer(float64(0.0)),
				},
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(10.0)),
					Longitude_WGS84: utils.GetPointer(float64(10.0)),
				},
			},
			1000,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
			},
		},
		{
			"two close and two far buildings",
			[]i.Building{
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(0.0)),
					Longitude_WGS84: utils.GetPointer(float64(0.0)),
				},
				{
					NameEn:          utils.GetPointer("second closest"),
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
				{
					NameEn:          &nameEn,
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(10.0)),
					Longitude_WGS84: utils.GetPointer(float64(10.0)),
				},
				{
					NameEn:          utils.GetPointer("closest"),
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.3601)),
					Longitude_WGS84: utils.GetPointer(float64(24.7501)),
				},
			},
			1000,
			60.36,
			24.75,
			10,
			0,
			[]i.Building{
				{
					NameEn:          utils.GetPointer("closest"),
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.3601)),
					Longitude_WGS84: utils.GetPointer(float64(24.7501)),
				},
				{
					NameEn:          utils.GetPointer("second closest"),
					Address:         address,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			for _, b := range tt.storedBuildings {
				saved, err := storage.Add(ctx, b)
				require.NoError(t, err)
				defer func() {
					if err := storage.Remove(ctx, *saved); err != nil {
						log.Printf(
							"can not remove a building '%v' after the test: %v",
							saved.ID,
							err,
						)
					}
				}()
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
			ignoreOption := cmpopts.IgnoreFields(
				i.Building{},
				"ID",
				"Address.ID",
				"Address.CreatedAt",
				"CreatedAt",
			)
			require.Equal(t, cmp.Diff(tt.expectedBuildings, got, ignoreOption), "")
		})
	}
}
