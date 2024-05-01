package bus

import (
	"context"
	"reflect"

	uow "github.com/niko-labs/libs-go/uow"
)

func (b *bus) RegisterCommandHandler(command CommandHandler, handler CommandHandlerFunc) error {
	cmdName := reflect.TypeOf(command).Name()
	if _, ok := b.commands[cmdName]; ok {
		return ErrorCommandHandlerAlreadyRegistered
	}
	b.commands[cmdName] = handler
	return nil
}

func (b *bus) RemoveCommandHandler(command CommandHandler) error {
	cmdName := reflect.TypeOf(command).Name()
	if _, ok := b.commands[cmdName]; ok {
		delete(b.commands, cmdName)
		return nil
	}
	return ErrorCommandHandlerNotFound
}

func (b *bus) SendCommand(ctx context.Context, command CommandHandler, uow *uow.UnitOfWork) (data any, erro error) {
	cmdName := reflect.TypeOf(command).Name()
	if handler, ok := b.commands[cmdName]; ok {
		return handler(ctx, uow, command)
	}
	return nil, ErrorCommandHandlerNotFound
}

func (b *bus) NumberOfCommandHandlers() int {
	return len(b.commands)
}
