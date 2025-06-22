package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"slices"
	. "talk/internal/lib/logger"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() AbstractLogger {
	logrusLogger := logrus.New()

	logrusLogger.SetLevel(logrus.DebugLevel)
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000",
		SortingFunc:     sortingFunc,
	})

	return &LogrusLogger{
		logger: logrusLogger,
	}
}

func (l *LogrusLogger) Log(mode LogMode, message string, fields LogFields) {
	if fields == nil {
		fields = LogFields{}
	}

	// runtime.Caller(2) чтобы подняться выше в место вызова лога
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fields["file"] = fmt.Sprintf("%s:%d", path.Base(file), line)
	}

	entry := l.logger.WithFields(logrus.Fields(fields))

	switch mode {
	case ErrorLogMode:
		entry.Error(message)
	case WarnLogMode:
		entry.Warn(message)
	case InfoLogMode:
		entry.Info(message)
	case DebugLogMode:
		entry.Debug(message)
	default:
		entry.Info(message)
	}
}

// Переносит поле file в конец списка полей
func sortingFunc(fields []string) {
	index := slices.Index(fields, "file")
	if index >= 0 {
		fields[index], fields[len(fields)-1] = fields[len(fields)-1], fields[index]
	}
}
