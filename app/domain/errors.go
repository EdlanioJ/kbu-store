package domain

import "errors"

var (
	// ErrNotFound not found
	ErrNotFound = errors.New("entity not found")
	// ErrBlocked entity is blocked
	ErrBlocked = errors.New("entity is blocked")
	// ErrIsPending entity is pending
	ErrPending = errors.New("entity is pending")
	// ErrActived entity is active
	ErrActived = errors.New("entity is active")
	// ErrInactived entity is inactive
	ErrInactived = errors.New("entity is disable")
	// ErrBadRequest bad request
	ErrBadRequest = errors.New("bad request")
	// ErrInternal internal server error
	ErrInternal = errors.New("internal server error")
)
