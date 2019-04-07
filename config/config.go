package config

import "fmt"

// Configer 配置文件接口类型。
type Configer interface {
	// Get 获取配置的string类型的值。
	Get(key string) (value string)
	// Int 获取配置的int类型的值。
	Int(key string) (int, error)
	// Int64 获取配置的int64类型的值。
	Int64(key string) (int64, error)
	// Bool 获取配置的bool类型的值。
	Bool(key string) (bool, error)
	// Float64 获取配置的float64类型的值。
	Float64(key string) (float64, error)
	// Set 设置配置的值。
	Set(key, value string) error
	// ParseFile 读取给定文件中配置的值。
	ParseFile(filename string) (Configer, error)
}

// Config 配置文件的实现对象。
var Config Configer

// New 新建配置信息。
// @param adapterName 配置的名字；
// @param filenam 配置文件名称。
// @return Configer 配置的实体；
// @return error 错误信息。
func New(adapterName, filename string) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, genError(fmt.Sprintf("Unknow adapter name: %s", adapterName))
	}
	return adapter.ParseFile(filename)
}

// Default 默认配置。
// 默认使用ini类型的配置。
func Default(filename string) (err error) {
	Config, err = New(AdapterName_INI, filename)
	return err
}

// adapters 配置类型的容器。
var adapters = make(map[string]Configer)

// Register 注册配置信息。
func Register(adapterName string, adapter Configer) {
	if adapter == nil {
		panic(genError("register config is nil"))
	}
	if _, ok := adapters[adapterName]; ok {
		panic(genError(fmt.Sprintf("%s has registed", adapterName)))
	}
	adapters[adapterName] = adapter
}
