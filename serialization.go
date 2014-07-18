package eventstore

type EventConverter interface {
	EventSerializer
	EventDeserializer
}

type EventSerializer interface {
	SerializeEvent(Event) ([]byte, error)
}

type EventDeserializer interface {
	DeserializeEvent([]byte) (Event, error)
}

type InformedConverter interface {
	EventSerializer
	InformedDeserializer
}

type InformedDeserializer interface {
	InformedDeserializeEvent(uint64, []byte) (Event, error)
}
