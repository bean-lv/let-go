package config

import (
	"errors"
	"fmt"
)

var ErrorPrefix = "config error"

func genError(msg string) error {
	return errors.New(fmt.Sprintf("%s: %s", ErrorPrefix, msg))
}
