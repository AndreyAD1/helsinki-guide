package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/specifications"
	i "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/stretchr/testify/require"
)

func testUpdateAbsentBuilding(t *testing.T) {
	storageN := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := i.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	storage := repositories.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	streetAddress := "test street"
	building := i.Building{
		NameEn: &nameEn,
		Address: i.Address{
			StreetAddress:   streetAddress,
			NeighbourhoodID: &savedNeighbour.ID,
		},
	}
	updated, err := storage.Update(context.Background(), building)
	require.ErrorIs(t, err, repositories.ErrNotExist)
	require.Nil(t, updated)
}

func testManageRemovedBuilding(t *testing.T) {
	storageN := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := i.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	storage := repositories.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	streetAddress := "test street"
	building := i.Building{
		NameEn: &nameEn,
		Address: i.Address{
			StreetAddress:   streetAddress,
			NeighbourhoodID: &savedNeighbour.ID,
		},
	}
	ctx := context.Background()
	saved, err := storage.Add(ctx, building)
	require.NoError(t, err)

	err = storage.Remove(ctx, *saved)
	require.NoError(t, err)

	_, err = storage.Update(ctx, building)
	require.ErrorIs(t, err, repositories.ErrNotExist)

	spec := s.NewBuildingSpecificationByAddress(streetAddress)
	buildings, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equalf(
		t,
		0,
		len(buildings),
		"unexpected building number: %v",
		buildings,
	)
}
