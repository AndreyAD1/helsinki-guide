package handlers

import (
	c "context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestHandlerContainer_returnAddresses(t *testing.T) {
	type fields struct {
		buildingService    *services.Buildings_mock
		bot                *internalBot_mock
		HandlersPerCommand map[string]CommandHandler
		handlersPerButton  map[string]internalButtonHandler
		commandsForHelp    string
		metrics            *metrics.Metrics
	}
	type args struct {
		ctx     c.Context
		chatID  int64
		address string
		limit   int
		offset  int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		buildingPreviews []services.BuildingPreview
		buildingError error
		expectedMsg tgbotapi.MessageConfig
	}{
		{
			"a building error",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123},
			[]services.BuildingPreview{},
			errors.New("some error"),
			tgbotapi.NewMessage(123, "Internal error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.buildingService.EXPECT().GetBuildingPreviews(
				tt.args.ctx,
				tt.args.address,
				tt.args.limit,
				tt.args.offset,
			).Return(tt.buildingPreviews, tt.buildingError)
			tt.fields.bot.EXPECT().
			Send(tt.expectedMsg).Return(tgbotapi.Message{}, nil)
			h := HandlerContainer{
				buildingService:    tt.fields.buildingService,
				bot:                tt.fields.bot,
				HandlersPerCommand: tt.fields.HandlersPerCommand,
				handlersPerButton:  tt.fields.handlersPerButton,
				commandsForHelp:    tt.fields.commandsForHelp,
				metrics:            tt.fields.metrics,
			}
			h.returnAddresses(tt.args.ctx, tt.args.chatID, tt.args.address, tt.args.limit, tt.args.offset)
		})
	}
}
