package log

import (
	"fmt"
)

type Logger interface {
	SetLogger(adapterName string, config ...string) error

	Trace(i ...interface{})
	Debug(i ...interface{})
	Info(i ...interface{})
	Warn(i ...interface{})
	Error(i ...interface{})
	Fatal(i ...interface{})
}

var Log Logger

func New(adapterName string) (Logger, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, genError(fmt.Sprintf("Unknown adapter name: %s", adapterName))
	}
	return adapter, nil
}

func Default() (err error) {
	Log, err = New(AdapterName_File)
	return err
}

var adapters = make(map[string]Logger)

func Register(adapterName string, adapter Logger) {
	if adapter == nil {
		panic(genError("register logger is nil"))
	}

	if _, ok := adapters[adapterName]; ok {
		panic(genError(fmt.Sprintf("%s has registed", adapterName)))
	}

	adapters[adapterName] = adapter
}
