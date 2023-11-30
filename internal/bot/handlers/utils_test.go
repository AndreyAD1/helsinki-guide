package handlers

import (
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	"github.com/stretchr/testify/require"
)

var extendedBuilding = services.BuildingDTO{
	NameFi: utils.GetPointer("name fi"),
	NameEn: utils.GetPointer("name en"),
	NameRu: utils.GetPointer("name ru"),
	Address: "test address",
	CompletionYear: utils.GetPointer(2023),
	Authors: &[]string{"Author 1", "Author2"},
	HistoryFi: utils.GetPointer("history fi"),
	HistoryEn: utils.GetPointer("history en"),
	HistoryRu: utils.GetPointer("history ru"),
}

func TestSerializeIntoMessage_positive(t *testing.T) {
	dummyBuilding := services.BuildingDTO{Address: "osoite"}
	type args struct {
		object         any
		outputLanguage outputLanguage
	}
	tests := []struct {
		name    string
		args    args
		expected    string
	}{
		{
			"dummy fi",
			args{dummyBuilding, Finnish},
			`Nimi: no data
Katuosoite: osoite
Käyttöönottovuosi: no data
Suunnittelijat: no data
Rakennushistoria: no data`,
		},
		{
			"dummy en",
			args{dummyBuilding, English},
			`Name: no data
Address: osoite
Completion year: no data
Authors: no data
Building history: no data`,
		},
		{
			"dummy ru",
			args{dummyBuilding, Russian},
			`Имя: нет данных
Адрес: osoite
Год постройки: нет данных
Авторы: нет данных
История здания: нет данных`,
		},
		{
			"extended fi",
			args{extendedBuilding, Finnish},
			`Nimi: name fi
Katuosoite: test address
Käyttöönottovuosi: 2023
Suunnittelijat: Author 1, Author2
Rakennushistoria: history fi`,
		},
		{
			"extended en",
			args{extendedBuilding, English},
			`Name: name en
Address: test address
Completion year: 2023
Authors: Author 1, Author2
Building history: history en`,
		},
		{
			"extended ru",
			args{extendedBuilding, Russian},
			`Имя: name ru
Адрес: test address
Год постройки: 2023
Авторы: Author 1, Author2
История здания: history ru`,
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
	type args struct {
		object         any
		outputLanguage outputLanguage
	}
	tests := []struct {
		name    string
		args    args
		expectedError    error
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
			"a structure with unexpected field tag",
			args{testTypeStruct{'A'}, English},
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