package handlers

import (
	c "context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DEFAULT_DISTANCE                  = 100
	noNearestBuildingsEnglishTemplate = "Unfortunately, I'm not aware of any buildings located within %v metres of your location."
	noNearestBuildingsFinnishTemplate = "Valitettavasti minulla ei ole tietoa yhdestäkään rakennuksesta %v metrin säteellä sinusta."
	noNearestBuildingsRussianTemplate = "К сожалению, у меня нет информации ни об одном здании в радиусе %v метров от вас."
	nearestBuildingsEnglishTemplate   = "Here are the closest buildings I'm aware of, situated within %v meters of your location:"
	nearestBuildingsFinnishTemplate   = "Luettelo rakennuksista, jotka sijaitsevat %v metrin säteellä sinusta:"
	nearestBuildingsRussianTemplate   = "Список зданий, расположенных в радиусе %v метров от вас:"
)

func (h HandlerContainer) getNearestAddresses(ctx c.Context, message *tgbotapi.Message) error {
	if message.Chat == nil {
		return ErrNoChat
	}
	location := message.Location
	if location == nil {
		return ErrNoLocation
	}
	buildings, err := h.buildingService.GetNearestBuildings(
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
	language := h.getPreferredLanguage(ctx, message.From)
	if len(buildings) == 0 {
		responseTemplate := noNearestBuildingsEnglishTemplate
		switch language {
		case services.Finnish:
			responseTemplate = noNearestBuildingsFinnishTemplate
		case services.Russian:
			responseTemplate = noNearestBuildingsRussianTemplate
		}
		msg := fmt.Sprintf(responseTemplate, DEFAULT_DISTANCE)
		return h.SendMessage(ctx, message.Chat.ID, msg, "")
	}
	titleTemplate := nearestBuildingsEnglishTemplate
	switch language {
	case services.Finnish:
		titleTemplate = nearestBuildingsFinnishTemplate
	case services.Russian:
		titleTemplate = nearestBuildingsRussianTemplate
	}
	title := fmt.Sprintf(titleTemplate, DEFAULT_DISTANCE)
	msg := tgbotapi.NewMessage(message.Chat.ID, title)
	keyboardRows, err := getBuildingButtonRows(ctx, language, buildings)
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
