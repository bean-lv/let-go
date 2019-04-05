package log

func init() {
	Register(AdapterName_File, &FileLogger{})
}

type FileLogger struct {
	folder      string // 日志目录。
	prefix      string // 日志前缀。
	suffix      string // 日志后缀。
	daily       bool   // 是否按天生成日志文件。
	currentDate int    // 打开日志的日期。
	compress    int    // 打包日志的天数。
	delete      int    // 删除日志的天数。
	dev         bool   // 是否开发环境。
	level       int    // 日志级别。
}

func (l *FileLogger) SetLogger(adapterName string, config ...string) error {
	return nil
}

func (l *FileLogger) Trace(i ...interface{}) {

}

func (l *FileLogger) Debug(i ...interface{}) {

}

func (l *FileLogger) Info(i ...interface{}) {

}

func (l *FileLogger) Warn(i ...interface{}) {

}

func (l *FileLogger) Error(i ...interface{}) {

}

func (l *FileLogger) Fatal(i ...interface{}) {

}
