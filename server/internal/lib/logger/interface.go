package logger

type AbstractLogger interface {
	Log(mode LogMode, message string, fields LogFields)
}

type LoggerFactory func() AbstractLogger
type LogMode struct{ string } // в структуре для строгой типизации
type LogFields map[string]interface{}

var (
	ErrorLogMode = LogMode{"error"}
	WarnLogMode  = LogMode{"warn"}
	DebugLogMode = LogMode{"debug"}
	InfoLogMode  = LogMode{"info"}
)
