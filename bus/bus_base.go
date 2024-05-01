package bus

import (
	"log"
	"sync"
)

var (
	globalBusInstance *bus
	once              sync.Once
)

type bus struct {
	commands map[string]CommandHandlerFunc
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
	once.Do(func() {
		log.Default().Println("--> Creating new bus instance")
		globalBusInstance = &bus{
			commands: make(map[string]CommandHandlerFunc),
			events:   make(map[string][]EventHandlerFunc),
		}
	})
	return globalBusInstance

}
