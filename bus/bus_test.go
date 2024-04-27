package bus

import (
	"context"
	"libs/uow"
	"log"
	"testing"
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

func baseCommandHandler01(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error) {
	// Command handler that does nothing
	return nil, nil
}
func baseCommandHandler02(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error) {
	// Command handler that does nothing
	return nil, nil
}
func baseCommandHandler03(ctx context.Context, uow *uow.UnitOfWork, cmd CommandHandler) (data any, erro error) {
	// Command handler that will be removed
	return nil, nil
}

func TestMain(t *testing.M) {
	// The BUS uses a global to store the handlers and events
	// This creates a problem when running tests, because the global is not reset between tests

	log.Default().Println("--> Starting Mocks...")

	bus := GetGlobal()
	bus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	bus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	bus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)
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
	localBus := createBase()
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
	localBus := createBase()
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
	localBus := createBase()

	localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)
	localBus.RegisterCommandHandler(cmdEx02{}, baseCommandHandler02)
	localBus.RegisterCommandHandler(cmdEx03{}, baseCommandHandler03)

	if localBus.NumberOfCommandHandlers() != 3 {
		t.Error("---> The number of command handlers should be 3")
	}
}

func TestBus_SendCommandToNonExistingHandler(t *testing.T) {
	localBus := createBase()
	localBus.RegisterCommandHandler(cmdEx01{}, baseCommandHandler01)

	_, err := localBus.SendCommand(context.Background(), cmdEx02{}, nil)
	if err == nil {
		t.Error("---> The command handler should not be found")
	}
}

func TestBus_SendCommandToExistingHandler(t *testing.T) {
	localBus := createBase()
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
