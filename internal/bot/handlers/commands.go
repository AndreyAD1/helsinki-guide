package handlers

import (
	c "context"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	unexpectedTextTmpl = "an unexpected message text for a button " +
		"'next': %s: message_id: %v: chat id: %v"
	defaultLimit = 10
)

func NewCommandContainer(
	bot *tgbotapi.BotAPI,
	service services.BuildingService,
) HandlerContainer {
	handlersPerCommand := map[string]CommandHandler{
		"start":     {HandlerContainer.start, "Start the bot"},
		"help":      {HandlerContainer.help, "Get help"},
		"settings":  {HandlerContainer.settings, "Configure settings"},
		"addresses": {HandlerContainer.getAllAdresses, "Get all available addresses"},
		"building":  {HandlerContainer.getBuilding, "Get building by address"},
	}
	handlersPerButton := map[string]ButtonHandler{
		"next": HandlerContainer.next,
	}
	availableCommands := []string{}
	for command := range handlersPerCommand {
		availableCommands = append(availableCommands, "/"+command)
	}
	slices.Sort(availableCommands)
	commandsForHelp := strings.Join(availableCommands, ", ")
	return HandlerContainer{
		service,
		bot,
		handlersPerCommand,
		handlersPerButton,
		commandsForHelp,
	}
}

func (h HandlerContainer) GetHandler(command string) (CommandHandler, bool) {
	handler, ok := h.HandlersPerCommand[command]
	return handler, ok
}

func (h HandlerContainer) GetButtonHandler(buttonName string) (ButtonHandler, bool) {
	handler, ok := h.handlersPerButton[buttonName]
	return handler, ok
}

func (h HandlerContainer) SendMessage(ctx c.Context, chatId int64, msgText string) {
	msg := tgbotapi.NewMessage(chatId, msgText)
	if _, err := h.bot.Send(msg); err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can no send a message to %v: %v", chatId, msgText),
			slog.Any(logger.ErrorKey, err),
		)
	}
}

func (h HandlerContainer) start(ctx c.Context, message *tgbotapi.Message) {
	startMsg := "Hello! I'm a bot that helps you to understand Helsinki better."
	h.SendMessage(ctx, message.Chat.ID, startMsg)
}

func (h HandlerContainer) help(ctx c.Context, message *tgbotapi.Message) {
	helpMsg := fmt.Sprintf("Available commands: %s", h.commandsForHelp)
	h.SendMessage(ctx, message.Chat.ID, helpMsg)
}

func (h HandlerContainer) settings(ctx c.Context, message *tgbotapi.Message) {
	settingsMsg := "No settings yet."
	h.SendMessage(ctx, message.Chat.ID, settingsMsg)
}

func (h HandlerContainer) getAllAdresses(ctx c.Context, message *tgbotapi.Message) {
	address := message.CommandArguments()
	h.returnAddresses(ctx, message.Chat.ID, address, defaultLimit, 0)
}

var prefixTemplate = "Search address: %s\nAvailable building addresses and names:"

func (h HandlerContainer) returnAddresses(
	ctx c.Context,
	chatID int64,
	address string,
	limit,
	offset int,
) {
	buildings, err := h.buildingService.GetBuildingPreviews(
		ctx,
		address,
		limit,
		offset,
	)
	if err != nil {
		h.SendMessage(ctx, chatID, "Internal error")
		return
	}
	items := make([]string, len(buildings)+1)
	items[0] = fmt.Sprintf(prefixTemplate, address)
	template := "%v. %s - %s"
	for i, building := range buildings {
		items[i+1] = fmt.Sprintf(
			template,
			offset+i+1,
			building.Address,
			building.Name,
		)
	}
	response := strings.Join(items, "\n")
	if len(buildings) < limit {
		response += "\nEnd"
		msg := tgbotapi.NewMessage(chatID, response)
		if _, err := h.bot.Send(msg); err != nil {
			slog.WarnContext(
				ctx,
				fmt.Sprintf(
					"Can not send a response %v to the chat %v",
					response,
					chatID,
				),
				slog.Any(logger.ErrorKey, err),
			)
		}
		return
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
		return
	}
	buttonData := tgbotapi.NewInlineKeyboardButtonData(
		button.label,
		string(buttonCallbackData),
	)
	moreAddressesMenuMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonData),
	)
	msg.ReplyMarkup = moreAddressesMenuMarkup
	if _, err := h.bot.Send(msg); err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not send an inline keyboard to: %v", chatID),
			slog.Any(logger.ErrorKey, err),
		)
	}
}

func (h HandlerContainer) getBuilding(ctx c.Context, message *tgbotapi.Message) {
	address := message.CommandArguments()
	if address == "" {
		h.SendMessage(ctx, message.Chat.ID, "Please add an address to this command.")
		return
	}
	buildings, err := h.buildingService.GetBuildingsByAddress(ctx, address)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("can not get building by address '%s'", address),
			slog.Any(logger.ErrorKey, err),
		)
		h.SendMessage(ctx, message.Chat.ID, "Internal error.")
		return
	}
	userLanguage := "en"
	if user := message.From; user != nil {
		userLanguage = user.LanguageCode
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
	h.SendMessage(ctx, message.Chat.ID, response)
}

func (h HandlerContainer) next(ctx c.Context, query *tgbotapi.CallbackQuery) {
	msgID, chatID := query.Message.MessageID, query.Message.Chat.ID
	var button Button
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		logMsg := fmt.Sprintf(
			"unexpected callback data %v from a message %v and the chat %v",
			query.Data,
			msgID,
			chatID,
		)
		slog.ErrorContext(ctx, logMsg, slog.Any(logger.ErrorKey, err))
		return
	}
	// I need to extract an address from a message text
	//  instead of using query data because the Telegram API specifies that
	//  query data should be less than 64 bytes.
	firstRow, _, found := strings.Cut(query.Message.Text, "\n")
	logMsg := fmt.Sprintf(unexpectedTextTmpl, query.Message.Text, msgID, chatID)
	if !found {
		slog.ErrorContext(ctx, logMsg)
		return
	}
	_, address, found := strings.Cut(firstRow, ":")
	if !found {
		slog.ErrorContext(ctx, logMsg)
		return
	}
	address = strings.TrimSpace(address)
	h.returnAddresses(ctx, chatID, address, button.Limit, button.Offset)
	// Telegram asks a bot server to explicitly answer every callback call
	callbackAnswer := tgbotapi.NewCallback(query.ID, "")
	_, err := h.bot.Request(callbackAnswer)
	if err != nil {
		slog.WarnContext(
			ctx,
			fmt.Sprintf("could not answer to a callback %v", query.ID),
			slog.Any(logger.ErrorKey, err),
		)
	}
}
