package handlers

import (
	c "context"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Function    func(HandlerContainer, c.Context, *tgbotapi.Message)
	Description string
}

type HandlerContainer struct {
	buildingService    services.BuildingService
	bot                *tgbotapi.BotAPI
	HandlersPerCommand map[string]Handler
	commandsForHelp    string
}

type CallBackQuery struct {
	Name   string `json:"name"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
}
