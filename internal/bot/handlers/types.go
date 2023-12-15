package handlers

import (
	c "context"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/middlewares"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	Function    func(HandlerContainer, c.Context, *tgbotapi.Message) error
	Description string
}
type internalButtonHandler func(HandlerContainer, c.Context, *tgbotapi.CallbackQuery) error
type ButtonHandler func(c.Context, *tgbotapi.CallbackQuery) error
type HandlerContainer struct {
	buildingService    services.Buildings
	bot                InternalBot
	HandlersPerCommand map[string]CommandHandler
	handlersPerButton  map[string]internalButtonHandler
	commandsForHelp    string
	metrics            *metrics.Metrics
	allHandlers        map[string]CommandHandler
}
type Button struct {
	label  string
	Name   string `json:"name"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}
type BotWithMetrics struct {
	clientName string
	*tgbotapi.BotAPI
	m *metrics.Metrics
}

func NewBotWithMetrics(bot *tgbotapi.BotAPI, m *metrics.Metrics) *BotWithMetrics {
	return &BotWithMetrics{"Telegram", bot, m}
}

func (b *BotWithMetrics) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	result, err := middlewares.Duration(
		func() (interface{}, error) { return b.BotAPI.Send(c) },
		b.m,
		b.clientName,
		"Send",
	)
	message := result.(tgbotapi.Message)
	return message, err
}

func (b *BotWithMetrics) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	result, err := middlewares.Duration(
		func() (interface{}, error) { return b.BotAPI.Request(c) },
		b.m,
		b.clientName,
		"Request",
	)
	response := result.(*tgbotapi.APIResponse)
	return response, err
}
