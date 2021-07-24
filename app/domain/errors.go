package domain

import "errors"

var (
	// ErrNotFound not found
	ErrNotFound = errors.New("Not found")
	// ErrBlocked already blocked
	ErrBlocked = errors.New("Already blocked")
	// ErrIsPending is still pending
	ErrIsPending = errors.New("Is still pending")
	// ErrActived is already active
	ErrActived = errors.New("Is already active")
	// ErrInactived is already inactive
	ErrInactived = errors.New("Is already inactive")
	// ErrBadParam invalid param
	ErrBadParam = errors.New("Param is not valid")
	// ErrInternal internal server error
	ErrInternal = errors.New("internal server error")
)
