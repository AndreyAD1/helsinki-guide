package handlers

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestHandlerContainer_building_errWithoutResponse(t *testing.T) {
	tests := []struct {
		name         string
		calbackQuery *tgbotapi.CallbackQuery
		queryID      string
	}{
		{
			"empty callback query",
			&tgbotapi.CallbackQuery{ID: "123"},
			"123",
		},
		{
			"no chat",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{},
				Message: &tgbotapi.Message{},
			},
			"123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			botMock := NewInternalBot_mock(t)
			botMock.EXPECT().
				Request(tgbotapi.NewCallback(tt.queryID, "")).Return(nil, nil)
			h := HandlerContainer{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				botMock,
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				metrics.NewMetrics(prometheus.NewRegistry()),
				map[string]CommandHandler{},
			}
			err := h.building(context.Background(), tt.calbackQuery)
			require.ErrorIs(t, err, ErrUnexpectedCallback)
		})
	}
}

func TestHandlerContainer_building_unexpectedButtonData(t *testing.T) {
	calbackQuery := &tgbotapi.CallbackQuery{
		ID:      "123",
		From:    &tgbotapi.User{},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
	}
	tests := []struct {
		name       string
		buttonData string
	}{
		{"no data", ""},
		{"unexpected data", `{"unknown json": 123}`},
		{"unexpected building ID", `{"name": "building", "id": "123.3456"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			botMock := NewInternalBot_mock(t)
			botMock.EXPECT().
				Send(tgbotapi.NewMessage(calbackQuery.Message.Chat.ID, "Internal error")).
				Return(tgbotapi.Message{}, nil).
				On("Request", tgbotapi.NewCallback(calbackQuery.ID, "")).Return(nil, nil)
			h := HandlerContainer{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				botMock,
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				metrics.NewMetrics(prometheus.NewRegistry()),
				map[string]CommandHandler{},
			}
			calbackQuery.Data = tt.buttonData
			err := h.building(context.Background(), calbackQuery)
			require.Error(t, err)
		})
	}
}

func TestHandlerContainer_building_serviceError(t *testing.T) {
	calbackQuery := &tgbotapi.CallbackQuery{
		ID:      "123",
		From:    &tgbotapi.User{},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
		Data:    `{"name": "building", "id": "123"}`,
	}
	botMock := NewInternalBot_mock(t)
	buildingMock := services.NewBuildings_mock(t)
	botMock.EXPECT().
		Send(tgbotapi.NewMessage(calbackQuery.Message.Chat.ID, "Internal error")).
		Return(tgbotapi.Message{}, nil).
		On("Request", tgbotapi.NewCallback(calbackQuery.ID, "")).
		Return(nil, nil)
	ctx := context.Background()
	buildingMock.EXPECT().GetBuildingByID(ctx, int64(123)).Return(nil, errors.New("test"))
	h := HandlerContainer{
		buildingMock,
		services.NewUsers_mock(t),
		botMock,
		map[string]CommandHandler{},
		map[string]internalButtonHandler{},
		"",
		metrics.NewMetrics(prometheus.NewRegistry()),
		map[string]CommandHandler{},
	}
	err := h.building(ctx, calbackQuery)
	require.Error(t, err)
}

func TestHandlerContainer_building_noBuilding(t *testing.T) {
	calbackQuery := &tgbotapi.CallbackQuery{
		ID:      "123",
		From:    &tgbotapi.User{},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
		Data:    `{"name": "building", "id": "123"}`,
	}
	botMock := NewInternalBot_mock(t)
	buildingMock := services.NewBuildings_mock(t)
	botMock.EXPECT().
		Send(tgbotapi.NewMessage(calbackQuery.Message.Chat.ID, "Can not find the building.")).
		Return(tgbotapi.Message{}, nil).
		On("Request", tgbotapi.NewCallback(calbackQuery.ID, "")).
		Return(nil, nil)
	ctx := context.Background()
	buildingMock.EXPECT().GetBuildingByID(ctx, int64(123)).Return(nil, nil)
	h := HandlerContainer{
		buildingMock,
		services.NewUsers_mock(t),
		botMock,
		map[string]CommandHandler{},
		map[string]internalButtonHandler{},
		"",
		metrics.NewMetrics(prometheus.NewRegistry()),
		map[string]CommandHandler{},
	}
	err := h.building(ctx, calbackQuery)
	require.ErrorIs(t, err, ErrUnexpectedCallback)
}

func TestHandlerContainer_building_serializationError(t *testing.T) {
	calbackQuery := &tgbotapi.CallbackQuery{
		ID:      "123",
		From:    &tgbotapi.User{ID: 555},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
		Data:    `{"name": "building", "id": "123"}`,
	}
	botMock := NewInternalBot_mock(t)
	buildingMock := services.NewBuildings_mock(t)
	userMock := services.NewUsers_mock(t)
	botMock.EXPECT().
		Send(tgbotapi.NewMessage(calbackQuery.Message.Chat.ID, "Internal error.")).
		Return(tgbotapi.Message{}, nil).
		On("Request", tgbotapi.NewCallback(calbackQuery.ID, "")).
		Return(nil, nil)
	ctx := context.Background()
	buildingMock.EXPECT().GetBuildingByID(ctx, int64(123)).
		Return(&services.BuildingDTO{}, nil)
	l := services.Language("unknown")
	userMock.EXPECT().GetPreferredLanguage(ctx, calbackQuery.From.ID).
		Return(&l, nil)
	h := HandlerContainer{
		buildingMock,
		userMock,
		botMock,
		map[string]CommandHandler{},
		map[string]internalButtonHandler{},
		"",
		metrics.NewMetrics(prometheus.NewRegistry()),
		map[string]CommandHandler{},
	}
	err := h.building(ctx, calbackQuery)
	require.Error(t, err)
}

func TestHandlerContainer_building_ok_noLanguageCheck(t *testing.T) {
	callbackQuery := &tgbotapi.CallbackQuery{
		ID:      "123",
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
		Data:    `{"name": "building", "id": "123"}`,
	}
	botMock := NewInternalBot_mock(t)
	buildingMock := services.NewBuildings_mock(t)
	userMock := services.NewUsers_mock(t)
	expectedMessage := tgbotapi.NewMessage(
		callbackQuery.Message.Chat.ID,
		`<b>Name:</b> no data
<b>Address:</b> test address
<b>Description:</b> no data
<b>Completion year:</b> no data
<b>Authors:</b> no data
<b>Facades:</b> no data
<b>Interesting details:</b> no data
<b>Notable features:</b> no data
<b>Surroundings:</b> no data
<b>Building history:</b> no data`,
	)
	expectedMessage.ParseMode = tgbotapi.ModeHTML
	botMock.EXPECT().
		Send(expectedMessage).
		Return(tgbotapi.Message{}, nil).
		On("Request", tgbotapi.NewCallback(callbackQuery.ID, "")).
		Return(nil, nil)
	ctx := context.Background()
	buildingMock.EXPECT().GetBuildingByID(ctx, int64(123)).
		Return(&services.BuildingDTO{Address: "test address"}, nil)
	h := HandlerContainer{
		buildingMock,
		userMock,
		botMock,
		map[string]CommandHandler{},
		map[string]internalButtonHandler{},
		"",
		metrics.NewMetrics(prometheus.NewRegistry()),
		map[string]CommandHandler{},
	}
	err := h.building(ctx, callbackQuery)
	require.NoError(t, err)
}

func TestHandlerContainer_building_ok(t *testing.T) {
	tests := []struct {
		name              string
		callbackQuery     *tgbotapi.CallbackQuery
		expectedMessage   tgbotapi.MessageConfig
		building          *services.BuildingDTO
		buildingError     error
		preferredLanguage *services.Language
		languageError     error
	}{
		{
			"default en",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{ID: 555, LanguageCode: "en"},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
				Data:    `{"name": "building", "id": "123"}`,
			},
			tgbotapi.NewMessage(
				99,
				`<b>Name:</b> no data
<b>Address:</b> test address
<b>Description:</b> no data
<b>Completion year:</b> no data
<b>Authors:</b> no data
<b>Facades:</b> no data
<b>Interesting details:</b> no data
<b>Notable features:</b> no data
<b>Surroundings:</b> no data
<b>Building history:</b> no data`,
			),
			&services.BuildingDTO{Address: "test address"},
			nil,
			nil,
			nil,
		},
		{
			"default ru",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{ID: 555, LanguageCode: "ru"},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
				Data:    `{"name": "building", "id": "123"}`,
			},
			tgbotapi.NewMessage(
				99,
				`<b>Имя:</b> тестовое имя
<b>Адрес:</b> test address
<b>Описание:</b> нет данных
<b>Год постройки:</b> нет данных
<b>Авторы:</b> нет данных
<b>Фасады:</b> нет данных
<b>Интересные детали:</b> нет данных
<b>Примечательные особенности:</b> нет данных
<b>Окрестности:</b> нет данных
<b>История здания:</b> нет данных`,
			),
			&services.BuildingDTO{
				Address: "test address",
				NameRu:  utils.GetPointer("тестовое имя"),
			},
			nil,
			nil,
			nil,
		},
		{
			"default fi",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{ID: 555, LanguageCode: "fi"},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
				Data:    `{"name": "building", "id": "123"}`,
			},
			tgbotapi.NewMessage(
				99,
				`<b>Nimi:</b> testi rakennus
<b>Katuosoite:</b> test address
<b>Kerrosluku:</b> ei tietoja
<b>Käyttöönottovuosi:</b> ei tietoja
<b>Suunnittelijat:</b> ei tietoja
<b>Julkisivut:</b> ei tietoja
<b>Erityispiirteet:</b> ei tietoja
<b>Huomattavia ominaisuuksia:</b> ei tietoja
<b>Ympäristönkuvaus:</b> ei tietoja
<b>Rakennushistoria:</b> ei tietoja`,
			),
			&services.BuildingDTO{
				Address: "test address",
				NameFi:  utils.GetPointer("testi rakennus"),
				NameEn:  utils.GetPointer("test building"),
				NameRu:  utils.GetPointer("тестовое имя"),
			},
			nil,
			nil,
			nil,
		},
		{
			"unknown default language",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{ID: 555, LanguageCode: "unknown"},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
				Data:    `{"name": "building", "id": "123"}`,
			},
			tgbotapi.NewMessage(
				99,
				`<b>Name:</b> test building
<b>Address:</b> test address
<b>Description:</b> no data
<b>Completion year:</b> no data
<b>Authors:</b> no data
<b>Facades:</b> no data
<b>Interesting details:</b> no data
<b>Notable features:</b> no data
<b>Surroundings:</b> no data
<b>Building history:</b> no data`,
			),
			&services.BuildingDTO{
				Address: "test address",
				NameFi:  utils.GetPointer("testi rakennus"),
				NameEn:  utils.GetPointer("test building"),
				NameRu:  utils.GetPointer("тестовое имя"),
			},
			nil,
			nil,
			nil,
		},
		{
			"preferred Finnish",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{ID: 555, LanguageCode: "en"},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
				Data:    `{"name": "building", "id": "123"}`,
			},
			tgbotapi.NewMessage(
				99,
				`<b>Nimi:</b> testi rakennus
<b>Katuosoite:</b> test address
<b>Kerrosluku:</b> ei tietoja
<b>Käyttöönottovuosi:</b> ei tietoja
<b>Suunnittelijat:</b> ei tietoja
<b>Julkisivut:</b> ei tietoja
<b>Erityispiirteet:</b> ei tietoja
<b>Huomattavia ominaisuuksia:</b> ei tietoja
<b>Ympäristönkuvaus:</b> ei tietoja
<b>Rakennushistoria:</b> ei tietoja`,
			),
			&services.BuildingDTO{
				Address: "test address",
				NameFi:  utils.GetPointer("testi rakennus"),
				NameEn:  utils.GetPointer("test building"),
				NameRu:  utils.GetPointer("тестовое имя"),
			},
			nil,
			&services.Finnish,
			nil,
		},
		{
			"language service error",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{ID: 555, LanguageCode: "ru"},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
				Data:    `{"name": "building", "id": "123"}`,
			},
			tgbotapi.NewMessage(
				99,
				`<b>Имя:</b> тестовое имя
<b>Адрес:</b> test address
<b>Описание:</b> нет данных
<b>Год постройки:</b> нет данных
<b>Авторы:</b> нет данных
<b>Фасады:</b> нет данных
<b>Интересные детали:</b> нет данных
<b>Примечательные особенности:</b> нет данных
<b>Окрестности:</b> нет данных
<b>История здания:</b> нет данных`,
			),
			&services.BuildingDTO{
				Address: "test address",
				NameFi:  utils.GetPointer("testi rakennus"),
				NameEn:  utils.GetPointer("test building"),
				NameRu:  utils.GetPointer("тестовое имя"),
			},
			nil,
			nil,
			errors.New("some language error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			botMock := NewInternalBot_mock(t)
			buildingMock := services.NewBuildings_mock(t)
			userMock := services.NewUsers_mock(t)
			tt.expectedMessage.ParseMode = tgbotapi.ModeHTML
			botMock.EXPECT().
				Send(tt.expectedMessage).
				Return(tgbotapi.Message{}, nil).
				On("Request", tgbotapi.NewCallback(tt.callbackQuery.ID, "")).
				Return(nil, nil)
			ctx := context.Background()
			id, err := strconv.ParseInt(tt.callbackQuery.ID, 0, 64)
			require.NoError(t, err)
			buildingMock.EXPECT().GetBuildingByID(ctx, id).
				Return(tt.building, tt.buildingError)
			userMock.EXPECT().GetPreferredLanguage(ctx, tt.callbackQuery.From.ID).
				Return(tt.preferredLanguage, tt.languageError)
			h := HandlerContainer{
				buildingMock,
				userMock,
				botMock,
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				metrics.NewMetrics(prometheus.NewRegistry()),
				map[string]CommandHandler{},
			}
			err = h.building(ctx, tt.callbackQuery)
			require.NoError(t, err)
		})
	}
}
