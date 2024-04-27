package bus

import "errors"

var (
	ErrorCommandHandlerAlreadyRegistered = errors.New("--> Command handler already registered")
	ErrorCommandHandlerNotFound          = errors.New("--> Command handler not found")
)
