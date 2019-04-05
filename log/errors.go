package log

import (
	"errors"
	"fmt"
)

var errorPrefix = "log error"

func genError(msg string) error {
	return errors.New(fmt.Sprintf("%s: %s", errorPrefix, msg))
}
