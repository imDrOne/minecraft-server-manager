package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func Init() Validator {
	return Validator{validator.New()}
}

func (v *Validator) StructWithErrors(stc interface{}) *ValidateError {
	err := v.Struct(stc)
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		var fieldErrors = make([]FieldError, 0, len(validateErrs))
		for _, err := range validateErrs {
			fieldErrors = append(fieldErrors, FieldError{
				Name:    err.Field(),
				Message: err.Tag(),
			})
		}
		return &ValidateError{fieldErrors}
	}

	return nil
}
