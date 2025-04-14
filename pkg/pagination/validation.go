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
		return fmt.Errorf("page must be gt than zero: %w", ErrPagePaginationInvalid)
	}
	return nil
}

func validateSize(v uint64) error {
	if v < 1 || v > 30 {
		return fmt.Errorf("page-size out of range 1 - 30: %w", ErrPagePaginationInvalid)
	}
	return nil
}
