package handlers

import (
	c "context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const DEFAULT_DISTANCE = 200

func (h HandlerContainer) getNearestAddresses(ctx c.Context, message *tgbotapi.Message) error {
	if message.Chat == nil {
		return ErrNoChat
	}
	location := message.Location
	if location == nil {
		return ErrNoLocation
	}
	buildings, err := h.buildingService.GetNearestBuildingPreviews(
		ctx,
		DEFAULT_DISTANCE,
		location.Latitude,
		location.Longitude,
		defaultLimit,
		0,
	)
	if err != nil {
		sendErr := h.SendMessage(ctx, message.Chat.ID, "Internal error", "")
		return errors.Join(sendErr, err)
	}
	title := fmt.Sprintf("Nearest buildings in %v meters:", DEFAULT_DISTANCE)
	msg := tgbotapi.NewMessage(message.Chat.ID, title)
	keyboardRows, err := getBuildingButtonRows(ctx, buildings)
	if err != nil {
		return err
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
	_, err = h.bot.Send(msg)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send nearest addresses to: %v", message.Chat.ID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}
