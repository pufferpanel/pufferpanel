package errors

import (
	gErrors "errors"
)

func New(msg string) error {
	return gErrors.New(msg)
}
