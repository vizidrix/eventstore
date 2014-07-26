package eventstore

type EventHandler interface {
	LoadEvents() ([]Event, error)
	ApplyEvent(Event) error
}

type EventHandlerMemento struct {
	eventStore EventStoreReaderWriter
	eventTypes map[uint32]uint64
}
