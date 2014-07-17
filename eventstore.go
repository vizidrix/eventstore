package eventstore

import (
	"errors"
	"time"
)

var (
	ErrInvalidVersion = errors.New("invalid aggregate version")
	ErrUsedTimestamp  = errors.New("timestamp used")
	ErrUsedKey        = errors.New("datastore key used")
)

// Responsible for creating valid Application and Domain unique Ids for Aggregates
type AggregateIdGenerater interface {
	GenerateAggregateId(uint32, uint32) (uint64, error)
}

// Responsible for persisting Events to the EventStore
type EventWriter interface {
	AppendEvent(Event) (time.Time, error)
}

// Responsible for serving Streams as queries against the EventStore
type StreamReader interface {
	LoadEventStreamByAggregate(uint32, uint32, uint64) ([]Event, error)
	LoadEventStreamByDomain(uint32, uint32) ([]Event, error)
}

/*
type SerializerEvent interface {
	Event
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
}

func UnmarshalEvent(jsonEvent []byte) (event Event, err error) {
	var memento *EventMemento
	err = json.Unmarshal(jsonEvent, memento)
	if err != nil {
		return
	}
	return memento, err
}
*/
