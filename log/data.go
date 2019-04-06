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
	LevelStatus
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
	LevelCrash
)

var levelPrefix = [LevelCrash]string{"[T]", "[D]", "[I]", "[S]", "[N]", "[W]", "[E]", "[F]", "[C]"}
