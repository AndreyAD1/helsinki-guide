package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var HandlersPerCommand = map[string]func(*tgbotapi.BotAPI, *tgbotapi.Message) {
	"start": start,
	"help": help,
	"settings": settings,
}

func start(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	startMsg := "Hello! I'm a bot that helps you to understand Helsinki better."
	msg := tgbotapi.NewMessage(message.Chat.ID, startMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func help(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	helpMsg := "No help yet."
	msg := tgbotapi.NewMessage(message.Chat.ID, helpMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func settings(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	settingMsg := "No settings yet."
	msg := tgbotapi.NewMessage(message.Chat.ID, settingMsg)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}