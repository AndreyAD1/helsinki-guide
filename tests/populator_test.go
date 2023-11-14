package integrationtests

import (
	"context"
	"log"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/repositories"
	s "github.com/AndreyAD1/helsinki-guide/internal/infrastructure/specifications"
	"github.com/AndreyAD1/helsinki-guide/internal/populator"
	u "github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/require"
)

var (
	testSheet  = "test"
	fiFilename = filepath.Join(".", "test_data", "test_fi.xlsx")
	enFilename = filepath.Join(".", "test_data", "test_en.xlsx")
	ruFilename = filepath.Join(".", "test_data", "test_ru.xlsx")

	expectedNeighbourhoods = []internal.Neighbourhood{
		{Name: "Lauttasaari", Municipality: u.GetStrPointer("Helsinki")},
		{Name: "Munkkiniemi", Municipality: u.GetStrPointer("Helsinki")},
	}

	architectFi     = "Arkkitehti"
	architectEn     = "Architect"
	architectRu     = "Архитектор"
	expectedAuthors = []internal.Actor{
		{Name: "Claus Tandefelt", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Kauko Kokko", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Niilo Kokko", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Rudolf Lanste"},
	}
	expectedBuildings = []internal.Building{
		{
			ID:     int64(1),
			Code:   u.GetStrPointer("09103100030008001"),
			NameFi: u.GetStrPointer("As Oy Meripuistotie 5"),
			NameEn: u.GetStrPointer("As Oy Meripuistotie 5"),
			NameRu: u.GetStrPointer("As Oy Meripuistotie 5"),
			Address: internal.Address{
				ID:              int64(1),
				StreetAddress:   "Meripuistotie 5",
				NeighbourhoodID: u.GetInt64Pointer(1),
			},
			ConstructionStartYear: u.GetIntPointer(1954),
			CompletionYear: u.GetIntPointer(1955),

		},
		{
			ID:     int64(2),
			Code:   u.GetStrPointer("09103100030009001"),
			NameFi: u.GetStrPointer("Gården Sjöallen 7"),
			NameEn: u.GetStrPointer("Gården Sjöallen 7"),
			NameRu: u.GetStrPointer("Gården Sjöallen 7"),
			Address: internal.Address{
				ID:              int64(2),
				StreetAddress:   "Meripuistotie 7",
				NeighbourhoodID: u.GetInt64Pointer(2),
			},
			ConstructionStartYear: nil,
			CompletionYear: u.GetIntPointer(1978),
		},
		{
			ID:     int64(3),
			Code:   u.GetStrPointer("09103100030001001"),
			NameFi: u.GetStrPointer("As Oy Pohjoiskaari 8"),
			NameEn: u.GetStrPointer("As Oy Pohjoiskaari 8"),
			NameRu: u.GetStrPointer("As Oy Pohjoiskaari 8"),
			Address: internal.Address{
				ID:              int64(3),
				StreetAddress:   "Pohjoiskaari 8",
				NeighbourhoodID: u.GetInt64Pointer(1),
			},
			ConstructionStartYear: u.GetIntPointer(1952),
			CompletionYear: u.GetIntPointer(1955),
		},
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
		len(expectedBuildings),
		len(buildings),
		"unexpected building number: %v",
		buildings,
	)
	for i, expected := range expectedBuildings {
		actual := buildings[i]
		require.Equal(t, expected.ID, actual.ID)
		validatePointerField(t, expected.Code, actual.Code)
		validatePointerField(t, expected.NameFi, actual.NameFi)
		validatePointerField(t, expected.NameEn, actual.NameEn)
		validatePointerField(t, expected.NameRu, actual.NameRu)
		require.Equal(t, expected.Address.ID, actual.Address.ID)
		require.Equal(t, expected.Address.StreetAddress, actual.Address.StreetAddress)
		validatePointerField(t, expected.Address.NeighbourhoodID, actual.Address.NeighbourhoodID)
		validatePointerField(t, expected.ConstructionStartYear, actual.ConstructionStartYear)
		validatePointerField(t, expected.CompletionYear, actual.CompletionYear)
		validatePointerField(t, expected.ComplexFi, actual.ComplexFi)
		validatePointerField(t, expected.ComplexEn, actual.ComplexEn)
		validatePointerField(t, expected.ComplexRu, actual.ComplexRu)
		validatePointerField(t, expected.HistoryFi, actual.HistoryFi)
		validatePointerField(t, expected.HistoryEn, actual.HistoryEn)
		validatePointerField(t, expected.HistoryRu, actual.HistoryRu)
		validatePointerField(t, expected.ReasoningFi, actual.ReasoningFi)
		validatePointerField(t, expected.ReasoningEn, actual.ReasoningEn)
		validatePointerField(t, expected.ReasoningRu, actual.ReasoningRu)
		validatePointerField(t, expected.ProtectionStatusFi, actual.ProtectionStatusFi)
		validatePointerField(t, expected.ProtectionStatusEn, actual.ProtectionStatusEn)
		validatePointerField(t, expected.ProtectionStatusRu, actual.ProtectionStatusRu)
		validatePointerField(t, expected.InfoSourceFi, actual.InfoSourceFi)
		validatePointerField(t, expected.InfoSourceEn, actual.InfoSourceEn)
		validatePointerField(t, expected.InfoSourceRu, actual.InfoSourceRu)
		validatePointerField(t, expected.SurroundingsFi, actual.SurroundingsFi)
		validatePointerField(t, expected.SurroundingsEn, actual.SurroundingsEn)
		validatePointerField(t, expected.SurroundingsRu, actual.SurroundingsRu)
		validatePointerField(t, expected.FoundationFi, actual.FoundationFi)
		validatePointerField(t, expected.FoundationFi, actual.FoundationFi)
		validatePointerField(t, expected.FoundationRu, actual.FoundationRu)
		validatePointerField(t, expected.FrameFi, actual.FrameFi)
		validatePointerField(t, expected.FrameEn, actual.FrameEn)
		validatePointerField(t, expected.FrameRu, actual.FrameRu)
		validatePointerField(t, expected.FloorDescriptionFi, actual.FloorDescriptionFi)
		validatePointerField(t, expected.FloorDescriptionEn, actual.FloorDescriptionEn)
		validatePointerField(t, expected.FloorDescriptionRu, actual.FloorDescriptionRu)
		validatePointerField(t, expected.FacadesFi, actual.FacadesFi)
		validatePointerField(t, expected.FacadesEn, actual.FacadesEn)
		validatePointerField(t, expected.FacadesRu, actual.FacadesRu)
		validatePointerField(t, expected.SpecialFeaturesFi, actual.SpecialFeaturesFi)
		validatePointerField(t, expected.SpecialFeaturesEn, actual.SpecialFeaturesEn)
		validatePointerField(t, expected.SpecialFeaturesRu, actual.SpecialFeaturesRu)
		validatePointerField(t, expected.Latitude_ETRSGK25, actual.Latitude_ETRSGK25)
		validatePointerField(t, expected.Longitude_ERRSGK25, actual.Longitude_ERRSGK25)
		require.Equal(t, expected.AuthorIDs, actual.AuthorIDs)
		require.Equal(t, len(expected.InitialUses), len(actual.InitialUses))
		for i, expectedUse := range expected.InitialUses {
			actualUse := actual.InitialUses[i]
			require.Equal(t, expectedUse.ID, actualUse.ID)
			require.Equal(t, expectedUse.NameFi, actualUse.NameFi)
			require.Equal(t, expectedUse.NameEn, actualUse.NameEn)
			require.Equal(t, expectedUse.NameRu, actualUse.NameRu)
		}
		require.Equal(t, len(expected.CurrentUses), len(actual.CurrentUses))
		for i, expectedUse := range expected.CurrentUses {
			actualUse := actual.CurrentUses[i]
			require.Equal(t, expectedUse.ID, actualUse.ID)
			require.Equal(t, expectedUse.NameFi, actualUse.NameFi)
			require.Equal(t, expectedUse.NameEn, actualUse.NameEn)
			require.Equal(t, expectedUse.NameRu, actualUse.NameRu)
		}
	}
	for i, expected := range expectedBuildings {
		validateBuildingStructs(
			t,
			reflect.ValueOf(expected),
			reflect.ValueOf(buildings[i]),
		)
	}
}

func validatePointerField[P string|int|int64|float32](
	t *testing.T,
	expected, 
	actual *P,
) {
	if expected == nil {
		require.Nil(t, actual)
		return
	}
	require.Equal(t, *expected, *actual)
}

func validateBuildingStructs(t *testing.T, expected, actual reflect.Value) {
Out:
	for i := 0; i < actual.NumField(); i++ {
		log.Printf("check a field %v", expected.Field(i).Type().Name())
		if expected.Field(i).Type().Name() == "Timestamps" {
			continue
		}
		expectedValue := expected.Field(i)
		actualValue := actual.Field(i)
		switch actualValue.Kind() {
		case reflect.Struct:
			validateBuildingStructs(t, expectedValue, actualValue)
		case reflect.Array:
			if expectedValue.Len() == 0 {
				continue Out
			}
			switch expectedValue.Index(0).Kind() {
			case reflect.Struct:
				for i := 0; i < expectedValue.Len(); i++ {
					validateBuildingStructs(
						t,
						expectedValue.Index(i),
						actualValue.Index(i),
					)
				}
			case reflect.Int, reflect.Int64:
				for i := 0; i < expectedValue.Len(); i++ {
					require.Equal(
						t,
						expectedValue.Index(i).Int(),
						actualValue.Index(i).Int(),
					)
				}
			}
		case reflect.Int, reflect.Int64:
			require.Equal(t, expectedValue.Int(), actualValue.Int())
		case reflect.String:
			expectedStr := expectedValue.String()
			require.Equal(t, expectedStr, actualValue.String())
		case reflect.Pointer:
			if expectedValue.IsNil() {
				require.True(t, actualValue.IsNil())
				continue Out
			}
			pointerValue := actualValue.Elem()
			switch pointerValue.Kind() {
			case reflect.Struct:
				for i := 0; i < expectedValue.Len(); i++ {
					validateBuildingStructs(
						t,
						expectedValue.Index(i),
						actualValue.Index(i),
					)
				}
			case reflect.String:
				expectedStr := expectedValue.Elem().String()
				require.Equal(t, expectedStr, actualValue.Elem().String())
			case reflect.Int:
				expectedInt := expectedValue.Elem().Int()
				require.Equal(t, expectedInt, actualValue.Elem().Int())
			case reflect.Float32:
				expectedStr := expectedValue.Elem().Float()
				require.Equal(t, expectedStr, actualValue.Elem().Float())
			}
		}
	}
}
