package configuration

type StartupConfig struct {
	BotAPIToken     string `env:"BotAPIToken"`
	DatabaseURL     string `env:"DatabaseURL,required"`
	Debug           bool   `env:"Debug"`
	TGUpdateTimeout int    `env:"UpdateTimeout" envDefault:"60"`
}
