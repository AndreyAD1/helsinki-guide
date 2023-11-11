package integrationtests

import (
	"context"
	"testing"

	i "github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func testBuildingRepository(t *testing.T) {
	storage1 := repositories.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := i.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storage1.Add(context.Background(), neighbourbourhood)

	storage := repositories.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	building := i.Building{
		NameEn: &nameEn,
		Address: i.Address{
			StreetAddress: "test street", 
			NeighbourhoodID: &savedNeighbour.ID,
		},
	}
	saved, err := storage.Add(context.Background(), building)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved.ID)
}