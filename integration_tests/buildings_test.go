package integrationtests

import (
	"context"
	"testing"

	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/stretchr/testify/require"
)

func testAddNewBuilding(t *testing.T) {
	storageN := r.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := r.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	actorStorage := r.NewActorRepo(dbpool)
	titleEn := "test title en"
	author1 := r.Actor{Name: "test1", TitleEn: &titleEn}
	author2 := r.Actor{Name: "test2", TitleEn: &titleEn}
	savedAuthor1, err := actorStorage.Add(context.Background(), author1)
	require.NoError(t, err)
	savedAuthor2, err := actorStorage.Add(context.Background(), author2)
	require.NoError(t, err)

	storage := r.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	streetAddress := "test street"
	building := r.Building{
		NameEn: &nameEn,
		Address: r.Address{
			StreetAddress:   streetAddress,
			NeighbourhoodID: &savedNeighbour.ID,
		},
		AuthorIDs: []int64{savedAuthor1.ID, savedAuthor2.ID},
		InitialUses: []r.UseType{
			{NameFi: "use1 fi", NameEn: "use1 en", NameRu: "use1 ru"},
			{NameFi: "use2 fi", NameEn: "use2 en", NameRu: "use2 ru"},
		},
		CurrentUses: []r.UseType{
			{NameFi: "use2 fi", NameEn: "use2 en", NameRu: "use1 ru"},
		},
	}
	saved1, err := storage.Add(context.Background(), building)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved1.ID)
	require.NotEqualValues(t, 0, saved1.Address.ID)

	spec := r.NewBuildingSpecificationByAddress(streetAddress)
	buildings, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equalf(
		t,
		1,
		len(buildings),
		"unexpected building number: %v",
		buildings,
	)
	require.Equal(t, *saved1, buildings[0])

	// save a similar building
	name2 := "rakennus"
	building.NameFi = &name2
	saved2, err := storage.Add(context.Background(), building)
	require.NoError(t, err)
	require.NotEqualValues(t, 0, saved1.ID)
	require.NotEqualValues(t, 0, saved1.Address.ID)
	buildings, err = storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equalf(
		t,
		2,
		len(buildings),
		"unexpected building number: %v",
		buildings,
	)
	require.Equal(t, *saved2, buildings[0])
	require.Equal(t, *saved1, buildings[1])
}

func testAddNewBuildingAddressError(t *testing.T) {
	storage := r.NewBuildingRepo(dbpool)
	nameEn := "test_building"
	streetAddress := "test street"
	absentNeghbourhoodID := int64(999)
	building := r.Building{
		NameEn: &nameEn,
		Address: r.Address{
			StreetAddress:   streetAddress,
			NeighbourhoodID: &absentNeghbourhoodID,
		},
	}
	saved, err := storage.Add(context.Background(), building)
	require.ErrorIs(t, err, r.ErrNoDependency)
	require.Nil(t, saved)

	spec := r.NewBuildingSpecificationByAddress(streetAddress)
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

func testAddNewBuildingAuthorError(t *testing.T) {
	storageN := r.NewNeighbourhoodRepo(dbpool)
	neighbourbourhood := r.Neighbourhood{Name: "test neighbourhood"}
	savedNeighbour, err := storageN.Add(context.Background(), neighbourbourhood)
	require.NoError(t, err)

	storage := r.NewBuildingRepo(dbpool)
	nameEn := "test_building author error"
	streetAddress := "test street author error"
	building := r.Building{
		NameEn: &nameEn,
		Address: r.Address{
			StreetAddress:   streetAddress,
			NeighbourhoodID: &savedNeighbour.ID,
		},
		AuthorIDs: []int64{10, 20},
	}
	saved, err := storage.Add(context.Background(), building)
	require.ErrorIs(t, err, r.ErrNoDependency)
	require.Nil(t, saved)

	spec := r.NewBuildingSpecificationByAddress(streetAddress)
	buildings, err := storage.Query(context.Background(), spec)
	require.NoError(t, err)
	require.Equalf(
		t,
		0,
		len(buildings),
		"unexpected building number: %v",
		len(buildings),
	)
}
