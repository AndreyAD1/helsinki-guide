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

func TestHandlerContainer_button_errWithoutResponse(t *testing.T) {
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
