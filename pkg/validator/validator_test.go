package validator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Node struct {
	Id   string `json:"id" validate:"required"`
	Code string `json:"name" validate:"required"`
}

func TestValidateWithError(t *testing.T) {
	validate := Init()
	err := validate.StructWithErrors(Node{
		Code: "name",
	})

	require.ErrorIs(t, err, &ValidateError{
		Errors: []FieldError{
			{
				"Id", "required1",
			},
		},
	})

}
