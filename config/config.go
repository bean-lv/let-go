package config

type Config interface {
	LoadConfig(filename string)
	Get(key string) (value interface{})
	Set(key string, value interface{})
}

type myConfig struct {
	data map[string]interface{}
}

func New() Config {
	return &myConfig{}
}

func Default() Config {
	config := New()

	return config
}

func (config *myConfig) LoadConfig(filename string) {

}

func (config *myConfig) Get(key string) interface{} {
	return config.data[key]
}

func (config *myConfig) Set(key string, value interface{}) {
	config.data[key] = value
}
