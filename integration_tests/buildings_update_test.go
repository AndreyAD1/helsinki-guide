package integrationtests

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	i "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories/types"
	"github.com/stretchr/testify/require"
)

func testUpdateBuildings(t *testing.T) {
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