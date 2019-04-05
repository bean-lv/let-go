package letgo

import (
	"letgo/config"
)

var Config config.Configer

func DefaultConfig(filename string) error {
	err := config.Default(filename)
	if err != nil {
		return err
	}
	Config = config.Config
	return nil
}
