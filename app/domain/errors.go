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
	// ErrBadRequest bad request
	ErrBadRequest = errors.New("Bad request")
	// ErrInternal internal server error
	ErrInternal = errors.New("internal server error")
)
