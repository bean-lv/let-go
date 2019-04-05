package log

const (
	AdapterName_File      = "file"
	AdapterName_MultiFile = "multifile"
)

type Level int

const (
	LevelTrace Level = iota + 1
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)
