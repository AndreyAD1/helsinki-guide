package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
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
