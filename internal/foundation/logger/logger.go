package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	LogEntry *logrus.Entry
}

func NewLogger(cfg *LConfig) (*Logger, error) {
	logger := logrus.New()
	logger.ReportCaller = true
	// set loggin level.
	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse log level")
	}
	logger.SetLevel(lvl)
	// set log format.
	switch cfg.Format {
	case JSON:
		logger.SetFormatter(&logrus.JSONFormatter{})
	case TEXT:
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		return nil, errors.New("Unknown log format: " + string(cfg.Format))
	}
	// set common fields.
	entry := logger.WithField("app_name", cfg.AppName)
	return &Logger{entry}, nil
}
