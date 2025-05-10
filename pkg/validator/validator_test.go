package validator

import (
	"testing"
	"time"
)

type Node struct {
	Id   string    `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required"`
	Time time.Time `json:"time" `
}

func TestInit(t *testing.T) {

}
