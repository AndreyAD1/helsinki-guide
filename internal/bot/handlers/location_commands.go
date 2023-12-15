package handlers

import (
	c "context"
	"fmt"
	"log/slog"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h HandlerContainer) location(ctx c.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	msg := tgbotapi.NewMessage(
		chatID,
		"Please share your location to get nearest buildings.",
	)
	locationButton := tgbotapi.NewKeyboardButtonLocation("Share my location")
	keyboardMarkup := tgbotapi.NewOneTimeReplyKeyboard(
		[]tgbotapi.KeyboardButton{locationButton},
	)
	msg.ReplyMarkup = keyboardMarkup
	_, err := h.bot.Send(msg)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send an inline keyboard to: %v", chatID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}

func (h HandlerContainer) getNearestAddresses(ctx c.Context, message *tgbotapi.Message) error {
	if message.Chat == nil {
		return ErrNoChat
	}
	chatID := message.Chat.ID

	// buildings, err := h.buildingService.GetNearestBuildingPreviews(
	// 	ctx,
	// 	latitude,
	// 	longitude,
	// )
	// responseText := getAddressResponse(buildings)
	// h.SendMessage(ctx, message.Chat.ID, responseText)
	return h.SendMessage(ctx, chatID, "dummy nearest buildings")
}
