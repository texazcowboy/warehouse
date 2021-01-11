package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type Logger struct {
	LogEntry *LogEntryWithStackTrace
}

type LogEntryWithStackTrace struct {
	*logrus.Entry
}

func (l *LogEntryWithStackTrace) WithError(err error) *logrus.Entry {
	entry := l.Entry
	if stackErr, ok := err.(stackTracer); ok {
		entry = entry.WithField("stacktrace", stackErr.StackTrace())
	}
	return entry.WithError(err)
}

func NewLogger(cfg *LConfig) (*Logger, error) {
	logger := logrus.New()
	logger.ReportCaller = true
	// set logging level.
	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, errors.Wrapf(err, "logger: NewLogger -> logrus.ParseLevel(%v)", cfg.Level)
	}
	logger.SetLevel(lvl)
	// set log format.
	switch cfg.Format {
	case JSON:
		logger.SetFormatter(&logrus.JSONFormatter{})
	case TEXT:
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		return nil, errors.New("logger: Unknown log format: " + string(cfg.Format))
	}
	// set common fields.
	entry := logger.WithField("app_name", cfg.AppName)
	return &Logger{&LogEntryWithStackTrace{entry}}, nil
}
