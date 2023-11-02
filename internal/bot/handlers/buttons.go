package handlers

import (
	c "context"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewButtonContainer(
	bot *tgbotapi.BotAPI,
	service services.BuildingService,
) ButtonHandlerContainer {
	handlersPerButton := map[string]ButtonHandler{
		"next": ButtonHandlerContainer.next,
	}
	return ButtonHandlerContainer{service, bot, handlersPerButton}
}

func (h ButtonHandlerContainer) GetHandler(button string) (ButtonHandler, bool) {
	handler, ok := h.handlersPerButton[button]
	return handler, ok
}

func (b ButtonHandlerContainer) next(ctx c.Context, query *tgbotapi.CallbackQuery) {
	return
}
