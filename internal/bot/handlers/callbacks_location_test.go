package handlers

import (
	"context"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestHandlerContainer_language_unexpectedCallback(t *testing.T) {
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
			"no sender",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{}},
			},
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
		{
			"no button data",
			&tgbotapi.CallbackQuery{
				ID:      "123",
				From:    &tgbotapi.User{},
				Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{}},
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
			err := h.language(context.Background(), tt.calbackQuery)
			require.ErrorIs(t, err, ErrUnexpectedCallback)
		})
	}
}

func TestHandlerContainer_language_unexpectedButtonLanguage(t *testing.T) {
	calbackQuery := &tgbotapi.CallbackQuery{
		ID:      "123",
		From:    &tgbotapi.User{},
		Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 99}},
		Data:    `{"name":"language","value":"xx"}`,
	}
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
	err := h.language(context.Background(), calbackQuery)
	require.Error(t, err)
}
