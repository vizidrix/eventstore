package eventstore

type EventSerialConverter interface {
	EventSerializer
	EventDeserializer
}

type EventSerializer interface {
	SerializeEvent(Event) ([]byte, error)
}

type EventDeserializer interface {
	DeserializeEvent([]byte) (Event, error)
}

type InformedSerialConverter interface {
	EventSerializer
	InformedDeserializer
}

type InformedDeserializer interface {
	InformedDeserializeEvent(uint64, []byte) (Event, error)
}

/*
func UnmarshalEvent(jsonEvent []byte) (event Event, err error) {
	var memento *EventMemento
	err = json.Unmarshal(jsonEvent, memento)
	if err != nil {
		return
	}
	return memento, err
}
*/
