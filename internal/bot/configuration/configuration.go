package configuration

type StartupConfig struct {
	BotAPIToken string `env:"BotAPIToken"`
	DatabaseURL string `env:"DatabaseURL,required"`
}
