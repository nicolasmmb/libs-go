package bus

import (
	"context"
	"reflect"

	uow "github.com/niko-labs/libs-go/uow"
)

func (b *bus) RegisterEventHandler(event EventHandler, handler EventHandlerFunc) error {
	eName := reflect.TypeOf(event).Name()
	if handlers, ok := b.events[eName]; ok {
		for _, h := range handlers {
			if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
				return ErrorEventHandlerAlreadyRegistered
			}
		}
		b.events[eName] = append(handlers, handler)
	} else {
		b.events[eName] = []EventHandlerFunc{handler}
	}
	return nil
}

func (b *bus) RemoveEventHandler(event EventHandler) error {
	eName := reflect.TypeOf(event).Name()
	if _, ok := b.events[eName]; ok {
		delete(b.events, eName)
		return nil
	}
	return ErrorEventHandlerNotFound
}

func (b *bus) SendEvent(ctx context.Context, event EventHandler, uow *uow.UnitOfWork) error {
	eName := reflect.TypeOf(event).Name()
	if handlers, ok := b.events[eName]; ok {
		for _, handler := range handlers {
			if err := handler(ctx, uow, event); err != nil {
				return err
			}
		}
		return nil
	}
	return ErrorEventHandlerNotFound
}

func (b *bus) SendEventAsync(ch chan error, ctx context.Context, event EventHandler, uow *uow.UnitOfWork) {
	if ch == nil {
		go b.SendEvent(ctx, event, uow)
	}
	go func() { ch <- b.SendEvent(ctx, event, uow) }()
}

func (b *bus) NumberOfEventHandlers() int {
	return len(b.events)
}
