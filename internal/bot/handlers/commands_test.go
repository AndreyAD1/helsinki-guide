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
		},
		{
			"no buildings - no address",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
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
				NewInternalBot_mock(t),
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
				NewInternalBot_mock(t),
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
				NewInternalBot_mock(t),
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

func TestHandlerContainer_getBuilding_noLanguageCheck(t *testing.T) {
	type fields struct {
		buildingService    *services.Buildings_mock
		bot                *InternalBot_mock
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
		name              string
		fields            fields
		args              args
		address           string
		buildings         []services.BuildingDTO
		buildingError     error
		expectedMsg       string
		expectedParseMode string
	}{
		{
			"no address",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
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
			"",
		},
		{
			"service error",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
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
			"",
		},
		{
			"no buildings",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
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
			tgbotapi.ModeHTML,
		},
		{
			"one building en",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
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
			tgbotapi.ModeHTML,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.address != "" {
				tt.fields.buildingService.EXPECT().
					GetBuildingsByAddress(tt.args.ctx, tt.address).
					Return(tt.buildings, tt.buildingError)
			}
			expectedMessage := tgbotapi.NewMessage(tt.args.message.Chat.ID, tt.expectedMsg)
			expectedMessage.ParseMode = tt.expectedParseMode
			tt.fields.bot.EXPECT().
				Send(expectedMessage).
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

func TestHandlerContainer_getBuilding_withLanguageCheck(t *testing.T) {
	type fields struct {
		buildingService    *services.Buildings_mock
		userService       *services.Users_mock
		bot                *InternalBot_mock
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
		name              string
		fields            fields
		args              args
		address           string
		buildings         []services.BuildingDTO
		buildingError     error
		preferredLanguage *services.Language
		languageError     error
		expectedMsg       string
		expectedParseMode string
	}{
		{
			"one building, default en",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					From: &tgbotapi.User{LanguageCode: "en"},
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
			nil,
			nil,
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
			tgbotapi.ModeHTML,
		},
		{
			"one building, default ru",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
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
			nil,
			nil,
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
			tgbotapi.ModeHTML,
		},
		{
			"one building, default fi",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
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
			nil,
			nil,
			`<b>Nimi:</b> testi rakennus
<b>Katuosoite:</b> test address
<b>Kerrosluku:</b> no data
<b>Käyttöönottovuosi:</b> no data
<b>Suunnittelijat:</b> no data
<b>Julkisivut:</b> no data
<b>Erityispiirteet:</b> no data
<b>Huomattavia ominaisuuksia:</b> no data
<b>Ymparistonkuvaus:</b> no data
<b>Rakennushistoria:</b> no data`,
			tgbotapi.ModeHTML,
		},
		{
			"two buildings, unknown default language",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{
				context.Background(),
				&tgbotapi.Message{
					From: &tgbotapi.User{LanguageCode: "unknown language"},
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
			nil,
			nil,
			`<b>Name:</b> test building
<b>Address:</b> test address
<b>Description:</b> no data
<b>Completion year:</b> no data
<b>Authors:</b> no data
<b>Facades:</b> no data
<b>Interesting details:</b> no data
<b>Notable features:</b> no data
<b>Surroundings:</b> no data
<b>Building history:</b> no data

<b>Name:</b> test building 2
<b>Address:</b> test address 2
<b>Description:</b> no data
<b>Completion year:</b> 1973
<b>Authors:</b> no data
<b>Facades:</b> no data
<b>Interesting details:</b> no data
<b>Notable features:</b> no data
<b>Surroundings:</b> no data
<b>Building history:</b> no data`,
			tgbotapi.ModeHTML,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.buildingService.EXPECT().
				GetBuildingsByAddress(tt.args.ctx, tt.address).
				Return(tt.buildings, tt.buildingError)
			tt.fields.userService.EXPECT().
				GetPreferredLanguage(tt.args.ctx, tt.args.message.From.ID).
				Return(tt.preferredLanguage, tt.languageError)
			expectedMessage := tgbotapi.NewMessage(tt.args.message.Chat.ID, tt.expectedMsg)
			expectedMessage.ParseMode = tt.expectedParseMode
			tt.fields.bot.EXPECT().
				Send(expectedMessage).
				Return(tgbotapi.Message{}, nil)
			h := HandlerContainer{
				buildingService:    tt.fields.buildingService,
				userService:        tt.fields.userService,
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