package logger

type LogFormat string

const (
	JSON LogFormat = "json"
	TEXT LogFormat = "text"
)

type LConfig struct {
	AppName string    `yaml:"app-name" validate:"required"`
	Level   string    `yaml:"level" validate:"required"`
	Format  LogFormat `yaml:"format" validate:"required"`
}
