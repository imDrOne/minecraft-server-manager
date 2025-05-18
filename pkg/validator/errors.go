package validator

import (
	"errors"
	"fmt"
	"slices"
)

type ValidateError struct {
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (e *ValidateError) Error() string {
	var str string
	for _, err := range e.Errors {
		str += fmt.Sprintf("\t%s: %s\n", err.Name, err.Message)
	}
	return str
}

func (e *FieldError) Error() string {
	return e.Message
}

func (e *ValidateError) Is(target error) bool {

	var err *ValidateError
	ok := errors.As(target, &err)
	if !ok || e == nil || err == nil {
		return false
	}

	return slices.Equal(e.Errors, err.Errors)
}
