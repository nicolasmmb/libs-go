package bus

import (
	"errors"

	"github.com/INTERNAL-CODE/libs-go/ddd/uow"
)

var global *BUS

type BUS struct {
	Commands map[string]func(uow uow.UOWInterface, command CommandInterface) (any, error)
	Event    map[string]func(uow uow.UOWInterface, event EventInterface)
}

func New() *BUS {
	if global != nil {
		return global
	}
	global = &BUS{
		Commands: make(map[string]func(uow uow.UOWInterface, command CommandInterface) (any, error)),
		Event:    make(map[string]func(uow uow.UOWInterface, event EventInterface)),
	}
	return global
}

func Get() (*BUS, error) {
	if global == nil {
		return nil, errors.New("--> BUS not initialized")
	}
	return global, nil
}

func (bus *BUS) HandleCommand(uow uow.UOWInterface, command CommandInterface) (any, error) {
	handler, ok := bus.Commands[command.Identifier()]
	if !ok {
		return nil, errors.New("--> Command not found")
	}
	return handler(uow, command)
}

func (bus *BUS) HandleEvent(uow uow.UOWInterface, event EventInterface) error {
	handler, ok := bus.Event[event.Identifier()]
	if !ok {
		return errors.New("--> Event not found")
	}
	go handler(uow, event)
	return nil
}

func (bus *BUS) RegisterCommand(command CommandInterface, handler func(uow uow.UOWInterface, command CommandInterface) (any, error)) error {
	if bus.Commands[command.Identifier()] != nil {
		return errors.New("--> Command already registered: " + command.Identifier())
	}
	bus.Commands[command.Identifier()] = handler
	return nil
}

func (bus *BUS) RegisterEvent(event EventInterface, handler func(uow uow.UOWInterface, event EventInterface)) error {
	if bus.Event[event.Identifier()] != nil {
		return errors.New("--> Event already registered: " + event.Identifier())
	}
	bus.Event[event.Identifier()] = handler
	return nil
}
