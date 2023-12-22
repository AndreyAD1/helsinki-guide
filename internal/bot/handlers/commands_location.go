package handlers

import (
	c "context"
	"errors"
	"fmt"
	"strings"

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
		sendErr := h.SendMessage(ctx, message.Chat.ID, "Internal error")
		return errors.Join(sendErr, err)
	}
	items := make([]string, len(buildings)+1)
	items[0] = "Nearest buildings:"
	for i, building := range buildings {
		items[i+1] = fmt.Sprintf(
			lineTemplate,
			i+1,
			building.Address,
			building.Name,
		)
	}
	response := strings.Join(items, "\n")
	return h.SendMessage(ctx, message.Chat.ID, response)
}
