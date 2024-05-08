package log

// Logger .
type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	Fatal(string, ...any)
	Panic(string, ...any)
	Log(...any)
}
