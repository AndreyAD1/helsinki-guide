package handlers

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
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
		chatID    int64
		latitude  float64
		longitude float64
	}
	serviceError := errors.New("some error")
	tests := []struct {
		name             string
		fields           fields
		args             args
		buildingPreviews []services.BuildingDTO
		buildingError    error
		expectedMsg      tgbotapi.MessageConfig
		expectedError    error
	}{
		{
			"building error",
			fields{
				services.NewBuildings_mock(t),
				NewInternalBot_mock(t),
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				nil,
			},
			args{chatID: 123, latitude: 3, longitude: 3},
			[]services.BuildingDTO{},
			serviceError,
			tgbotapi.NewMessage(123, "Internal error"),
			serviceError,
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
			args{chatID: 123, latitude: 3, longitude: 3},
			[]services.BuildingDTO{},
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
			},
			nil,
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
			args{chatID: 123, latitude: 3, longitude: 3},
			[]services.BuildingDTO{
				{
					Address: "test 1",
					NameEn:  utils.GetPointer("test name 1"),
					ID:      999,
				},
			},
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           123,
					ReplyToMessageID: 0,
					ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
							tgbotapi.InlineKeyboardButton{
								Text:         "test 1 - test name 1",
								CallbackData: utils.GetPointer(`{"name":"building","id":"999"}`),
							},
						}},
					},
				},
				Text:                  fmt.Sprintf(nearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
				DisableWebPagePreview: false,
			},
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
			[]services.BuildingDTO{
				{Address: "test 1", NameEn: utils.GetPointer("test name 1"), ID: 1000},
				{Address: "test 2", NameEn: utils.GetPointer("test name 2"), ID: 999},
			},
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           123,
					ReplyToMessageID: 0,
					ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
							tgbotapi.InlineKeyboardButton{
								Text:         "test 1 - test name 1",
								CallbackData: utils.GetPointer(`{"name":"building","id":"1000"}`),
							}},
							{tgbotapi.InlineKeyboardButton{
								Text:         "test 2 - test name 2",
								CallbackData: utils.GetPointer(`{"name":"building","id":"999"}`),
							}},
						},
					},
				},
				Text:                  fmt.Sprintf(nearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
				DisableWebPagePreview: false,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tt.fields.buildingService.EXPECT().GetNearestBuildings(
				ctx,
				DEFAULT_DISTANCE,
				tt.args.latitude,
				tt.args.longitude,
				defaultLimit,
				0,
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
			message := tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: tt.args.chatID},
			}
			if tt.args.latitude != 0 {
				message.Location = &tgbotapi.Location{
					Latitude:  tt.args.latitude,
					Longitude: tt.args.longitude,
				}
			}
			err := h.getNearestAddresses(ctx, &message)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestHandlerContainer_getNearestAddresses_languages(t *testing.T) {
	type args struct {
		userID       int64
		userLanguage string
	}
	tests := []struct {
		name           string
		args           args
		buildings      []services.BuildingDTO
		storedLanguage *services.Language
		userError      error
		expectedMsg    tgbotapi.MessageConfig
		expectedError  error
	}{
		{
			"no buildings - English - no configured language",
			args{userID: int64(123), userLanguage: "en"},
			[]services.BuildingDTO{},
			nil,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		{
			"no buildings - French - no configured language",
			args{userID: int64(123), userLanguage: "fr"},
			[]services.BuildingDTO{},
			nil,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		{
			"no buildings - Russian - no configured language",
			args{userID: int64(123), userLanguage: "ru"},
			[]services.BuildingDTO{},
			nil,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsRussianTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		{
			"no buildings - English - configured English",
			args{userID: int64(123), userLanguage: "en"},
			[]services.BuildingDTO{},
			&services.English,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		{
			"no buildings - English - configured Finnish",
			args{userID: int64(123), userLanguage: "en"},
			[]services.BuildingDTO{},
			&services.Finnish,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsFinnishTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		{
			"no buildings - Finnish - configured Russian",
			args{userID: int64(123), userLanguage: "fi"},
			[]services.BuildingDTO{},
			&services.Russian,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsRussianTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		{
			"no buildings - French - configured English",
			args{userID: int64(123), userLanguage: "fr"},
			[]services.BuildingDTO{},
			&services.English,
			nil,
			tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: 123},
				Text:     fmt.Sprintf(noNearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
			},
			nil,
		},
		// {
		// 	"one building",
		// 	args{chatID: 123, latitude: 3, longitude: 3},
		// 	[]services.BuildingDTO{
		// 		{
		// 			Address: "test 1",
		// 			NameEn:  utils.GetPointer("test name 1"),
		// 			ID:      999,
		// 		},
		// 	},
		// 	tgbotapi.MessageConfig{
		// 		BaseChat: tgbotapi.BaseChat{
		// 			ChatID:           123,
		// 			ReplyToMessageID: 0,
		// 			ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
		// 				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
		// 					tgbotapi.InlineKeyboardButton{
		// 						Text:         "test 1 - test name 1",
		// 						CallbackData: utils.GetPointer(`{"name":"building","id":"999"}`),
		// 					},
		// 				}},
		// 			},
		// 		},
		// 		Text:                  fmt.Sprintf(nearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
		// 		DisableWebPagePreview: false,
		// 	},
		// 	nil,
		// },
		// {
		// 	"two buildings",
		// 	args{chatID: 123, latitude: 3, longitude: 0},
		// 	[]services.BuildingDTO{
		// 		{Address: "test 1", NameEn: utils.GetPointer("test name 1"), ID: 1000},
		// 		{Address: "test 2", NameEn: utils.GetPointer("test name 2"), ID: 999},
		// 	},
		// 	tgbotapi.MessageConfig{
		// 		BaseChat: tgbotapi.BaseChat{
		// 			ChatID:           123,
		// 			ReplyToMessageID: 0,
		// 			ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
		// 				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
		// 					tgbotapi.InlineKeyboardButton{
		// 						Text:         "test 1 - test name 1",
		// 						CallbackData: utils.GetPointer(`{"name":"building","id":"1000"}`),
		// 					}},
		// 					{tgbotapi.InlineKeyboardButton{
		// 						Text:         "test 2 - test name 2",
		// 						CallbackData: utils.GetPointer(`{"name":"building","id":"999"}`),
		// 					}},
		// 				},
		// 			},
		// 		},
		// 		Text:                  fmt.Sprintf(nearestBuildingsEnglishTemplate, DEFAULT_DISTANCE),
		// 		DisableWebPagePreview: false,
		// 	},
		// 	nil,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildingService := services.NewBuildings_mock(t)
			userService := services.NewUsers_mock(t)
			bot := NewInternalBot_mock(t)

			ctx := context.Background()
			buildingService.EXPECT().GetNearestBuildings(
				ctx,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(tt.buildings, nil)
			bot.EXPECT().Send(tt.expectedMsg).Return(tgbotapi.Message{}, nil)

			userService.EXPECT().GetPreferredLanguage(ctx, tt.args.userID).
				Return(tt.storedLanguage, tt.userError)

			h := HandlerContainer{
				buildingService:    buildingService,
				userService:        userService,
				bot:                bot,
				HandlersPerCommand: map[string]CommandHandler{},
				handlersPerButton:  map[string]internalButtonHandler{},
				commandsForHelp:    "",
				metrics:            nil,
			}
			message := tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: int64(123)},
				From: &tgbotapi.User{ID: tt.args.userID, LanguageCode: tt.args.userLanguage},
				Location: &tgbotapi.Location{
					Latitude:  float64(60),
					Longitude: float64(30),
				},
			}
			err := h.getNearestAddresses(ctx, &message)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
