package handlers

import (
	c "context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestHandlerContainer_returnAddresses(t *testing.T) {
	type fields struct {
		buildingService    *services.Buildings_mock
		userService        *services.Users_mock
		bot                *InternalBot_mock
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
		user    *tgbotapi.User
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		buildings      []services.BuildingDTO
		buildingError  error
		storedLanguage *services.Language
		userError      error
		expectedMsg    tgbotapi.MessageConfig
	}{
		{
			"a building error",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123},
			[]services.BuildingDTO{},
			errors.New("some error"),
			nil,
			nil,
			tgbotapi.NewMessage(123, "Internal error"),
		},
		{
			"no buildings - no address",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 1, user: &tgbotapi.User{ID: int64(3), LanguageCode: "en"}},
			[]services.BuildingDTO{},
			nil,
			nil,
			nil,
			tgbotapi.NewMessage(123, `Search address: 
Available building addresses and names:
No buildings were found.`),
		},
		{
			"several buildings and address, no offset",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 3, offset: 0, address: "test"},
			[]services.BuildingDTO{
				{ID: 1, Address: "test 1", NameEn: utils.GetPointer("test name 1")},
				{ID: 2, Address: "test 2", NameEn: utils.GetPointer("test name 2")},
			},
			nil,
			nil,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: 123,
					ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"test 1 - test name 1",
								`{"name":"building","id":"1"}`,
							),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"test 2 - test name 2",
								`{"name":"building","id":"2"}`,
							),
						),
					),
				},
				Text: `Search address: test
Available building addresses and names:`,
			},
		},
		{
			"several buildings and address, offset",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 3, offset: 1, address: "test"},
			[]services.BuildingDTO{
				{ID: 2, Address: "test 1", NameEn: utils.GetPointer("test name 1")},
				{ID: 3, Address: "test 2", NameEn: utils.GetPointer("test name 2")},
			},
			nil,
			nil,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: 123,
					ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"test 1 - test name 1",
								`{"name":"building","id":"2"}`,
							),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"test 2 - test name 2",
								`{"name":"building","id":"3"}`,
							),
						),
					),
				},
				Text: `Search address: test
Available building addresses and names:`,
			},
		},
		{
			"several buildings and address, offset, button",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 2, offset: 1, address: "test"},
			[]services.BuildingDTO{
				{ID: 1, Address: "test 1", NameEn: utils.GetPointer("test name 1")},
				{ID: 2, Address: "test 2", NameEn: utils.GetPointer("test name 2")},
			},
			nil,
			nil,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: 123,
					ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"test 1 - test name 1",
								`{"name":"building","id":"1"}`,
							),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"test 2 - test name 2",
								`{"name":"building","id":"2"}`,
							),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"Next 2 buildings",
								`{"name":"next","limit":2,"offset":3}`,
							),
						),
					),
				},
				Text: `Search address: test
Available building addresses and names:`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.buildingService.EXPECT().GetBuildings(
				tt.args.ctx,
				tt.args.address,
				tt.args.limit,
				tt.args.offset,
			).Return(tt.buildings, tt.buildingError)
			tt.fields.bot.EXPECT().
				Send(tt.expectedMsg).Return(tgbotapi.Message{}, nil)
			if tt.args.user != nil {
				tt.fields.userService.EXPECT().
					GetPreferredLanguage(tt.args.ctx, tt.args.user.ID).
					Return(tt.storedLanguage, tt.userError)
			}
			h := HandlerContainer{
				buildingService:    tt.fields.buildingService,
				userService:        tt.fields.userService,
				bot:                tt.fields.bot,
				HandlersPerCommand: tt.fields.HandlersPerCommand,
				handlersPerButton:  tt.fields.handlersPerButton,
				commandsForHelp:    tt.fields.commandsForHelp,
				metrics:            tt.fields.metrics,
			}
			h.returnAddresses(
				tt.args.ctx,
				tt.args.chatID,
				tt.args.user,
				tt.args.address,
				tt.args.limit,
				tt.args.offset,
			)
		})
	}
}
