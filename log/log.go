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

func NewLogger() Logger {
	return &myLogger{}
}

func (l *myLogger) Trace(i ...interface{}) {}
func (l *myLogger) Debug(i ...interface{}) {}
func (l *myLogger) Info(i ...interface{})  {}
func (l *myLogger) Warn(i ...interface{})  {}
func (l *myLogger) Error(i ...interface{}) {}
func (l *myLogger) Fatal(i ...interface{}) {}
