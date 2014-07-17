package eventstore

type EventSerializerDeserializer interface {
	EventSerializer
	EventDeserializer
}

type EventSerializer interface {
	SerializeEvent(Event) ([]byte, error)
}

type EventDeserializer interface {
	DeserializeEvent([]byte) (Event, error)
}

type InformedSerializerDeserializer interface {
	EventSerializer
	InformedDeserializer
}

type InformedDeserializer interface {
	InformedDeserializeEvent(uint64, []byte) (Event, error)
}
