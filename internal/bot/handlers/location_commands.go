package handlers

import (
	c "context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
