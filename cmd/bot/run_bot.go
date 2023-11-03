package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/AndreyAD1/helsinki-guide/internal/bot"
	"github.com/AndreyAD1/helsinki-guide/internal/bot/configuration"
	"github.com/caarlos0/env/v9"
	"github.com/spf13/cobra"
)

var (
	botToken string
	BotCmd   = &cobra.Command{
		Use:   "bot",
		Short: "Run a Telegram bot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	panicCounter   int
	panicThreshold = 10
	logLevel       = new(slog.LevelVar)
)

func init() {
	BotCmd.Flags().StringVarP(
		&botToken,
		"token",
		"t",
		"",
		"a token of Telegram bot",
	)
}

func run() error {
	ctx := context.Background()
	handlerOptions := slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &handlerOptions))
	slog.SetDefault(logger)

	config := configuration.StartupConfig{}
	err := env.Parse(&config)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"a configuration error",
			slog.Any("error", err),
		)
		return fmt.Errorf("a configuration error: %w", err)
	}
	if config.Debug {
		logLevel.Set(slog.LevelDebug)
	}
	if botToken != "" {
		config.BotAPIToken = botToken
	}
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		slog.ErrorContext(
			ctx,
			fmt.Sprintf("catch a panic"),
			slog.Any("panic", p),
		)
		panicCounter++
		if panicCounter >= panicThreshold {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("too many panics: %v", panicCounter),
				slog.Any("panic", p),
			)
			return
		}
		run()
	}()
	server, err := bot.NewServer(config)
	if err != nil {
		return fmt.Errorf("can not create a new server: %w", err)
	}
	defer server.Shutdown()
	return server.RunBot(ctx)
}
