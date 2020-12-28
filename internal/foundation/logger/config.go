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
	AppName string    `yaml:"app-name"`
	Level   LogLevel  `yaml:"level"`
	Format  LogFormat `yaml:"format"`
}
