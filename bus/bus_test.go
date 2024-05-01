package bus

import (
	"context"
	"errors"
	"log"
	"testing"

	uow "github.com/niko-labs/libs-go/uow"
)

type cmdEx01 struct{}
type cmdEx02 struct{}
type cmdEx03 struct{}

func (c cmdEx01) IsCommand() {}
func (c cmdEx02) IsCommand() {}
func (c cmdEx03) IsCommand() {}
func (c cmdEx01) Data() any  { return &c }
func (c cmdEx02) Data() any  { return &c }
func (c cmdEx03) Data() any  { return &c }

type eventEx01 struct{}
type eventEx02 struct{}
type eventEx03 struct{}

func (e eventEx01) IsEvent() {}
func (e eventEx02) IsEvent() {}
func (e eventEx03) IsEvent() {}

func (e eventEx01) Data() any { return &e }
func (e eventEx02) Data() any { return &e }
func (e eventEx03) Data() any { return &e }

func baseCommandHandler01(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error) {
	// Command handler that does nothing
	return nil, nil
}
func baseCommandHandler02(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error) {
	// Command handler that does nothing
	return nil, nil
}
func baseCommandHandler03(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error) {
	return nil, nil
}

func baseEventHandler01(ctx context.Context, uow *uow.UnitOfWork, event EventHandler) error {
	return nil
}
func baseEventHandler02(ctx context.Context, uow *uow.UnitOfWork, event EventHandler) error {
	return nil
}
func baseEventHandler03(ctx context.Context, uow *uow.UnitOfWork, event EventHandler) error {
	return nil
}

func createBaseBusMock() *bus {
	return &bus{
		commands: make(map[string]CommandHandlerFunc),
		events:   make(map[string][]EventHandlerFunc),
	}
}

func TestMain(t *testing.M) {
	// The BUS uses a global to store the handlers and events
	// This creates a problem when running tests, because the global is not reset between tests

	log.Default().Println("--> Starting Mocks...")

	bus := GetGlobal()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
	bus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	bus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	bus.RegisterEventHandler(eventEx03{}, baseEventHandler03)
	log.Default().Println("--> Mocks Ready!...")
	t.Run()
}

func TestBus_CreateBus(t *testing.T) {
	localBus := createBase()
	if localBus == nil {
		t.Error("---> Bus is nil")
	}
}

func TestBus_RemoveCommandHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)

	before := localBus.NumberOfCommandHandlers()
	err := localBus.RemoveCommandHandler(cmdEx01{})
	if err != nil {
		t.Error("If the command handler is registered, it should be removed without errors")
	}
	after := localBus.NumberOfCommandHandlers()
	if before == after {
		t.Error("---> The number of command handlers should decrease")
	}
}

func TestBus_RegisterExistingCommandHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)

	before := localBus.NumberOfCommandHandlers()
	err := localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	if err == nil {
		t.Error("---> If the command handler is already registered, it should return an error")
	}
	after := localBus.NumberOfCommandHandlers()

	if before != after {
		t.Error("---> The number of command handlers should not change")
	}
}

func TestBus_CountCommandHandlers(t *testing.T) {
	localBus := createBaseBusMock()

	localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	localBus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	localBus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)

	if localBus.NumberOfCommandHandlers() != 3 {
		t.Error("---> The number of command handlers should be 3")
	}
}

func TestBus_SendCommandToNonExistingHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RemoveCommandHandler(cmdEx02{})

	_, err := localBus.SendCommand(context.Background(), cmdEx02{}, nil)
	if err == nil {
		t.Error("---> The command handler should not be found")
	}
}

func TestBus_SendCommandToExistingHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)

	_, err := localBus.SendCommand(context.Background(), cmdEx01{}, nil)
	if err != nil {
		t.Error("---> The command handler should be found")
	}
}

func BenchmarkBus_ToGetGlobal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetGlobal()
	}
}

func BenchmarkBus_ToCreateBus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createBase()
	}
}

func BenchmarkBus_ToRegisterCommandHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
}

func BenchmarkBus_ToRemoveCommandHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
	bus.RemoveCommandHandler(cmdEx03{})
	bus.RemoveCommandHandler(cmdEx02{})
	bus.RemoveCommandHandler(cmdEx01{})
}

func BenchmarkBus_ToSendCommandToNonExistingHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
	for i := 0; i < b.N; i++ {
		bus.SendCommand(context.Background(), cmdEx01{}, nil)
	}
}

func BenchmarkBus_ToSendCommandToExistingHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
	for i := 0; i < b.N; i++ {
		bus.SendCommand(context.Background(), cmdEx03{}, nil)
	}
}

func BenchmarkBus_ToCountCommandHandlers(b *testing.B) {
	bus := createBase()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
	for i := 0; i < b.N; i++ {
		bus.NumberOfCommandHandlers()
	}
}

// registe event handler
func TestBus_RegisterEventHandler(t *testing.T) {
	localBus := createBaseBusMock()
	before := localBus.NumberOfEventHandlers()

	err := localBus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	if err != nil {
		t.Error("If the event handler is registered, it should be removed without errors")
	}
	after := localBus.NumberOfEventHandlers()
	if before == after {
		t.Error("---> The number of event handlers should increase, not be the same as: ", before)
	}
}

func TestBus_RemoveEventHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterEventHandler(eventEx01{}, baseEventHandler01)

	before := localBus.NumberOfEventHandlers()
	err := localBus.RemoveEventHandler(eventEx01{})
	if err != nil {
		t.Error("If the event handler is registered, it should be removed without errors")
	}
	after := localBus.NumberOfEventHandlers()
	if before == after {
		t.Error("---> The number of event handlers should decrease")
	}
}

func TestBus_RegisterExistingEventHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterEventHandler(eventEx01{}, baseEventHandler01)

	before := localBus.NumberOfEventHandlers()
	err := localBus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	if err == nil {
		t.Error("---> If the event handler is already registered, it should return an error")
	}
	after := localBus.NumberOfEventHandlers()

	if before != after {
		t.Error("---> The number of event handlers should not change")
	}
}

func TestBus_CountEventHandlers(t *testing.T) {
	localBus := createBaseBusMock()

	localBus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	localBus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	localBus.RegisterEventHandler(eventEx03{}, baseEventHandler03)

	if localBus.NumberOfEventHandlers() != 3 {
		t.Error("---> The number of event handlers should be 3")
	}
}

func TestBus_SendEventToNonExistingHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RemoveEventHandler(eventEx02{})

	err := localBus.SendEvent(context.Background(), eventEx02{}, nil)
	if err == nil {
		t.Error("---> The event handler should not be found")
	}
}

func TestBus_SendEventToExistingHandler(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterEventHandler(eventEx01{}, baseEventHandler01)

	err := localBus.SendEvent(context.Background(), eventEx01{}, nil)
	if err != nil {
		t.Error("---> The event handler should be found")
	}
}

func TestBus_SendEventAsync(t *testing.T) {
	localBus := createBaseBusMock()
	localBus.RegisterEventHandler(eventEx01{}, func(ctx context.Context, uow *uow.UnitOfWork, event EventHandler) error {
		log.Println("---> Mock: Event Handler")
		log.Println("---> Mock: Event Handler Done")
		return errors.New("Mock: Error")
	})

	ch := make(chan error)
	localBus.SendEventAsync(ch, context.Background(), eventEx01{}, nil)

	err := <-ch
	if err == nil {
		t.Error("---> The event handler should return an error")
	}
	if err.Error() != "Mock: Error" {
		t.Error("---> The event handler should return the error message")
	}
}

func BenchmarkBus_ToRegisterEventHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	bus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	bus.RegisterEventHandler(eventEx03{}, baseEventHandler03)
}

func BenchmarkBus_ToRemoveEventHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	bus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	bus.RegisterEventHandler(eventEx03{}, baseEventHandler03)
	bus.RemoveEventHandler(eventEx03{})
	bus.RemoveEventHandler(eventEx02{})
	bus.RemoveEventHandler(eventEx01{})
}

func BenchmarkBus_ToSendEventToNonExistingHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	bus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	bus.RegisterEventHandler(eventEx03{}, baseEventHandler03)
	for i := 0; i < b.N; i++ {
		bus.SendEvent(context.Background(), eventEx01{}, nil)
	}
}

func BenchmarkBus_ToSendEventToExistingHandler(b *testing.B) {
	bus := createBase()
	bus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	bus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	bus.RegisterEventHandler(eventEx03{}, baseEventHandler03)
	for i := 0; i < b.N; i++ {
		bus.SendEvent(context.Background(), eventEx03{}, nil)
	}
}

func BenchmarkBus_ToCountEventHandlers(b *testing.B) {
	bus := createBase()
	bus.RegisterEventHandler(eventEx01{}, baseEventHandler01)
	bus.RegisterEventHandler(eventEx02{}, baseEventHandler02)
	bus.RegisterEventHandler(eventEx03{}, baseEventHandler03)
	for i := 0; i < b.N; i++ {
		bus.NumberOfEventHandlers()
	}
}

func BenchmarkBus_ToSendEventAsyncWithNilChannel(b *testing.B) {
	bus := GetGlobal()

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		bus.SendEventAsync(make(chan error), ctx, eventEx01{}, nil)
	}
}
