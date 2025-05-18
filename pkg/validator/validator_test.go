package validator

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	idError = &ValidateError{
		Errors: []FieldError{
			{
				"Id", "required",
			},
		},
	}
)

type Node struct {
	Id   string `json:"id" validate:"required"`
	Code string `json:"name" validate:"required"`
}

func TestValidateStructWithErrors_Error(t *testing.T) {
	validate := Init()
	err := validate.StructWithErrors(Node{
		Code: "name",
	})

	require.Error(t, err)

	require.ErrorIs(t, err, idError)
}

func TestValidateStructWithErrors_NoError(t *testing.T) {
	validate := Init()
	err := validate.StructWithErrors(Node{
		Id:   "1",
		Code: "name",
	})

	require.NoError(t, err)
}
