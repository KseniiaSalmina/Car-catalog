package config

type Logger struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
