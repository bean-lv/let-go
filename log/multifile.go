package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type MultifileLogger struct {
	folder      string   // 日志目录。
	prefix      string   // 日志前缀。
	suffix      string   // 日志后缀。
	fileWriter  *os.File // 日志文件。
	openDate    int      // 打开日志的日期。
	daily       bool     // 是否按天生成日志文件。
	compress    int      // 打包日志的天数。
	delete      int      // 删除日志的天数。
	dev         bool     // 是否开发环境。
	level       Level    // 日志级别。
	enableDepth bool     // 是否打印调用方法所在文件及行数。
	callerDepth int      // 调用方法的深度。

	sync.RWMutex
}

func newMultifileLogger() Logger {
	return &MultifileLogger{
		folder: "log",
		suffix: ".log",
		daily:  true,
		delete: 30,
		level:  LevelInfo,
	}
}

func (l *MultifileLogger) Init(jsonConfig string) error {
	if len(jsonConfig) > 0 {
		err := json.Unmarshal([]byte(jsonConfig), l)
		if err != nil {
			return genError(fmt.Sprintf("init logger error: %v", err))
		}
	}
	if len(l.folder) == 0 {
		l.folder = "log"
	}
	if len(l.suffix) > 0 {
		if !strings.HasPrefix(l.suffix, ".") {
			l.suffix = "." + l.suffix
		}
	} else {
		l.suffix = ".log"
	}
	if l.enableDepth && l.callerDepth == 0 {
		l.callerDepth = 1
	}

	return nil
}

func (l *MultifileLogger) initNewFile() error {
	// err := os.MkdirAll(l.folder, os.ModePerm)
	// if err != nil {
	// 	return genError(fmt.Sprintf("make dir %s error: %v", l.folder, err))
	// }
	// now := time.Now()

	// filename := l.getFilename(now)
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	return genError(fmt.Sprintf("open file %s error: %v", filename, err))
	// }
	// l.fileWriter = f
	// l.openDate = getDate(now)

	return nil
}

func (l *MultifileLogger) Trace(f interface{}, args ...interface{}) {
	if l.level > LevelTrace {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelTrace, msg, time.Now())
}

func (l *MultifileLogger) Debug(f interface{}, args ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelDebug, msg, time.Now())
}

func (l *MultifileLogger) Info(f interface{}, args ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelInfo, msg, time.Now())
}

func (l *MultifileLogger) Status(f interface{}, args ...interface{}) {
	if l.level > LevelStatus {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelStatus, msg, time.Now())
}

func (l *MultifileLogger) Notice(f interface{}, args ...interface{}) {
	if l.level > LevelNotice {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelNotice, msg, time.Now())
}

func (l *MultifileLogger) Warn(f interface{}, args ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelWarn, msg, time.Now())
}

func (l *MultifileLogger) Error(f interface{}, args ...interface{}) {
	if l.level > LevelError {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelError, msg, time.Now())
}

func (l *MultifileLogger) Fatal(f interface{}, args ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelFatal, msg, time.Now())
}

func (l *MultifileLogger) Crash(f interface{}, args ...interface{}) {
	if l.level > LevelCrash {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelCrash, msg, time.Now())
}

func (l *MultifileLogger) writeMsg(lvl Level, msg string, when time.Time) {}
