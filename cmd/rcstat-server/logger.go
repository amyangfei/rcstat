package main

import (
	"github.com/Sirupsen/logrus"
	"io"
)

type LogConfig struct {
	Output io.Writer
	Level  logrus.Level
	Format logrus.Formatter
}

// LogString2Level converts the log level string to logrus Level.
func LogString2Level(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	}

	return logrus.InfoLevel
}
