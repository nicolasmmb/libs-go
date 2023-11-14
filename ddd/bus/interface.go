package bus

import "github.com/INTERNAL-CODE/libs-go/ddd/uow"

type BUSInterface interface {
	New() *BUS
	Get() (*BUS, error)
	HandleCommand(uow uow.UOWInterface, command CommandInterface) (any, error)
	HandleEvent(uow uow.UOWInterface, event EventInterface) error
	RegisterCommand(command CommandInterface, handler func(uow uow.UOWInterface, command CommandInterface) (any, error)) error
	RegisterEvent(event EventInterface, handler func(uow uow.UOWInterface, event EventInterface)) error
}

type CommandInterface interface {
	Identifier() string
	GetData() any
}

type EventInterface interface {
	Identifier() string
	GetData() any
}
