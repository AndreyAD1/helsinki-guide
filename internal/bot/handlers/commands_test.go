package handlers

import (
	"context"
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
		name             string
		fields           fields
		args             args
		buildingPreviews []services.BuildingPreview
		buildingError    error
		expectedMsg      tgbotapi.MessageConfig
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
		{
			"no buildings - no address",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 1},
			[]services.BuildingPreview{},
			nil,
			tgbotapi.NewMessage(123, `Search address: 
Available building addresses and names:
End`),
		},
		{
			"several buildings and address, no offset",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 3, offset: 0, address: "test"},
			[]services.BuildingPreview{
				{Address: "test 1", Name: "test name 1"},
				{Address: "test 2", Name: "test name 2"},
			},
			nil,
			tgbotapi.NewMessage(123, `Search address: test
Available building addresses and names:
1. test 1 - test name 1
2. test 2 - test name 2
End`),
		},
		{
			"several buildings and address, offset",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 3, offset: 1, address: "test"},
			[]services.BuildingPreview{
				{Address: "test 1", Name: "test name 1"},
				{Address: "test 2", Name: "test name 2"},
			},
			nil,
			tgbotapi.NewMessage(123, `Search address: test
Available building addresses and names:
2. test 1 - test name 1
3. test 2 - test name 2
End`),
		},
		{
			"several buildings and address, offset, button",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, limit: 2, offset: 1, address: "test"},
			[]services.BuildingPreview{
				{Address: "test 1", Name: "test name 1"},
				{Address: "test 2", Name: "test name 2"},
			},
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: 123,
					ReplyMarkup: tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(
								"Next 2 buildings",
								`{"name":"next","limit":2,"offset":3}`,
							),
						),
					),
				},
				Text: `Search address: test
Available building addresses and names:
2. test 1 - test name 1
3. test 2 - test name 2`,
			},
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

func TestHandlerContainer_getBuilding(t *testing.T) {
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
		message *tgbotapi.Message
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		address       string
		buildings     []services.BuildingDTO
		buildingError error
		expectedMsg   string
	}{
		{
			"no address",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					Text: "test",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 4},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"",
			[]services.BuildingDTO{},
			nil,
			"Please add an address to this command.",
		},
		{
			"service error",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					Text: "\\address test address",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 8},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"test address",
			[]services.BuildingDTO{},
			errors.New("building error"),
			"Internal error.",
		},
		{
			"no buildings",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					Text: "\\address test address",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 8},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"test address",
			[]services.BuildingDTO{},
			nil,
			"Unfortunately, I don't know this address.",
		},
		{
			"one building en",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					Text: "\\address test address",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 8},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"test address",
			[]services.BuildingDTO{
				{
					NameEn:  utils.GetPointer("test building"),
					Address: "test address",
				},
			},
			nil,
			`Name: test building
Address: test address
Description: no data
Completion year: no data
Authors: no data
Building history: no data
Notable features: no data
Facades: no data
Interesting details: no data
Surroundings: no data`,
		},
		{
			"one building ru",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					From: &tgbotapi.User{LanguageCode: "ru"},
					Text: "\\address test address",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 8},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"test address",
			[]services.BuildingDTO{
				{
					NameEn:  utils.GetPointer("test building"),
					NameRu:  utils.GetPointer("тестовое имя"),
					Address: "test address",
				},
			},
			nil,
			`Имя: тестовое имя
Адрес: test address
Описание: нет данных
Год постройки: нет данных
Авторы: нет данных
История здания: нет данных
Примечательные особенности: нет данных
Фасады: нет данных
Интересные детали: нет данных
Окрестности: нет данных`,
		},
		{
			"one building ru",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					From: &tgbotapi.User{LanguageCode: "fi"},
					Text: "\\address test address",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 8},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"test address",
			[]services.BuildingDTO{
				{
					NameFi:  utils.GetPointer("testi rakennus"),
					NameEn:  utils.GetPointer("test building"),
					NameRu:  utils.GetPointer("тестовое имя"),
					Address: "test address",
				},
			},
			nil,
			`Nimi: testi rakennus
Katuosoite: test address
Kerrosluku: no data
Käyttöönottovuosi: no data
Suunnittelijat: no data
Rakennushistoria: no data
Huomattavia ominaisuuksia: no data
Julkisivut: no data
Erityispiirteet: no data
Ymparistonkuvaus: no data`,
		},
		{
			"two buildings",
			fields{
				services.NewBuildings_mock(t),
				newInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					Text: "\\address test address",
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Length: 8},
					},
					Chat: &tgbotapi.Chat{ID: 123},
				},
			},
			"test address",
			[]services.BuildingDTO{
				{
					NameEn:  utils.GetPointer("test building"),
					Address: "test address",
				},
				{
					NameEn:         utils.GetPointer("test building 2"),
					Address:        "test address 2",
					CompletionYear: utils.GetPointer(1973),
				},
			},
			nil,
			`Name: test building
Address: test address
Description: no data
Completion year: no data
Authors: no data
Building history: no data
Notable features: no data
Facades: no data
Interesting details: no data
Surroundings: no data

Name: test building 2
Address: test address 2
Description: no data
Completion year: 1973
Authors: no data
Building history: no data
Notable features: no data
Facades: no data
Interesting details: no data
Surroundings: no data`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.address != "" {
				tt.fields.buildingService.EXPECT().
					GetBuildingsByAddress(tt.args.ctx, tt.address).
					Return(tt.buildings, tt.buildingError)
			}

			tt.fields.bot.EXPECT().
				Send(tgbotapi.NewMessage(tt.args.message.Chat.ID, tt.expectedMsg)).
				Return(tgbotapi.Message{}, nil)
			h := HandlerContainer{
				buildingService:    tt.fields.buildingService,
				bot:                tt.fields.bot,
				HandlersPerCommand: tt.fields.HandlersPerCommand,
				handlersPerButton:  tt.fields.handlersPerButton,
				commandsForHelp:    tt.fields.commandsForHelp,
				metrics:            tt.fields.metrics,
			}
			h.getBuilding(tt.args.ctx, tt.args.message)
		})
	}
}
