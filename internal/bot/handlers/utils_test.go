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
Rakennushistoria: no data`,
		test_name_fi,
		address,
	)
	dto := services.BuildingDTO{NameFi: &test_name_fi, Address: address}
	result, err := SerializeIntoMessage(dto, "fi")
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}
