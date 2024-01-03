package handlers

import (
	c "context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/logger"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/metrics"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	unexpectedTextTmpl = "an unexpected message text for a button " +
		"'next': '%s': message_id: '%v': chat id: '%v'"
	defaultLimit = 10
)

func NewCommandContainer(
	bot InternalBot,
	service services.BuildingService,
	userService services.UserService,
	metricsContainer *metrics.Metrics,
) HandlerContainer {
	handlersPerButton := map[string]internalButtonHandler{
		"next":     HandlerContainer.next,
		"language": HandlerContainer.language,
	}
	availableCommands := []string{}
	for command := range handlersPerCommand {
		availableCommands = append(availableCommands, "/"+command)
	}
	slices.Sort(availableCommands)
	commandsForHelp := strings.Join(availableCommands, ", ")
	allHandlers := make(map[string]CommandHandler)
	for command, handler := range handlersPerCommand {
		allHandlers[command] = handler
	}
	allHandlers["nearestAddresses"] = CommandHandler{
		HandlerContainer.getNearestAddresses,
		"get nearest addresses",
	}
	return HandlerContainer{
		service,
		userService,
		bot,
		handlersPerCommand,
		handlersPerButton,
		commandsForHelp,
		metricsContainer,
		allHandlers,
	}
}

func (h HandlerContainer) GetCommandHandler(command string) (func(c.Context, *tgbotapi.Message) error, bool) {
	handler, ok := h.allHandlers[command]
	if !ok {
		return nil, false
	}
	metricWrapper := func(ctx c.Context, message *tgbotapi.Message) error {
		now := time.Now()
		err := handler.Function(h, ctx, message)
		h.metrics.CommandDuration.With(
			prometheus.Labels{"command_name": command},
		).Observe(time.Since(now).Seconds())
		if err != nil {
			h.metrics.HandlerErrors.With(
				prometheus.Labels{"handler_name": command},
			).Inc()
		}
		return err
	}
	return metricWrapper, ok
}

func (h HandlerContainer) GetButtonHandler(buttonName string) (ButtonHandler, bool) {
	handler, ok := h.handlersPerButton[buttonName]
	if !ok {
		return nil, false
	}
	metricWrapper := func(ctx c.Context, query *tgbotapi.CallbackQuery) error {
		now := time.Now()
		err := handler(h, ctx, query)
		h.metrics.ButtonDuration.With(
			prometheus.Labels{"button_name": buttonName},
		).Observe(time.Since(now).Seconds())
		if err != nil {
			h.metrics.HandlerErrors.With(
				prometheus.Labels{"handler_name": buttonName},
			).Inc()
		}
		return err
	}
	return metricWrapper, ok
}

func (h HandlerContainer) SendMessage(ctx c.Context, chatId int64, msgText string, parseMode string) error {
	msg := tgbotapi.NewMessage(chatId, msgText)
	msg.ParseMode = parseMode
	_, err := h.bot.Send(msg)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send a message to %v: %v", chatId, msgText),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}

func (h HandlerContainer) start(ctx c.Context, message *tgbotapi.Message) error {
	if message.Chat == nil {
		return ErrNoChat
	}
	chatID := message.Chat.ID
	startMsg := "Hello! I'm a bot that provides information about Helsinki buildings."
	msg := tgbotapi.NewMessage(
		chatID,
		startMsg+"\n\n"+helpMessage,
	)
	locationButton := tgbotapi.NewKeyboardButtonLocation(
		"Share my location and get the nearest buildings",
	)
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

func (h HandlerContainer) help(ctx c.Context, message *tgbotapi.Message) error {
	return h.SendMessage(ctx, message.Chat.ID, helpMessage, "")
}

func (h HandlerContainer) settings(ctx c.Context, message *tgbotapi.Message) error {
	if message.Chat == nil {
		return ErrNoChat
	}
	chatID := message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Choose a preferable language:")
	buttons := []tgbotapi.InlineKeyboardButton{}
	languageButtons := []LanguageButton{
		{Button{"Finnish", "language"}, "fi"},
		{Button{"English", "language"}, "en"},
		{Button{"Russian", "language"}, "ru"},
	}
	for _, button := range languageButtons {
		buttonCallbackData, err := json.Marshal(button)
		if err != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("can not create a button %v", button),
				slog.Any(logger.ErrorKey, err),
			)
			sendErr := h.SendMessage(ctx, chatID, "Internal error", "")
			return errors.Join(sendErr, err)
		}
		buttons = append(
			buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				button.label,
				string(buttonCallbackData),
			),
		)
	}
	settingsMenuMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
	msg.ReplyMarkup = settingsMenuMarkup
	_, err := h.bot.Send(msg)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send a settings keyboard to: %v", chatID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}

func (h HandlerContainer) getAllAdresses(ctx c.Context, message *tgbotapi.Message) error {
	address := message.CommandArguments()
	return h.returnAddresses(ctx, message.Chat.ID, address, defaultLimit, 0)
}

func (h HandlerContainer) returnAddresses(
	ctx c.Context,
	chatID int64,
	address string,
	limit,
	offset int,
) error {
	buildings, err := h.buildingService.GetBuildingPreviews(
		ctx,
		address,
		limit,
		offset,
	)
	if err != nil {
		sendErr := h.SendMessage(ctx, chatID, "Internal error", "")
		return errors.Join(sendErr, err)
	}
	items := make([]string, len(buildings)+1)
	items[0] = fmt.Sprintf(headerTemplate, address)
	for i, building := range buildings {
		items[i+1] = fmt.Sprintf(
			lineTemplate,
			offset+i+1,
			building.Address,
			building.Name,
		)
	}
	response := strings.Join(items, "\n")
	if len(buildings) < limit {
		response += "\nEnd"
		return h.SendMessage(ctx, chatID, response, "")
	}

	msg := tgbotapi.NewMessage(chatID, response)
	button := NextButton{
		Button{fmt.Sprintf("Next %v buildings", limit), "next"},
		limit,
		offset + len(buildings),
	}
	buttonCallbackData, err := json.Marshal(button)
	if err != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf("can not create a button %v", button),
			slog.Any(logger.ErrorKey, err),
		)
		return err
	}
	buttonData := tgbotapi.NewInlineKeyboardButtonData(
		button.label,
		string(buttonCallbackData),
	)
	moreAddressesMenuMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonData),
	)
	msg.ReplyMarkup = moreAddressesMenuMarkup
	_, err = h.bot.Send(msg)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send an inline keyboard to: %v", chatID),
			slog.Any(logger.ErrorKey, err),
		)
	}
	return err
}

func (h HandlerContainer) getBuilding(ctx c.Context, message *tgbotapi.Message) error {
	address := message.CommandArguments()
	if address == "" {
		return h.SendMessage(
			ctx,
			message.Chat.ID,
			"Please add an address to this command.",
			"",
		)
	}
	buildings, err := h.buildingService.GetBuildingsByAddress(ctx, address)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not get building by address '%s'", address),
			slog.Any(logger.ErrorKey, err),
		)
		sendErr := h.SendMessage(ctx, message.Chat.ID, "Internal error.", "")
		return errors.Join(sendErr, err)
	}
	if len(buildings) == 0 {
		response := "Unfortunately, I don't know this address."
		return h.SendMessage(ctx, message.Chat.ID, response, tgbotapi.ModeHTML)
	}
	userLanguage := services.English
	if user := message.From; user != nil {
		switch user.LanguageCode {
		case "fi":
			userLanguage = services.Finnish
		case "ru":
			userLanguage = services.Russian
		}
		preferredLanguage, err := h.userService.GetPreferredLanguage(
			ctx,
			user.ID,
		)
		if err == nil && preferredLanguage != nil {
			userLanguage = *preferredLanguage
		}
	}
	items := make([]string, len(buildings))
	for i, building := range buildings {
		serializedItem, err := SerializeIntoMessage(building, userLanguage)
		if err != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("can not serialize a building '%s'", address),
				slog.Any(logger.ErrorKey, err),
			)
			items[i] = "A building error."
			continue
		}
		items[i] = serializedItem
	}
	response := strings.Join(items, "\n\n")
	return h.SendMessage(ctx, message.Chat.ID, response, tgbotapi.ModeHTML)
}
