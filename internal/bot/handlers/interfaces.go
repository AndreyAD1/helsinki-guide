package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"


type InternalBot interface {
	Request(tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}