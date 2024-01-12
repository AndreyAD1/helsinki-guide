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
	"unicode/utf8"

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
		NEXT_BUTTON:     HandlerContainer.next,
		LANGUAGE_BUTTON: HandlerContainer.language,
		BUILDING_BUTTON: HandlerContainer.building,
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
		"get the nearest addresses",
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

func (h HandlerContainer) ProcessCommonMessage(ctx c.Context, message *tgbotapi.Message) error {
	filteredText := strings.Trim(message.Text, " ")
	if filteredText == "" {
		return h.SendMessage(
			ctx,
			message.Chat.ID,
			"Please enter any address.",
			tgbotapi.ModeHTML,
		)
	}
	if utf8.RuneCountInString(filteredText) > MAX_MESSAGE_LENGTH {
		return h.SendMessage(
			ctx,
			message.Chat.ID,
			fmt.Sprintf(
				"Please enter an address with less than %v characters.",
				MAX_MESSAGE_LENGTH,
			),
			tgbotapi.ModeHTML,
		)
	}
	return h.returnAddresses(ctx, message.Chat.ID, filteredText, defaultLimit, 0)
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
		{Button{"Finnish", LANGUAGE_BUTTON}, "fi"},
		{Button{"English", LANGUAGE_BUTTON}, "en"},
		{Button{"Russian", LANGUAGE_BUTTON}, "ru"},
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
	return h.returnAddresses(ctx, message.Chat.ID, "", defaultLimit, 0)
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
	title := fmt.Sprintf(headerTemplate, address)
	msg := tgbotapi.NewMessage(chatID, title)
	if len(buildings) == 0 {
		msg.Text += "\nNo buildings are found."
		_, err = h.bot.Send(msg)
		return err
	}
	keyboardRows, err := getBuildingButtonRows(ctx, buildings)
	if err != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf("can not create building rows: '%v'", address),
			slog.Any(logger.ErrorKey, err),
		)
		return err
	}
	if len(buildings) < limit {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
		_, err = h.bot.Send(msg)
		if err != nil {
			slog.WarnContext(ctx, err.Error())
		}
		return err
	}
	button := NextButton{
		Button{fmt.Sprintf("Next %v buildings", limit), NEXT_BUTTON},
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
	keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(buttonData))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
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
