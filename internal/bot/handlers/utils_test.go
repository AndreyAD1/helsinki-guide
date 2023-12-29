package handlers

import (
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/require"
)

var extendedBuilding = services.BuildingDTO{
	NameFi:            utils.GetPointer("name fi"),
	NameEn:            utils.GetPointer("name en"),
	NameRu:            utils.GetPointer("name ru"),
	Address:           "test address",
	DescriptionFi:     utils.GetPointer("description fi"),
	DescriptionEn:     utils.GetPointer("description en"),
	DescriptionRu:     utils.GetPointer("description ru"),
	CompletionYear:    utils.GetPointer(2023),
	Authors:           &[]string{"Author 1", "Author2"},
	HistoryFi:         utils.GetPointer("history fi"),
	HistoryEn:         utils.GetPointer("history en"),
	HistoryRu:         utils.GetPointer("history ru"),
	NotableFeaturesFi: utils.GetPointer("features fi"),
	NotableFeaturesEn: utils.GetPointer("features en"),
	NotableFeaturesRu: utils.GetPointer("features ru"),
	FacadesFi:         utils.GetPointer("facades fi"),
	FacadesEn:         utils.GetPointer("facades en"),
	FacadesRu:         utils.GetPointer("facades ru"),
	DetailsFi:         utils.GetPointer("details fi"),
	DetailsEn:         utils.GetPointer("details en"),
	DetailsRu:         utils.GetPointer("details ru"),
	SurroundingsFi:    utils.GetPointer("surroundings fi"),
	SurroundingsEn:    utils.GetPointer("surroundings en"),
	SurroundingsRu:    utils.GetPointer("surroundings ru"),
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
			`<b>Nimi:</b> no data
<b>Katuosoite:</b> osoite
<b>Kerrosluku:</b> no data
<b>Käyttöönottovuosi:</b> no data
<b>Suunnittelijat:</b> no data
<b>Julkisivut:</b> no data
<b>Erityispiirteet:</b> no data
<b>Huomattavia ominaisuuksia:</b> no data
<b>Ymparistonkuvaus:</b> no data
<b>Rakennushistoria:</b> no data`,
		},
		{
			"dummy en",
			args{dummyBuilding, English},
			`<b>Name:</b> no data
<b>Address:</b> osoite
<b>Description:</b> no data
<b>Completion year:</b> no data
<b>Authors:</b> no data
<b>Facades:</b> no data
<b>Interesting details:</b> no data
<b>Notable features:</b> no data
<b>Surroundings:</b> no data
<b>Building history:</b> no data`,
		},
		{
			"dummy ru",
			args{dummyBuilding, Russian},
			`<b>Имя:</b> нет данных
<b>Адрес:</b> osoite
<b>Описание:</b> нет данных
<b>Год постройки:</b> нет данных
<b>Авторы:</b> нет данных
<b>Фасады:</b> нет данных
<b>Интересные детали:</b> нет данных
<b>Примечательные особенности:</b> нет данных
<b>Окрестности:</b> нет данных
<b>История здания:</b> нет данных`,
		},
		{
			"extended fi",
			args{extendedBuilding, Finnish},
			`<b>Nimi:</b> name fi
<b>Katuosoite:</b> test address
<b>Kerrosluku:</b> description fi
<b>Käyttöönottovuosi:</b> 2023
<b>Suunnittelijat:</b> Author 1, Author2
<b>Julkisivut:</b> facades fi
<b>Erityispiirteet:</b> details fi
<b>Huomattavia ominaisuuksia:</b> features fi
<b>Ymparistonkuvaus:</b> surroundings fi
<b>Rakennushistoria:</b> history fi`,
		},
		{
			"extended en",
			args{extendedBuilding, English},
			`<b>Name:</b> name en
<b>Address:</b> test address
<b>Description:</b> description en
<b>Completion year:</b> 2023
<b>Authors:</b> Author 1, Author2
<b>Facades:</b> facades en
<b>Interesting details:</b> details en
<b>Notable features:</b> features en
<b>Surroundings:</b> surroundings en
<b>Building history:</b> history en`,
		},
		{
			"extended ru",
			args{extendedBuilding, Russian},
			`<b>Имя:</b> name ru
<b>Адрес:</b> test address
<b>Описание:</b> description ru
<b>Год постройки:</b> 2023
<b>Авторы:</b> Author 1, Author2
<b>Фасады:</b> facades ru
<b>Интересные детали:</b> details ru
<b>Примечательные особенности:</b> features ru
<b>Окрестности:</b> surroundings ru
<b>История здания:</b> history ru`,
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
