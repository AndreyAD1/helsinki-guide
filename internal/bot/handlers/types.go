package handlers

import (
	c "context"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	Function    func(HandlerContainer, c.Context, *tgbotapi.Message)
	Description string
}

type internalButtonHandler func(HandlerContainer, c.Context, *tgbotapi.CallbackQuery)
type ButtonHandler func(c.Context, *tgbotapi.CallbackQuery)

type HandlerContainer struct {
	buildingService    services.BuildingService
	bot                *tgbotapi.BotAPI
	HandlersPerCommand map[string]CommandHandler
	handlersPerButton  map[string]internalButtonHandler
	commandsForHelp    string
	metrics            *metrics.Metrics
}

type Button struct {
	label  string
	Name   string `json:"name"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}
