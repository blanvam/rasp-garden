package entities

import "errors"

var (

	// ErrInternalServer describe error
	ErrInternalServer = errors.New("Internal Server Error")

	// ErrNotFound describe error not found
	ErrNotFound = errors.New("Your requested Item is not found")

	// ErrBrokerReceived describe error not found
	ErrBrokerReceived = errors.New("Error decoding resource payload received from broker")

	// ErrConflict describe a conflict error
	ErrConflict = errors.New("Your Item already exist")

	// ErrStore describe a store error
	ErrStore = errors.New("Your Item could not be stored")

	// ErrInvalid describe a restriction error
	ErrInvalid = errors.New("Your Item is invalid, review restrictions")

	// ErrCancelled describe a client error
	ErrCancelled = errors.New("Operation was cancelled or timed out")

	// ErrNotConnected describe a client connection error
	ErrNotConnected = errors.New("Client not connected")

	// ErrConnected describe a client connection error
	ErrConnected = errors.New("Client connected")

	// ErrCtxDone describe Context error
	ErrCtxDone = errors.New("Context is done")

	// ErrRGPIO describe error connecting to Raspberry Pi GPIO
	ErrRGPIO = errors.New("Error connecting to Raspberry Pi GPIO")

	// ErrInvalidRGPIO describe error for invalid action given
	ErrInvalidRGPIO = errors.New("Invalid action, should be open or close")
)
