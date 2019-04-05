package config

import "fmt"

type Configer interface {
	Get(key string) (value string)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float64(key string) (float64, error)
	Set(key, value string) error
	ParseFile(filename string) (Configer, error)
}

var Config Configer

func New(adapterName, filename string) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, genError(fmt.Sprintf("Unknow adapter name: %s", adapterName))
	}
	return adapter.ParseFile(filename)
}

func Default(filename string) (err error) {
	Config, err = New(AdapterName_INI, filename)
	return err
}

var adapters = make(map[string]Configer)

func Register(adapterName string, adapter Configer) {
	if adapter == nil {
		panic(genError("register config is nil"))
	}
	if _, ok := adapters[adapterName]; ok {
		panic(genError(fmt.Sprintf("%s has registed", adapterName)))
	}
	adapters[adapterName] = adapter
}
