package shared

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

func GenerateValidationMessage(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		msg := make([]string, 0)
		for _, e := range errs {
			t := e.Field()+": "+e.ActualTag()
			if e.Param() != "" {
				t += " (" + e.Param() + ")"
			}
			msg = append(msg, t)
		}
		return errors.New(strings.Join(msg, ", "))
	}
	return nil
}
