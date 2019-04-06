package letgo

import (
	"letgo/log"
)

var Logger log.Logger

func DefaultLogger() error {
	err := log.Default()
	if err != nil {
		return err
	}
	Logger = log.Log
	return nil
}
