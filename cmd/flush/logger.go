package flush

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	logger logrus.Logger
}

// InitLogger configures the logrus logger with fancy formatting
func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&CustomFormatter{})
	logger.SetLevel(logrus.DebugLevel) // Enable all log levels for demo
	return logger
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor func(a ...any) string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = color.New(color.FgHiCyan).SprintFunc()
	case logrus.InfoLevel:
		levelColor = color.New(color.FgHiGreen).SprintFunc()
	case logrus.WarnLevel:
		levelColor = color.New(color.FgHiYellow).SprintFunc()
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgHiRed).SprintFunc()
	default:
		levelColor = color.New(color.FgHiWhite).SprintFunc()

	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")

	// Create message
	var msg string
	if len(entry.Data) > 0 {
		fields := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
		msg = fmt.Sprintf("%s [%s] %s (%s)", timestamp, levelColor(entry.Level), entry.Message, levelColor(fields))
	} else {
		msg = fmt.Sprintf("%s [%s] %s", timestamp, levelColor(entry.Level), entry.Message)
	}

	// Add emojis
	switch entry.Level {
	case logrus.InfoLevel:
		msg = "ℹ️  " + msg
	case logrus.WarnLevel:
		msg = "⚠️  " + msg
	case logrus.ErrorLevel, logrus.FatalLevel:
		msg = "❌ " + msg
	case logrus.DebugLevel:
		msg = "🔍 " + msg
	}

	return []byte(msg + "\n"), nil
}
