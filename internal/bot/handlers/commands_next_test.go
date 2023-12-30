package handlers

import (
	c "context"
	"errors"
	"testing"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
)

func TestHandlerContainer_next_positive(t *testing.T) {
	type fields struct {
		buildingService *services.Buildings_mock
		userService     *services.Users_mock
		bot             *InternalBot_mock
	}
	type args struct {
		ctx   c.Context
		query *tgbotapi.CallbackQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		queryID string
		address string
		limit   int
		offset  int
	}{
		{
			"valid text",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
			},
			args{
				c.Background(),
				&tgbotapi.CallbackQuery{
					ID: "123",
					Message: &tgbotapi.Message{
						Chat: &tgbotapi.Chat{},
						Text: "address: test address   \n",
					},
					Data: `{"name":"next","limit":2,"offset":3}`,
				},
			},
			"123",
			"test address",
			2,
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.buildingService.EXPECT().
				GetBuildingPreviews(tt.args.ctx, tt.address, tt.limit, tt.offset).
				Return([]services.BuildingPreview{}, errors.New("test"))

			tt.fields.bot.EXPECT().
				Request(tgbotapi.NewCallback(tt.queryID, "")).Return(nil, nil).
				On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).
				Return(tgbotapi.Message{}, nil)

			h := HandlerContainer{
				tt.fields.buildingService,
				tt.fields.userService,
				tt.fields.bot,
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				metrics.NewMetrics(prometheus.NewRegistry()),
				map[string]CommandHandler{},
			}
			h.next(tt.args.ctx, tt.args.query)
		})
	}
}

func TestHandlerContainer_next_negative(t *testing.T) {
	type fields struct {
		buildingService *services.Buildings_mock
		userService     *services.Users_mock
		bot             *InternalBot_mock
	}
	type args struct {
		ctx   c.Context
		query *tgbotapi.CallbackQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		queryID string
	}{
		{
			"empty callback query",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
			},
			args{
				c.Background(),
				&tgbotapi.CallbackQuery{ID: "123"},
			},
			"123",
		},
		{
			"invalid callback data",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
			},
			args{
				c.Background(),
				&tgbotapi.CallbackQuery{
					ID:      "123",
					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{}},
				},
			},
			"123",
		},
		{
			"invalid callback text",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
			},
			args{
				c.Background(),
				&tgbotapi.CallbackQuery{
					ID: "123",
					Message: &tgbotapi.Message{
						Chat: &tgbotapi.Chat{},
						Text: "one-line text",
					},
					Data: `{"name":"next","limit":2,"offset":3}`,
				},
			},
			"123",
		},
		{
			"valid text",
			fields{
				services.NewBuildings_mock(t),
				services.NewUsers_mock(t),
				NewInternalBot_mock(t),
			},
			args{
				c.Background(),
				&tgbotapi.CallbackQuery{
					ID: "123",
					Message: &tgbotapi.Message{
						Chat: &tgbotapi.Chat{},
						Text: "address: test address   ",
					},
					Data: `{"name":"next","limit":2,"offset":3}`,
				},
			},
			"123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.bot.EXPECT().
				Request(tgbotapi.NewCallback(tt.queryID, "")).Return(nil, nil)
			h := HandlerContainer{
				tt.fields.buildingService,
				tt.fields.userService,
				tt.fields.bot,
				map[string]CommandHandler{},
				map[string]internalButtonHandler{},
				"",
				metrics.NewMetrics(prometheus.NewRegistry()),
				map[string]CommandHandler{},
			}
			h.next(tt.args.ctx, tt.args.query)
		})
	}
}
