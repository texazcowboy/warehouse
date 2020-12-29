package logger

type LogLevel string

const (
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
)

type LConfig struct {
	AppName string    `yaml:"app-name" validate:"required"`
	Level   LogLevel  `yaml:"level" validate:"required"`
	Format  LogFormat `yaml:"format" validate:"required"`
}
