package handlers

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Function    func(HandlerContainer, *tgbotapi.BotAPI, *tgbotapi.Message)
	Description string
}

type HandlerContainer struct {
	addressService     services.AddressService
	HandlersPerCommand map[string]Handler
	commandsForHelp    string
}

func NewHandler(service services.AddressService) HandlerContainer {
	handlersPerCommand := map[string]Handler{
		"start":    {HandlerContainer.start, "Start the bot"},
		"help":     {HandlerContainer.help, "Get help"},
		"settings": {HandlerContainer.settings, "Configure settings"},
	}
	availableCommands := []string{}
	for command := range handlersPerCommand {
		availableCommands = append(availableCommands, "/"+command)
	}
	slices.Sort(availableCommands)
	commandsForHelp := strings.Join(availableCommands, ", ")
	return HandlerContainer{service, handlersPerCommand, commandsForHelp}
}

func (h HandlerContainer) GetHandler(command string) (Handler, bool) {
	handler, ok := h.HandlersPerCommand[command]
	return handler, ok
}

func (h HandlerContainer) start(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	startMsg := "Hello! I'm a bot that helps you to understand Helsinki better."
	msg := tgbotapi.NewMessage(message.Chat.ID, startMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func (h HandlerContainer) help(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	helpMsg := fmt.Sprintf("Available commands: %s", h.commandsForHelp)
	msg := tgbotapi.NewMessage(message.Chat.ID, helpMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func (h HandlerContainer) settings(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	settingMsg := "No settings yet."
	msg := tgbotapi.NewMessage(message.Chat.ID, settingMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func (h HandlerContainer) getAllAdresses(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {}
