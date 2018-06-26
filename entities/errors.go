package entities

import "errors"

var (

	// ErrInternalServer describre error
	ErrInternalServer = errors.New("Internal Server Error")

	// ErrNotFound describre error not found
	ErrNotFound = errors.New("Your requested Item is not found")

	// ErrConflict describre a conflict error
	ErrConflict = errors.New("Your Item already exist")

	// ErrStore describre a store error
	ErrStore = errors.New("Your Item could not be stored")

	// ErrInvalid describre a restriction error
	ErrInvalid = errors.New("Your Item is invalid, review restrictions")
)
