package main

import "errors"

var (
	// ErrInvalidFlag returned when some flag is invalid
	ErrInvalidFlag = errors.New("invalid flag")
)
