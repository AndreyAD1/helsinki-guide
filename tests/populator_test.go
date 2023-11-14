package integrationtests

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/populator"
	"github.com/stretchr/testify/require"
)

var (
	testSheet = "test" 
	fiFilename = filepath.Join(".", "test_data", "test_fi.xlsx")
	enFilename = filepath.Join(".", "test_data", "test_en.xlsx")
	ruFilename = filepath.Join(".", "test_data", "test_ru.xlsx")

	helsinki = "Helsinki"
	expectedNeighbourhoods = []internal.Neighbourhood{
		{Name: "Lauttasaari", Municipality: &helsinki},
		{Name: "Munkkiniemi", Municipality: &helsinki},
	}

	architectFi = "Arkkitehti"
	architectEn = "Architect"
	architectRu = "Архитектор"
	expectedAuthors = []internal.Actor{
		{Name: "Claus Tandefelt", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Kauko Kokko", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu,},
		{Name: "Niilo Kokko", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu,},
		{Name: "Rudolf Lanste"},
	}
)

func testRunPopulator(t *testing.T) {
	ctx := context.Background()
	config := configuration.PopulatorConfig{DatabaseURL: databaseUrl}
	populator, err := populator.NewPopulator(ctx, config)
	require.NoError(t, err)
	err = populator.Run(
		context.Background(), 
		testSheet, 
		fiFilename,
		enFilename,
		ruFilename,
	)
	require.NoError(t, err)

	neigbourhoodRepo := repositories.NewNeighbourhoodRepo(dbpool)
	spec := s.NewNeighbourhoodSpecificationAll(100, 0)
	storedNeighbourhoods, err := neigbourhoodRepo.Query(ctx, spec)
	require.NoError(t, err)
	require.Equal(t, 2, len(storedNeighbourhoods))
	for i, expected := range expectedNeighbourhoods {
		require.Equal(t, expected.Name, storedNeighbourhoods[i].Name)
		if expected.Municipality == nil {
			require.Nil(t, storedNeighbourhoods[i].Municipality)
			continue
		}
		require.Equal(t, expected.Municipality, storedNeighbourhoods[i].Municipality)
	}

	actorRepo := repositories.NewActorRepo(dbpool)
	spec = s.NewActorSpecificationAll(100, 0)
	authors, err := actorRepo.Query(ctx, spec)
	require.NoError(t, err)
	require.Equal(t, 4, len(authors))
	for i, expected := range expectedAuthors {
		require.Equal(t, expected.Name, authors[i].Name)
		if expected.TitleFi == nil {
			require.Nil(t, authors[i].TitleFi)
			require.Nil(t, authors[i].TitleEn)
			require.Nil(t, authors[i].TitleRu)
			continue
		}
		require.Equal(t, *expected.TitleFi, *authors[i].TitleFi)
		require.Equal(t, *expected.TitleEn, *authors[i].TitleEn)
		require.Equal(t, *expected.TitleRu, *authors[i].TitleRu)
	}
	buildingRepo := repositories.NewBuildingRepo(dbpool)
	spec = s.NewBuildingSpecificationByAlikeAddress("", 100, 0)
	buildings, err := buildingRepo.Query(ctx, spec)
	require.NoError(t, err)
	require.Equalf(
		t,
		3,
		len(buildings),
		"unexpected building number: %v",
		buildings,
	)
}