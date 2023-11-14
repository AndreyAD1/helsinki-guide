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
	testSheet = "test" 
	fiFilename = filepath.Join(".", "test_data", "test_fi.xlsx")
	enFilename = filepath.Join(".", "test_data", "test_en.xlsx")
	ruFilename = filepath.Join(".", "test_data", "test_ru.xlsx")

	expectedNeighbourhoods = []internal.Neighbourhood{
		{Name: "Lauttasaari", Municipality: u.StrToPointer("Helsinki")},
		{Name: "Munkkiniemi", Municipality: u.StrToPointer("Helsinki")},
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
	expectedBuildings = []internal.Building{
		{
			ID: int64(1),
			Code: u.StrToPointer("09103100030008001"),
			NameFi: u.StrToPointer("As Oy Meripuistotie 5"),
			NameEn: u.StrToPointer("As Oy Meripuistotie 5"),
			NameRu: u.StrToPointer("As Oy Meripuistotie 5"),
			Address: internal.Address{
				ID: int64(1),
				StreetAddress: "Meripuistotie 5",
				NeighbourhoodID: u.Int64ToPointer(1),
			},
		},
		{
			ID: int64(2),
			Code: u.StrToPointer("09103100030009001"),
			NameFi: u.StrToPointer("Gården Sjöallen 7"),
			NameEn: u.StrToPointer("Gården Sjöallen 7"),
			NameRu: u.StrToPointer("Gården Sjöallen 7"),
			Address: internal.Address{
				ID: int64(2),
				StreetAddress: "Meripuistotie 7",
				NeighbourhoodID: u.Int64ToPointer(2),
			},
		},
		{
			ID: int64(3),
			Code: u.StrToPointer("09103100030001001"),
			NameFi: u.StrToPointer("As Oy Pohjoiskaari 8"),
			NameEn: u.StrToPointer("As Oy Pohjoiskaari 8"),
			NameRu: u.StrToPointer("As Oy Pohjoiskaari 8"),
			Address: internal.Address{
				ID: int64(3),
				StreetAddress: "Pohjoiskaari 8",
				NeighbourhoodID: u.Int64ToPointer(2),
			},
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
		validateBuildingStructs(
			t, 
			reflect.ValueOf(expected), 
			reflect.ValueOf(buildings[i]),
		)
	}
}

func validateBuildingStructs(t *testing.T, expected, actual reflect.Value) {
Out:for i := 0; i < actual.NumField(); i++ {
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
