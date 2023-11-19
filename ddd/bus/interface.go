package bus

import "github.com/INTERNAL-CODE/libs-go/ddd/uow"

type BUSInterface interface {
	New() *BUS
	Get() (*BUS, error)
	HandleCommand(uow uow.UOW, command CommandInterface) (any, error)
	HandleEvent(uow uow.UOW, event EventInterface) error
	RegisterCommand(command CommandInterface, handler func(uow uow.UOW, command CommandInterface) (any, error)) error
	RegisterEvent(event EventInterface, handler func(uow uow.UOW, event EventInterface)) error
}

type CommandInterface interface {
	Identifier() string
	GetData() any
}

type EventInterface interface {
	Identifier() string
	GetData() any
}
