package bus

import "errors"

var (
	ErrorCommandHandlerAlreadyRegistered = errors.New("--> Command handler already registered")
	ErrorEventHandlerAlreadyRegistered   = errors.New("--> Event handler already registered")
	ErrorCommandHandlerNotFound          = errors.New("--> Command handler not found")
	ErrorEventHandlerNotFound            = errors.New("--> Event handler not found")
)
