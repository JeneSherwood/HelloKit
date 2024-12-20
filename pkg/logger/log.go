package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	LevelKey = "logLevel"

	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"

	FormatText = "text"
	FormatJson = "json"
)

type Config struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

func InitLogger(c *Config) {
	var level log.Level
	switch c.Level {
	case DebugLevel:
		level = log.DebugLevel
	case InfoLevel:
		level = log.InfoLevel
	case WarnLevel:
		level = log.WarnLevel
	case ErrorLevel:
		level = log.ErrorLevel
	case FatalLevel:
		level = log.FatalLevel
	default:
		log.Warnf("Invalid loglevel: %s, use default info level\n", c.Level)
		c.Level = InfoLevel
		level = log.InfoLevel
	}

	log.SetLevel(level)

	switch c.Format {
	case FormatText:
		log.SetFormatter(&log.TextFormatter{})
	case FormatJson:
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Warnf("Invalid log format: %s, use default text\n", c.Format)
		c.Format = FormatText
		log.SetFormatter(&log.TextFormatter{})
	}

	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
}
