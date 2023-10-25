package bot

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AndreyAD1/helsinki-guide/internal/bot/handlers"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/services"
	"github.com/AndreyAD1/helsinki-guide/internal/infrastructure/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Server struct {
	bot      *tgbotapi.BotAPI
	handlers handlers.HandlerContainer
}

func NewServer(botToken string) (*Server, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = false
	db, err := storage.NewStorage()
	if err != nil {
		return nil, err
	}
	buildingService := services.NewService(db)
	handlerContainer := handlers.NewHandler(buildingService)
	return &Server{bot, handlerContainer}, nil
}

func (s *Server) RunBot() {
	if err := s.setBotCommands(); err != nil {
		log.Fatalf("can not set commands: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	updates := s.bot.GetUpdatesChan(u)

	go s.receiveUpdates(ctx, updates)

	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (s *Server) setBotCommands() error {
	commands := []tgbotapi.BotCommand{}
	for commandName, handler := range s.handlers.HandlersPerCommand {
		command := tgbotapi.BotCommand{
			Command:     commandName,
			Description: handler.Description,
		}
		commands = append(commands, command)
	}
	setCommandsConfig := tgbotapi.NewSetMyCommands(commands...)
	result, err := s.bot.Request(setCommandsConfig)
	if err != nil {
		return fmt.Errorf("can not make a request to set commands: %w", err)
	}
	if !result.Ok {
		return fmt.Errorf(result.Description)
	}

	return nil
}

func (s *Server) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			s.handleUpdate(update)
		}
	}
}

func (s *Server) handleUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		s.handleMessage(update.Message)
	case update.CallbackQuery != nil:
		s.handleButton(update.CallbackQuery)
	}
}

func (s *Server) handleMessage(message *tgbotapi.Message) {
	user := message.From

	if user == nil {
		return
	}
	handler, ok := s.handlers.GetHandler(message.Command())
	if !ok {
		answer := fmt.Sprintf("I don't understand this message: %s", message.Text)
		msg := tgbotapi.NewMessage(message.Chat.ID, answer)
		if _, err := s.bot.Send(msg); err != nil {
			log.Printf("An error occured: %s", err.Error())
		}
		return
	}
	handler.Function(s.handlers, s.bot, message)
}

func (s *Server) handleButton(query *tgbotapi.CallbackQuery) {
	log.Println("a callback is not supported")
}
