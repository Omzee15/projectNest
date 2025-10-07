package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Custom formatter for cleaner output
type CleanFormatter struct {
	TimestampFormat string
	DisableColors   bool
}

func (f *CleanFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor string
	var resetColor string = "\033[0m"

	if !f.DisableColors {
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = "\033[36m" // Cyan
		case logrus.InfoLevel:
			levelColor = "\033[32m" // Green
		case logrus.WarnLevel:
			levelColor = "\033[33m" // Yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = "\033[31m" // Red
		default:
			levelColor = "\033[37m" // White
		}
	}

	timestamp := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())

	// Build the log message
	var message strings.Builder

	// Timestamp and level
	message.WriteString(timestamp)
	message.WriteString(" [")
	message.WriteString(levelColor)
	message.WriteString(level)
	message.WriteString(resetColor)
	message.WriteString("] ")

	// Component if available
	if component, ok := entry.Data["component"]; ok {
		message.WriteString("[")
		message.WriteString(component.(string))
		message.WriteString("] ")
	}

	// Main message
	message.WriteString(entry.Message)

	// Additional fields
	for key, value := range entry.Data {
		if key != "component" {
			message.WriteString(" ")
			message.WriteString(key)
			message.WriteString("=")
			message.WriteString(formatValue(value))
		}
	}

	message.WriteString("\n")

	return []byte(message.String()), nil
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func Init() {
	// Set custom formatter for cleaner output
	env := os.Getenv("APP_ENV")
	if env == "production" {
		// Use JSON formatter in production
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
	} else {
		// Use clean formatter for development
		logrus.SetFormatter(&CleanFormatter{
			TimestampFormat: "15:04:05",
			DisableColors:   false,
		})
	}

	// Set log output
	logrus.SetOutput(os.Stdout)

	// Set log level
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// Helper functions for structured logging
func WithComponent(component string) *logrus.Entry {
	return logrus.WithField("component", component)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

func WithRequest(method, path, clientIP string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"method":    method,
		"path":      path,
		"client_ip": clientIP,
	})
}
