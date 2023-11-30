package handlers

import (
	"fmt"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/stretchr/testify/require"
)

func Test_SerializeIntoMessage(t *testing.T) {
	test_name_fi := "fin name"
	address := "osoite"
	expectedResult := fmt.Sprintf(
		`Nimi: %v
Katuosoite: %v
Käyttöönottovuosi: no data
Suunnittelijat: no data
Rakennushistoria: no data`,
		test_name_fi,
		address,
	)
	dto := services.BuildingDTO{NameFi: &test_name_fi, Address: address}
	result, err := SerializeIntoMessage(dto, "fi")
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestSerializeIntoMessage(t *testing.T) {
	type args struct {
		object         any
		outputLanguage outputLanguage
	}
	tests := []struct {
		name    string
		args    args
		expected    string
		expectedErr bool
	}{
		{
			"dummy en",
			args{services.BuildingDTO{Address: "osoite"}, Finnish},
			`Nimi: no data
Katuosoite: osoite
Käyttöönottovuosi: no data
Suunnittelijat: no data
Rakennushistoria: no data`,
			false,
		},
		{
			"dummy en",
			args{services.BuildingDTO{Address: "osoite"}, English},
			`Name: no data
Address: osoite
Completion year: no data
Authors: no data
Building history: no data`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SerializeIntoMessage(tt.args.object, tt.args.outputLanguage)
			if tt.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}
