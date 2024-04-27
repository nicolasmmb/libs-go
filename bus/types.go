package bus

import (
	"context"
	"libs/uow"
)

type CommandHandlerFunc func(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error)
type EventHandlerFunc func(ctx context.Context, uow *uow.UnitOfWork, event EventHandler) error

type BusInterface interface {
	RegisterCommandHandler(command CommandHandler, handler CommandHandlerFunc) error
	RemoveCommandHandler(command CommandHandler) error

	RegisterEventHandler(event EventHandler, handler EventHandlerFunc) error
	RemoveEventHandler(event EventHandler) error

	SendCommand(command CommandHandler)
}

type CommandHandler interface {
	IsCommand()
	Data() any
}

type EventHandler interface {
	IsEvent()
	Data() any
}
