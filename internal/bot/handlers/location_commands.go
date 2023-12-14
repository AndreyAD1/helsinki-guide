package handlers

import (
	c "context"
	"fmt"
	"log/slog"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h HandlerContainer) location(ctx c.Context, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Please share your location.")
	locationButton := tgbotapi.NewKeyboardButtonLocation("Share my location")
	keyboardMarkup := tgbotapi.NewOneTimeReplyKeyboard(
		[]tgbotapi.KeyboardButton{locationButton},
	)
	msg.ReplyMarkup = keyboardMarkup
	if _, err := h.bot.Send(msg); err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send an inline keyboard to: %v", chatID),
			slog.Any(logger.ErrorKey, err),
		)
	}
}