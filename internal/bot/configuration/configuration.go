package configuration

type StartupConfig struct {
	BotAPIToken     string `env:"BotAPIToken,required,notEmpty"`
	DatabaseURL     string `env:"DatabaseURL,required,notEmpty"`
	Debug           bool   `env:"Debug"`
	TGUpdateTimeout int    `env:"UpdateTimeout" envDefault:"60"`
}

type PopulatorConfig struct {
	DatabaseURL    string `env:"DatabaseURL,required,notEmpty"`
	SourceFilename string `env:"PopulatorSource,required,notEmpty"`
}
