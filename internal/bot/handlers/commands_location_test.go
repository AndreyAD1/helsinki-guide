package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/require"
)

func TestHandlerContainer_getNearestAddresses(t *testing.T) {
	type fields struct {
		buildingService    *services.Buildings_mock
		bot                *InternalBot_mock
		HandlersPerCommand map[string]CommandHandler
		handlersPerButton  map[string]internalButtonHandler
		commandsForHelp    string
		metrics            *metrics.Metrics
	}
	type args struct {
		chatID  int64
		latitude float64
		longitude float64
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		buildingPreviews []services.BuildingPreview
		buildingError    error
		expectedMsg      tgbotapi.MessageConfig
		expectedError    error
	}{
		{
			"a building error",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123},
			[]services.BuildingPreview{},
			errors.New("some error"),
			tgbotapi.NewMessage(123, "Internal error"),
			nil,
		},
		{
			"no location",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123},
			[]services.BuildingPreview{},
			nil,
			tgbotapi.MessageConfig{},
			ErrNoLocation,
		},
		{
			"one building",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, latitude: 3, longitude: 0},
			[]services.BuildingPreview{
				{Address: "test 1", Name: "test name 1"},
			},
			nil,
			tgbotapi.NewMessage(123, `Nearest buildings:
1. test 1 - test name 1`),
			nil,
		},
		{
			"two buildings",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, latitude: 3, longitude: 0},
			[]services.BuildingPreview{
				{Address: "test 1", Name: "test name 1"},
				{Address: "test 2", Name: "test name 2"},
			},
			nil,
			tgbotapi.NewMessage(123, `Nearest buildings:
2. test 1 - test name 1
3. test 2 - test name 2`,
			),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.expectedError == nil {
				tt.fields.buildingService.EXPECT().GetNearestBuildingPreviews(
					ctx,
					DEFAULT_DISTANCE,
					tt.args.latitude,
					tt.args.longitude,
					defaultLimit,
					0,
				).Return(tt.buildingPreviews, tt.buildingError)
				
				tt.fields.bot.EXPECT().
				Send(tt.expectedMsg).Return(tgbotapi.Message{}, nil)
			}
			h := HandlerContainer{
				buildingService:    tt.fields.buildingService,
				bot:                tt.fields.bot,
				HandlersPerCommand: tt.fields.HandlersPerCommand,
				handlersPerButton:  tt.fields.handlersPerButton,
				commandsForHelp:    tt.fields.commandsForHelp,
				metrics:            tt.fields.metrics,
			}
			message := tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: tt.args.chatID},
			}
			if tt.args.latitude != 0 {
				message.Location = &tgbotapi.Location{
					Latitude: tt.args.latitude,
					Longitude: tt.args.longitude,
				}
			}
			err := h.getNearestAddresses(ctx, &message)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}