package validator

import "fmt"

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
	Node
	return str
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s %s", e.Name, e.Message)
}
