package integrationtests

import (
	"context"
	"path/filepath"
	"testing"

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

	actorRepo := repositories.NewActorRepo(dbpool)
	spec = s.NewActorSpecificationAll(100, 0)
	actors, err := actorRepo.Query(ctx, spec)
	require.NoError(t, err)
	require.Equal(t, 4, len(actors))

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