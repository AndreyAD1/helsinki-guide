package configuration

type StartupConfig struct {
	BotAPIToken     string `env:"BOT_TOKEN,required,notEmpty"`
	DatabaseURL     string `env:"DATABASE_URL,required,notEmpty"`
	Debug           bool   `env:"DEBUG"`
	TGUpdateTimeout int    `env:"UpdateTimeout" envDefault:"60"`
}

type PopulatorConfig struct {
	DatabaseURL string `env:"DATABASE_URL,required,notEmpty"`
}
