package bot

import (
	"log"

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
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
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

func run() {
	config := configuration.StartupConfig{}
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}
	if botToken != "" {
		config.BotAPIToken = botToken
	}
	server, err := bot.NewServer(config)
	if err != nil {
		log.Fatalf("can not run a server: %v", err)
	}
	defer server.Shutdown()
	server.RunBot()
}
