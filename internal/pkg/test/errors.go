package test

import "errors"

var (
	ErrInternalSql        = errors.New("DB error")
	ErrInternalConnection = errors.New("vault connection error")
)
