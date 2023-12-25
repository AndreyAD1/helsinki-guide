package integrationtests

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	r "github.com/AndreyAD1/helsinki-guide/internal/bot/infrastructure/repositories"
	"github.com/AndreyAD1/helsinki-guide/internal/populator"
	u "github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/require"
)

var (
	testSheet  = "test"
	fiFilename = filepath.Join(".", "testdata", "test_fi.xlsx")
	enFilename = filepath.Join(".", "testdata", "test_en.xlsx")
	ruFilename = filepath.Join(".", "testdata", "test_ru.xlsx")

	expectedNeighbourhoods = []r.Neighbourhood{
		{Name: "Lauttasaari", Municipality: u.GetPointer("Helsinki")},
		{Name: "Munkkiniemi", Municipality: u.GetPointer("Helsinki")},
	}

	architectFi     = "Arkkitehti"
	architectEn     = "Architect"
	architectRu     = "Архитектор"
	expectedAuthors = []r.Actor{
		{Name: "Claus Tandefelt", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Kauko Kokko", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Niilo Kokko", TitleFi: &architectFi, TitleEn: &architectEn, TitleRu: &architectRu},
		{Name: "Rudolf Lanste"},
	}
	expectedBuildings = []r.Building{
		{
			ID:     int64(1),
			Code:   u.GetPointer("09103100030008001"),
			NameFi: u.GetPointer("As Oy Meripuistotie 5"),
			NameEn: u.GetPointer("As Oy Meripuistotie 5"),
			NameRu: u.GetPointer("As Oy Meripuistotie 5"),
			Address: r.Address{
				ID:              int64(1),
				StreetAddress:   "Meripuistotie 5",
				NeighbourhoodID: u.GetPointer(int64(1)),
			},
			ConstructionStartYear: u.GetPointer(1954),
			CompletionYear:        u.GetPointer(1955),
			AuthorIDs:             []int64{1},
			InitialUses: []r.UseType{
				{
					ID:     1,
					NameFi: "kerrostalot",
					NameEn: "apartment buildings",
					NameRu: "многоквартирные дома",
				},
			},
			CurrentUses: []r.UseType{
				{
					ID:     1,
					NameFi: "kerrostalot",
					NameEn: "apartment buildings",
					NameRu: "многоквартирные дома",
				},
			},
			HistoryFi:          u.GetPointer("Vuonna 1988 rakennuksen kellariin suunniteltiin taloyhtiölle saunaosasto (Arkkitehti Hannu Lehto). Vuonna 1994 rakennuksen parvekkeet suunniteltiin varustettaviksi sivuun siirrettävillä lasiseinillä (arkkitehti Tapani Virkkala)."),
			HistoryEn:          u.GetPointer("In 1988, a sauna section was planned for the housing association in the basement of the building (Architect Hannu Lehto). In 1994, the building's balconies were designed to be equipped with glass walls that can be moved to the side (architect Tapani Virkkala)."),
			HistoryRu:          u.GetPointer("В 1988 году для жилищного товарищества в подвале здания было спроектировано помещение сауны (архитектор Ханну Лехто). В 1994 году балконы здания было спроектировано оборудованными стеклянными стенами, которые можно сдвигать в сторону (архитектор Тапани Вирккала)."),
			ReasoningFi:        u.GetPointer("Reasoning 1 fi"),
			ReasoningEn:        u.GetPointer("Reasoning 1 en"),
			ReasoningRu:        u.GetPointer("Reasoning 1 ru"),
			FoundationFi:       u.GetPointer("betoni"),
			FoundationEn:       u.GetPointer("concrete"),
			FoundationRu:       u.GetPointer("бетон"),
			FrameFi:            u.GetPointer("betoni"),
			FrameEn:            u.GetPointer("concrete"),
			FrameRu:            u.GetPointer("бетон"),
			FloorDescriptionFi: u.GetPointer("description fi 1"),
			FloorDescriptionEn: u.GetPointer("description en 1"),
			FloorDescriptionRu: u.GetPointer("description ru 1"),
			FacadesFi:          u.GetPointer("facade fi 1"),
			FacadesEn:          u.GetPointer("facade en 1"),
			FacadesRu:          u.GetPointer("facade ru 1"),
			Latitude_ETRSGK25:  u.GetPointer(float32(6671929)),
			Longitude_ETRSGK25: u.GetPointer(float32(25493834)),
		},
		{
			ID:     int64(2),
			Code:   u.GetPointer("09103100030009001"),
			NameFi: u.GetPointer("Gården Sjöallen 7"),
			NameEn: u.GetPointer("Gården Sjöallen 7"),
			NameRu: u.GetPointer("Gården Sjöallen 7"),
			Address: r.Address{
				ID:              int64(2),
				StreetAddress:   "Meripuistotie 7",
				NeighbourhoodID: u.GetPointer(int64(2)),
			},
			ConstructionStartYear: nil,
			CompletionYear:        u.GetPointer(1978),
			AuthorIDs:             []int64{2},
			InitialUses: []r.UseType{
				{
					ID:     1,
					NameFi: "kerrostalot",
					NameEn: "apartment buildings",
					NameRu: "многоквартирные дома",
				},
				{
					ID:     2,
					NameFi: "päiväkodit",
					NameEn: "kindergartens",
					NameRu: "детские сады",
				},
			},
			CurrentUses: []r.UseType{
				{
					ID:     2,
					NameFi: "päiväkodit",
					NameEn: "kindergartens",
					NameRu: "детские сады",
				},
			},
			FoundationFi:       u.GetPointer("betoni"),
			FoundationEn:       u.GetPointer("concrete"),
			FoundationRu:       u.GetPointer("бетон"),
			FrameFi:            u.GetPointer("betoni"),
			FrameEn:            u.GetPointer("concrete"),
			FrameRu:            u.GetPointer("бетон"),
			FloorDescriptionFi: u.GetPointer("description fi 2"),
			FloorDescriptionEn: u.GetPointer("description en 2"),
			FloorDescriptionRu: u.GetPointer("description ru 2"),
			FacadesFi:          u.GetPointer("facade fi 2"),
			FacadesEn:          u.GetPointer("facade en 2"),
			FacadesRu:          u.GetPointer("facade ru 2"),
			Latitude_ETRSGK25:  u.GetPointer(float32(6671951)),
			Longitude_ETRSGK25: u.GetPointer(float32(25493874)),
		},
		{
			ID:     int64(3),
			Code:   u.GetPointer("09103100030001001"),
			NameFi: u.GetPointer("As Oy Pohjoiskaari 8"),
			NameEn: u.GetPointer("As Oy Pohjoiskaari 8"),
			NameRu: u.GetPointer("As Oy Pohjoiskaari 8"),
			Address: r.Address{
				ID:              int64(3),
				StreetAddress:   "Pohjoiskaari 8",
				NeighbourhoodID: u.GetPointer(int64(1)),
			},
			ConstructionStartYear: u.GetPointer(1952),
			CompletionYear:        u.GetPointer(1955),
			AuthorIDs:             []int64{3, 4},
			InitialUses: []r.UseType{
				{
					ID:     1,
					NameFi: "kerrostalot",
					NameEn: "apartment buildings",
					NameRu: "многоквартирные дома",
				},
			},
			CurrentUses: []r.UseType{
				{
					ID:     1,
					NameFi: "kerrostalot",
					NameEn: "apartment buildings",
					NameRu: "многоквартирные дома",
				},
			},
			HistoryFi:          u.GetPointer("Vuonna 1996 rakennuksen ikkunat ja parvekeovet uusittiin puu-alumiinirakenteisiksi (Fenestra Oy)."),
			HistoryEn:          u.GetPointer("In 1996, the building's windows and balcony doors were renewed with wood-aluminum structures (Fenestra Oy)."),
			HistoryRu:          u.GetPointer("В 1996 году окна и балконные двери здания были заменены дерево-алюминиевыми конструкциями (Fenestra Oy)."),
			ReasoningFi:        u.GetPointer("Reasoning 3 fi"),
			ReasoningEn:        u.GetPointer("Reasoning 3 en"),
			ReasoningRu:        u.GetPointer("Reasoning 3 ru"),
			FoundationFi:       u.GetPointer("betoni"),
			FoundationEn:       u.GetPointer("concrete"),
			FoundationRu:       u.GetPointer("бетон"),
			FrameFi:            u.GetPointer("tiili"),
			FrameEn:            u.GetPointer("brick"),
			FrameRu:            u.GetPointer("кирпич"),
			FloorDescriptionFi: u.GetPointer("description fi 3"),
			FloorDescriptionEn: u.GetPointer("description en 3"),
			FloorDescriptionRu: u.GetPointer("description ru 3"),
			FacadesFi:          u.GetPointer("facade fi 3"),
			FacadesEn:          u.GetPointer("facade en 3"),
			FacadesRu:          u.GetPointer("facade ru 3"),
			Latitude_ETRSGK25:  u.GetPointer(float32(6671911)),
			Longitude_ETRSGK25: u.GetPointer(float32(25494008)),
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

	neigbourhoodRepo := r.NewNeighbourhoodRepo(dbpool)
	spec := r.NewNeighbourhoodSpecificationAll(100, 0)
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

	actorRepo := r.NewActorRepo(dbpool)
	spec = r.NewActorSpecificationAll(100, 0)
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
	buildingRepo := r.NewBuildingRepo(dbpool)
	spec = r.NewBuildingSpecificationByAlikeAddress("", 100, 0)
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
		validatePointerField(t, expected.Longitude_ETRSGK25, actual.Longitude_ETRSGK25)
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
}

func validatePointerField[P string | int | int64 | float32](
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
