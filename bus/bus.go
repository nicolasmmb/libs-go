package bus

import (
	"context"
	"reflect"

	uow "github.com/niko-labs/libs-go/uow"
)

var (
	globalBusInstance *bus
)

type bus struct {
	handlers map[string]CommandHandlerFunc
	events   map[string][]EventHandlerFunc
}

func GetGlobal() *bus {
	if globalBusInstance != nil {
		return globalBusInstance
	} else {
		globalBusInstance = createBase()
	}
	return globalBusInstance
}

func createBase() *bus {
	return &bus{
		handlers: make(map[string]CommandHandlerFunc),
		events:   make(map[string][]EventHandlerFunc),
	}
}

func (b *bus) RegisterCommandHandler(command CommandHandler, handler CommandHandlerFunc) error {
	cmdName := reflect.TypeOf(command).Name()
	if _, ok := b.handlers[cmdName]; ok {
		return ErrorCommandHandlerAlreadyRegistered
	}
	b.handlers[cmdName] = handler
	return nil
}

func (b *bus) RemoveCommandHandler(command CommandHandler) error {

	cmdName := reflect.TypeOf(command).Name()
	if _, ok := b.handlers[cmdName]; ok {
		delete(b.handlers, cmdName)
		return nil
	}
	return ErrorCommandHandlerNotFound
}

func (b *bus) SendCommand(ctx context.Context, command CommandHandler, uow *uow.UnitOfWork) (data any, erro error) {
	cmdName := reflect.TypeOf(command).Name()
	if handler, ok := b.handlers[cmdName]; ok {
		return handler(ctx, uow, command)
	}
	return nil, ErrorCommandHandlerNotFound
}

func (b *bus) NumberOfCommandHandlers() int {
	return len(b.handlers)
}
