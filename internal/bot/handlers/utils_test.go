package handlers

import (
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/require"
)

var extendedBuilding = services.BuildingDTO{
	NameFi:         utils.GetPointer("name fi"),
	NameEn:         utils.GetPointer("name en"),
	NameRu:         utils.GetPointer("name ru"),
	Address:        "test address",
	DescriptionFi: utils.GetPointer("description fi"),
	DescriptionEn: utils.GetPointer("description en"),
	DescriptionRu: utils.GetPointer("description ru"),
	CompletionYear: utils.GetPointer(2023),
	Authors:        &[]string{"Author 1", "Author2"},
	HistoryFi:      utils.GetPointer("history fi"),
	HistoryEn:      utils.GetPointer("history en"),
	HistoryRu:      utils.GetPointer("history ru"),
	NotableFeaturesFi: utils.GetPointer("features fi"),
	NotableFeaturesEn: utils.GetPointer("features en"),
	NotableFeaturesRu: utils.GetPointer("features ru"),
	FacadesFi: utils.GetPointer("facades fi"),
	FacadesEn: utils.GetPointer("facades en"),
	FacadesRu: utils.GetPointer("facades ru"),
	DetailsFi: utils.GetPointer("details fi"),
	DetailsEn: utils.GetPointer("details en"),
	DetailsRu: utils.GetPointer("details ru"),
	SurroundingsFi: utils.GetPointer("surroundings fi"),
	SurroundingsEn: utils.GetPointer("surroundings en"),
	SurroundingsRu: utils.GetPointer("surroundings ru"),
}

func TestSerializeIntoMessage_positive(t *testing.T) {
	dummyBuilding := services.BuildingDTO{Address: "osoite"}
	type args struct {
		object         any
		outputLanguage outputLanguage
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			"dummy fi",
			args{dummyBuilding, Finnish},
			`Nimi: no data
Katuosoite: osoite
Kerrosluku: no data
Käyttöönottovuosi: no data
Suunnittelijat: no data
Rakennushistoria: no data
Huomattavia ominaisuuksia: no data
Julkisivut: no data
Erityispiirteet: no data
Ymparistonkuvaus: no data`,
		},
		{
			"dummy en",
			args{dummyBuilding, English},
			`Name: no data
Address: osoite
Description: no data
Completion year: no data
Authors: no data
Building history: no data
Notable features: no data
Facades: no data
Interesting details: no data
Surroundings: no data`,
		},
		{
			"dummy ru",
			args{dummyBuilding, Russian},
			`Имя: нет данных
Адрес: osoite
Описание: нет данных
Год постройки: нет данных
Авторы: нет данных
История здания: нет данных
Примечательные особенности: нет данных
Фасады: нет данных
Интересные детали: нет данных
Окрестности: нет данных`,
		},
		{
			"extended fi",
			args{extendedBuilding, Finnish},
			`Nimi: name fi
Katuosoite: test address
Kerrosluku: description fi
Käyttöönottovuosi: 2023
Suunnittelijat: Author 1, Author2
Rakennushistoria: history fi
Huomattavia ominaisuuksia: features fi
Julkisivut: facades fi
Erityispiirteet: details fi
Ymparistonkuvaus: surroundings fi`,
		},
		{
			"extended en",
			args{extendedBuilding, English},
			`Name: name en
Address: test address
Description: description en
Completion year: 2023
Authors: Author 1, Author2
Building history: history en
Notable features: features en
Facades: facades en
Interesting details: details en
Surroundings: surroundings en`,
		},
		{
			"extended ru",
			args{extendedBuilding, Russian},
			`Имя: name ru
Адрес: test address
Описание: description ru
Год постройки: 2023
Авторы: Author 1, Author2
История здания: history ru
Примечательные особенности: features ru
Фасады: facades ru
Интересные детали: details ru
Окрестности: surroundings ru`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeIntoMessage(tt.args.object, tt.args.outputLanguage)
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestSerializeIntoMessage_negative(t *testing.T) {
	type testTagStruct struct {
		field1 string `valueLanguage:"en"`
	}
	type testTypeStruct struct {
		field1 rune `valueLanguage:"all" nameEn:"field1"`
	}
	type testPointerTypeStruct struct {
		field1 *rune `valueLanguage:"all" nameEn:"field1"`
	}
	type args struct {
		object         any
		outputLanguage outputLanguage
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			"not a structure",
			args{123, English},
			ErrUnexpectedType,
		},
		{
			"a structure with no field tags",
			args{args{}, English},
			ErrNoFieldTag,
		},
		{
			"a structure with no name tag",
			args{testTagStruct{"test"}, English},
			ErrNoNameTag,
		},
		{
			"a structure with an unexpected field type",
			args{testTypeStruct{'A'}, English},
			ErrUnexpectedFieldType,
		},
		{
			"a structure with an unexpected pointer field type",
			args{testPointerTypeStruct{utils.GetPointer('T')}, English},
			ErrUnexpectedFieldType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SerializeIntoMessage(tt.args.object, tt.args.outputLanguage)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
