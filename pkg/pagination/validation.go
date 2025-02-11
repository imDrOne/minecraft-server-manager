package pagination

import (
	"errors"
	"fmt"
)

var (
	ErrPagePaginationInvalid = errors.New("invalid page pagination")
)

func validatePage(v uint64) error {
	if v <= 0 {
		return fmt.Errorf("%w: page must be gt than zero", ErrPagePaginationInvalid)
	}
	return nil
}

func validateSize(v uint64) error {
	if v < 1 || v > 30 {
		return fmt.Errorf("%w: page-size out of range 1 - 30", ErrPagePaginationInvalid)
	}
	return nil
}
