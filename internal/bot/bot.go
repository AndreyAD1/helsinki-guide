package bot

import (
	"bufio"
	"context"
	"log"
	"os"

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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

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
		break

	case update.CallbackQuery != nil:
		log.Println("a callback is not supported")
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}
	var answer string
	switch message.Command() {
	case "start":
		answer = startMessage
	case "help":
		answer = "no help yet"
	case "settings":
		answer = "no settings yet"
	default:
		answer = "This is not a command I can understand"
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, answer)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
	log.Printf("%s wrote %s", user.FirstName, text)
}
