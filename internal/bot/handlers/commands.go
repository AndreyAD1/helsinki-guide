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
	metricsContainer *metrics.Metrics,
) HandlerContainer {
	handlersPerButton := map[string]internalButtonHandler{
		"next": HandlerContainer.next,
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

func (h HandlerContainer) SendMessage(ctx c.Context, chatId int64, msgText string) error {
	msg := tgbotapi.NewMessage(chatId, msgText)
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
		startMsg,
	)
	locationButton := tgbotapi.NewKeyboardButtonLocation("Share my location and get nearest buildings")
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
	helpMsg := fmt.Sprintf("Available commands: %s", h.commandsForHelp)
	return h.SendMessage(ctx, message.Chat.ID, helpMsg)
}

func (h HandlerContainer) settings(ctx c.Context, message *tgbotapi.Message) error {
	settingsMsg := "No settings yet."
	return h.SendMessage(ctx, message.Chat.ID, settingsMsg)
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
		sendErr := h.SendMessage(ctx, chatID, "Internal error")
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
		return h.SendMessage(ctx, chatID, response)
	}

	msg := tgbotapi.NewMessage(chatID, response)
	buttonLabel := fmt.Sprintf("Next %v buildings", limit)
	button := Button{buttonLabel, "next", limit, offset + len(buildings)}
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
		)
	}
	buildings, err := h.buildingService.GetBuildingsByAddress(ctx, address)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not get building by address '%s'", address),
			slog.Any(logger.ErrorKey, err),
		)
		sendErr := h.SendMessage(ctx, message.Chat.ID, "Internal error.")
		return errors.Join(sendErr, err)
	}
	userLanguage := English
	if user := message.From; user != nil {
		switch user.LanguageCode {
		case "fi":
			userLanguage = Finnish
		case "ru":
			userLanguage = Russian
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
	response := "Unfortunately, I don't know this address."
	if len(items) > 0 {
		response = strings.Join(items, "\n\n")
	}
	return h.SendMessage(ctx, message.Chat.ID, response)
}

func (h HandlerContainer) next(ctx c.Context, query *tgbotapi.CallbackQuery) error {
	// Telegram asks a bot server to explicitly answer every callback call
	defer func() {
		callbackAnswer := tgbotapi.NewCallback(query.ID, "")
		_, err := h.bot.Request(callbackAnswer)
		if err != nil {
			slog.WarnContext(
				ctx,
				fmt.Sprintf("could not answer to a callback %v", query.ID),
				slog.Any(logger.ErrorKey, err),
			)
		}
	}()

	message := query.Message
	if message == nil {
		errMsg := fmt.Sprintf("a callback has no message %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "a callback has no message"},
		).Inc()
		return nil
	}
	msgID := query.Message.MessageID
	chat := query.Message.Chat
	if chat == nil {
		errMsg := fmt.Sprintf("a callback has no chat %v", query.ID)
		slog.WarnContext(ctx, errMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "a callback has no chat"},
		).Inc()
		return nil
	}
	var button Button
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		logMsg := fmt.Sprintf(
			"unexpected callback data %v from a message %v and the chat %v",
			query.Data,
			msgID,
			chat.ID,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "unexpected callback data"},
		).Inc()
		return nil
	}
	// I need to extract an address from a message text
	//  instead of using query data because the Telegram API specifies that
	//  query data should be less than 64 bytes.
	firstRow, _, found := strings.Cut(query.Message.Text, "\n")
	logMsg := fmt.Sprintf(unexpectedTextTmpl, query.Message.Text, msgID, chat.ID)
	if !found {
		slog.ErrorContext(ctx, logMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "unexpected callback message"},
		).Inc()
		return nil
	}
	_, address, found := strings.Cut(firstRow, ":")
	if !found {
		slog.ErrorContext(ctx, logMsg)
		h.metrics.UnexpectedNextCallback.With(
			prometheus.Labels{"error": "unexpected callback message"},
		).Inc()
		return nil
	}
	address = strings.TrimSpace(address)
	return h.returnAddresses(ctx, chat.ID, address, button.Limit, button.Offset)
}
