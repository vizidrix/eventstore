package eventstore

import (
	"errors"
	"time"
)

var (
	ErrInvalidVersion = errors.New("invalid aggregate version")
)

type EventStoreReaderWriter interface {
	AggregateIdGenerater
	EventWriter
	StreamReader
}

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
	LoadEventStreamByEventType(uint32, uint64) ([]Event, error)
	LoadEventStreamByDomain(uint32, uint32) ([]Event, error)
}
