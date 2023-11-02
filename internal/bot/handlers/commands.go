package handlers

import (
	c "context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	return HandlerContainer{service, bot, handlersPerCommand, handlersPerButton, commandsForHelp}
}

func (h HandlerContainer) GetHandler(command string) (CommandHandler, bool) {
	handler, ok := h.HandlersPerCommand[command]
	return handler, ok
}

func (h HandlerContainer) GetButtonHandler(buttonName string) (ButtonHandler, bool) {
	handler, ok := h.handlersPerButton[buttonName]
	return handler, ok
}

func (h HandlerContainer) SendMessage(chatId int64, msgText string) {
	msg := tgbotapi.NewMessage(chatId, msgText)
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("An error occured: %v", err)
	}
}

func (h HandlerContainer) start(ctx c.Context, message *tgbotapi.Message) {
	startMsg := "Hello! I'm a bot that helps you to understand Helsinki better."
	h.SendMessage(message.Chat.ID, startMsg)
}

func (h HandlerContainer) help(ctx c.Context, message *tgbotapi.Message) {
	helpMsg := fmt.Sprintf("Available commands: %s", h.commandsForHelp)
	h.SendMessage(message.Chat.ID, helpMsg)
}

func (h HandlerContainer) settings(ctx c.Context, message *tgbotapi.Message) {
	settingsMsg := "No settings yet."
	h.SendMessage(message.Chat.ID, settingsMsg)
}

func (h HandlerContainer) getAllAdresses(ctx c.Context, message *tgbotapi.Message) {
	address := message.CommandArguments()
	limit := 2
	h.returnAddresses(ctx, message.Chat.ID, address, limit, 0)
}

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
		log.Printf("can not get addresses: %v", err)
		h.SendMessage(chatID, "Internal error")
		return
	}
	items := make([]string, len(buildings)+1)
	items[0] = "Available building addresses and names:"
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
			log.Printf(
				"Can not send a response %v to the chat %v: %v",
				response,
				chatID,
				err,
			)
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID, response)
	buttonLabel := fmt.Sprintf("Next %v buildings", limit)
	button := Button{buttonLabel, "next", limit, offset + len(buildings)}
	buttonCallbackData, err := json.Marshal(button)
	if err != nil {
		log.Printf("can not create a button %v: %v", button, err)
		return
	}
	log.Println(string(buttonCallbackData))
	buttonData := tgbotapi.NewInlineKeyboardButtonData(
		button.label,
		string(buttonCallbackData),
	)
	moreAddressesMenuMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonData),
	)
	msg.ReplyMarkup = moreAddressesMenuMarkup
	if _, err := h.bot.Send(msg); err != nil {
		log.Printf("An error occured: %v", err)
	}
}

func (h HandlerContainer) getBuilding(ctx c.Context, message *tgbotapi.Message) {
	address := message.CommandArguments()
	if address == "" {
		h.SendMessage(message.Chat.ID, "Please add an address to this command.")
		return
	}
	buildings, err := h.buildingService.GetBuildingsByAddress(ctx, address)
	if err != nil {
		log.Printf("can not get building by address '%s': %v", address, err)
		h.SendMessage(message.Chat.ID, "Internal error.")
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
			log.Printf("can not serialize a building '%s': %v", address, err)
			items[i] = "A building error."
			continue
		}
		items[i] = serializedItem
	}
	response := strings.Join(items, "\n\n")
	h.SendMessage(message.Chat.ID, response)
}

func (h HandlerContainer) next(ctx c.Context, query *tgbotapi.CallbackQuery) {
	var button Button
	if err := json.Unmarshal([]byte(query.Data), &button); err != nil {
		log.Printf("unexpected callback data %v: %v", query, err)
		return
	}
	h.returnAddresses(ctx, query.Message.Chat.ID, "", button.Limit, button.Offset)
}
