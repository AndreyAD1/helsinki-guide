package configuration

type StartupConfig struct {
	BotAPIToken     string `env:"BOT_TOKEN,required,notEmpty"`
	DatabaseURL     string `env:"DATABASE_URL,required,notEmpty"`
	Debug           bool   `env:"DEBUG"`
	TGUpdateTimeout int    `env:"UPDATE_TIMEOUT" envDefault:"60"`
	MetricsUser     string `env:"METRICS_USER,required,notEmpty"`
	MetricsPassword string `env:"METRICS_PASSWORD,required,notEmpty"`
	MetricsPort     int    `env:"METRICS_PORT,required,notEmpty"`
}

type PopulatorConfig struct {
	DatabaseURL string `env:"DATABASE_URL,required,notEmpty"`
	ConverterURL string `env:"CONVERTER_URL" envDefault:"https://epsg.io"`
}
