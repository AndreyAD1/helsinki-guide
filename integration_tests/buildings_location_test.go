package integrationtests

import (
	"context"
	"log"
	"testing"

	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func testGetNearestBuildings(t *testing.T) {
	storageN := r.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := r.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	storage := r.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	address1 := r.Address{
		StreetAddress:   "test street1",
		NeighbourhoodID: &savedNeighbour.ID,
	}
	address2 := r.Address{
		StreetAddress:   "test street2",
		NeighbourhoodID: &savedNeighbour.ID,
	}
	type buildingInfo struct {
		storedBuilding r.Building
		isExpected     bool
	}
	tests := []struct {
		name              string
		storedBuildings   []r.Building
		distance          int
		latitude          float64
		longitude         float64
		limit             int
		offset            int
		expectedBuildings []r.Building
	}{
		{
			"no buildings",
			[]r.Building{},
			100,
			60.36,
			24.75,
			10,
			0,
			[]r.Building{},
		},
		{
			"one building without coordinates",
			[]r.Building{
				{
					NameEn:  &nameEn,
					Address: address1,
				},
			},
			100,
			60.36,
			24.75,
			10,
			0,
			[]r.Building{},
		},
		{
			"one too far building",
			[]r.Building{
				{
					NameEn:          &nameEn,
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(0.0)),
					Longitude_WGS84: utils.GetPointer(float64(0.0)),
				},
			},
			100,
			60.36,
			24.75,
			10,
			0,
			[]r.Building{},
		},
		{
			"one very close building",
			[]r.Building{
				{
					NameEn:          &nameEn,
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(60.36)),
					Longitude_WGS84: utils.GetPointer(float64(24.75)),
				},
			},
			100,
			60.36,
			24.75,
			10,
			0,
			[]r.Building{
				{
					NameEn:          &nameEn,
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(60.36)),
					Longitude_WGS84: utils.GetPointer(float64(24.75)),
				},
			},
		},
		{
			"one close and two far buildings",
			[]r.Building{
				{
					NameEn:          &nameEn,
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(0.0)),
					Longitude_WGS84: utils.GetPointer(float64(0.0)),
				},
				{
					NameEn:          &nameEn,
					Address:         address2,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
				{
					NameEn:          &nameEn,
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(10.0)),
					Longitude_WGS84: utils.GetPointer(float64(10.0)),
				},
			},
			1000,
			60.36,
			24.75,
			10,
			0,
			[]r.Building{
				{
					NameEn:          &nameEn,
					Address:         address2,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
			},
		},
		{
			"two close and two far buildings",
			[]r.Building{
				{
					NameEn:          &nameEn,
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(0.0)),
					Longitude_WGS84: utils.GetPointer(float64(0.0)),
				},
				{
					NameEn:          utils.GetPointer("second closest"),
					Address:         address1,
					Latitude_WGS84:  utils.GetPointer(float64(60.361)),
					Longitude_WGS84: utils.GetPointer(float64(24.751)),
				},
				{
					NameEn:          &nameEn,
					Address:         address2,
					Latitude_WGS84:  utils.GetPointer(float64(10.0)),
					Longitude_WGS84: utils.GetPointer(float64(10.0)),
				},
				{
					NameEn:          utils.GetPointer("closest"),
					Address:         address2,
					Latitude_WGS84:  utils.GetPointer(float64(60.3601)),
					Longitude_WGS84: utils.GetPointer(float64(24.7501)),
				},
			},
			1000,
			60.36,
			24.75,
			10,
			0,
			[]r.Building{
				{
					NameEn:          utils.GetPointer("closest"),
					Address:         address2,
					Latitude_WGS84:  utils.GetPointer(float64(60.3601)),
					Longitude_WGS84: utils.GetPointer(float64(24.7501)),
				},
				{
					NameEn:          utils.GetPointer("second closest"),
					Address:         address1,
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
			specification := r.NewBuildingSpecificationNearest(
				tt.distance,
				tt.latitude,
				tt.longitude,
				tt.limit,
				tt.offset,
			)
			got, err := storage.Query(context.Background(), specification)
			require.NoError(t, err)
			ignoreOption := cmpopts.IgnoreFields(
				r.Building{},
				"ID",
				"Address.ID",
				"Address.CreatedAt",
				"CreatedAt",
			)
			require.Equal(
				t,
				cmp.Diff(
					tt.expectedBuildings,
					got,
					cmpopts.IgnoreUnexported(r.Timestamps{}),
					ignoreOption,
				),
				"",
			)
		})
	}
}
