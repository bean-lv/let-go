package log

import (
	"fmt"
)

type Logger interface {
	Init(jsonConfig string) error

	Trace(f interface{}, args ...interface{})
	Debug(f interface{}, args ...interface{})
	Info(f interface{}, args ...interface{})
	Status(f interface{}, args ...interface{})
	Notice(f interface{}, args ...interface{})
	Warn(f interface{}, args ...interface{})
	Error(f interface{}, args ...interface{})
	Fatal(f interface{}, args ...interface{})
	Crash(f interface{}, args ...interface{})
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
	Log.Init("")
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
