package log

type Logger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

type myLogger struct {
}

func New() Logger {
	return &myLogger{}
}

func Default() Logger {
	logger := New()

	return logger
}

func (l *myLogger) Trace(i ...interface{}) {}
func (l *myLogger) Debug(i ...interface{}) {}
func (l *myLogger) Info(i ...interface{})  {}
func (l *myLogger) Warn(i ...interface{})  {}
func (l *myLogger) Error(i ...interface{}) {}
func (l *myLogger) Fatal(i ...interface{}) {}
