package log

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func init() {
	Register(AdapterName_File, newFileLogger())
}

type FileLogger struct {
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

func newFileLogger() Logger {
	return &FileLogger{
		folder: "log",
		suffix: ".log",
		daily:  true,
		delete: 30,
		level:  LevelInfo,
	}
}

func (l *FileLogger) Init(jsonConfig string) error {
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

	err := l.initNewFile()

	return err
}

func (l *FileLogger) initNewFile() error {
	err := os.MkdirAll(l.folder, os.ModePerm)
	if err != nil {
		return genError(fmt.Sprintf("make dir %s error: %v", l.folder, err))
	}
	now := time.Now()

	filename := l.getFilename(now)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return genError(fmt.Sprintf("open file %s error: %v", filename, err))
	}
	l.fileWriter = f
	l.openDate = getDate(now)

	return nil
}

func (l *FileLogger) Trace(f interface{}, args ...interface{}) {
	if l.level > LevelTrace {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelTrace, msg, time.Now())
}

func (l *FileLogger) Debug(f interface{}, args ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelDebug, msg, time.Now())
}

func (l *FileLogger) Info(f interface{}, args ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelInfo, msg, time.Now())
}

func (l *FileLogger) Status(f interface{}, args ...interface{}) {
	if l.level > LevelStatus {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelStatus, msg, time.Now())
}

func (l *FileLogger) Notice(f interface{}, args ...interface{}) {
	if l.level > LevelNotice {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelNotice, msg, time.Now())
}

func (l *FileLogger) Warn(f interface{}, args ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelWarn, msg, time.Now())
}

func (l *FileLogger) Error(f interface{}, args ...interface{}) {
	if l.level > LevelError {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelError, msg, time.Now())
}

func (l *FileLogger) Fatal(f interface{}, args ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelFatal, msg, time.Now())
}

func (l *FileLogger) Crash(f interface{}, args ...interface{}) {
	if l.level > LevelCrash {
		return
	}
	msg := formatMsg(f, args...)
	l.writeMsg(LevelCrash, msg, time.Now())
}

func (l *FileLogger) writeMsg(lvl Level, msg string, when time.Time) {
	msg = l.formatLevelMsg(lvl, msg, when)

	l.Lock()
	defer l.Unlock()

	l.checkFileWriter(when)

	l.fileWriter.WriteString(msg)
}

func (l *FileLogger) checkFileWriter(when time.Time) {
	if !l.daily {
		return
	}
	if getDate(when) != l.openDate {
		filename := l.getFilename(when)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file writer: %s error: %v", filename, err)
			return
		}
		l.fileWriter.Close()
		l.fileWriter = f

		go l.handleExpireFiles(when)
	}
	return
}

func (l *FileLogger) handleExpireFiles(when time.Time) {
	// TODO:
	if l.compress > 0 {
		// 处理需要压缩的文件。
	}
	if l.delete > 0 {
		// 处理需要删除的文件。
	}
}

func (l *FileLogger) getFilename(when time.Time) string {
	var filename string
	if len(l.prefix) > 0 {
		filename += l.prefix + "-"
	}
	if l.daily {
		filename += fmt.Sprintf("%d", getDate(when))
	} else {
		filename += "log"
	}
	filename += l.suffix
	if len(l.folder) > 0 {
		filename = path.Join(l.folder, filename)
	}
	return filename
}

func (l *FileLogger) formatLevelMsg(lvl Level, msg string, when time.Time) string {
	from := ""
	if l.enableDepth {

	}
	msg = fmt.Sprintf("%s: %s %s%s\n", formatMsgTime(when), levelPrefix[lvl], from, msg)
	return msg
}

func getDate(t time.Time) int {
	yy, mm, dd := t.Date()
	return yy*10000 + int(mm)*100 + dd
}

func formatMsgTime(when time.Time) string {
	return when.Format("2006-01-02 15:04:05")
}

func formatMsg(f interface{}, args ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(args) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {

		} else {
			msg += strings.Repeat("%v ", len(args))
		}
	default:
		msg = fmt.Sprint(f)
		if len(args) == 0 {
			return msg
		}
		msg += strings.Repeat("%v ", len(args))
	}
	msg = strings.TrimSpace(msg)
	return fmt.Sprintf(msg, args...)
}
