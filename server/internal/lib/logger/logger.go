package logger

import (
	utils "talk/internal/utils"
)

type Logger struct {
	baseLogger AbstractLogger
	isRoot     bool
	fields     LogFields // only for no root logger
}

func (l *Logger) Error(message string, fields LogFields) {
	l.baseLogger.Log(ErrorLogMode, message, utils.MergeMaps(l.fields, fields))
}

func (l *Logger) Warn(message string, fields LogFields) {
	l.baseLogger.Log(WarnLogMode, message, utils.MergeMaps(l.fields, fields))
}

func (l *Logger) Debug(message string, fields LogFields) {
	l.baseLogger.Log(DebugLogMode, message, utils.MergeMaps(l.fields, fields))
}

func (l *Logger) Info(message string, fields LogFields) {
	l.baseLogger.Log(InfoLogMode, message, utils.MergeMaps(l.fields, fields))
}

func (l *Logger) Err(err error) LogFields {
	return LogFields{"error": err}
}

func (l *Logger) AttachFields(fields LogFields) {
	if l.isRoot {
		l.Error("no available set fields for root logger", nil)
		return
	}

	l.fields = fields
}

var factory LoggerFactory = nil

var Log *Logger // global logger

// Инициализация глобального логгера
func InitGlobalLogger(loggerFactory LoggerFactory) {
	factory = loggerFactory

	Log = &Logger{
		baseLogger: loggerFactory(),
		isRoot:     true,
	}
}

// Для использования внутри сущностей
func NewLogger() *Logger {
	return &Logger{
		baseLogger: factory(),
		isRoot:     false,
	}
}
