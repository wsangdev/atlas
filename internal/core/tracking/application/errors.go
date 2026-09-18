package application

import "errors"

var (
	// ErrInvalidRange: from debe ser anterior a to.
	ErrInvalidRange = errors.New("rango de fechas invalido")
)
