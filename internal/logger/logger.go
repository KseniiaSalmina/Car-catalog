package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/KseniiaSalmina/Car-catalog/internal/config"
)

type Logger struct {
	Logger *logrus.Logger
}

func NewLogger(cfg config.Logger) (*Logger, error) {
	l := logrus.New()

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parce log level: %w", err)
	}

	l.SetLevel(lvl)

	return &Logger{Logger: l}, nil
}
