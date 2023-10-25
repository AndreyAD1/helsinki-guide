package bot

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot *tgbotapi.BotAPI
	startMessage = "Hello! I'm a bot that helps you to understand Helsinki better."
)

func RunBot(botToken string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	// Set this to true to log all interactions with telegram servers
	bot.Debug = false
	
	if err := setBotCommands(); err != nil {
		log.Fatalf("can not set commands: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func setBotCommands() error {
	commands := []tgbotapi.BotCommand{}
	for commandName, handler := range handlers.HandlersPerCommand {
		command := tgbotapi.BotCommand{
			Command: commandName, 
			Description: handler.Description,
		}
		commands = append(commands, command)
	}
	setCommandsConfig := tgbotapi.NewSetMyCommands(commands...)
	result, err := bot.Request(setCommandsConfig)
	if err != nil {
		return fmt.Errorf("can not make a request to set commands: %w", err)
	}
	if !result.Ok {
		return fmt.Errorf(result.Description)
	}

	return nil
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		handleMessage(update.Message)
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From

	if user == nil {
		return
	}
	handler, ok := handlers.HandlersPerCommand[message.Command()]
	if !ok {
		answer := fmt.Sprintf("I don't understand this message: %s", message.Text)
		msg := tgbotapi.NewMessage(message.Chat.ID, answer)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("An error occured: %s", err.Error())
		}
		return
	}
	handler.Function(bot, message)
}

func handleButton(query *tgbotapi.CallbackQuery) {
	log.Println("a callback is not supported")
}
