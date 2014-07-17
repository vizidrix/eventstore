package eventstore

import (
	"errors"
	"fmt"
)

// Helper function for Domains to use when defining events
func E(domain uint32, version uint64, typeId uint64) uint64 {
	return (uint64(domain) << 32) | (version & 0x7FFF << 16) | (typeId & 0xFFFF)
}

var (
	ErrInvalidEvent = errors.New("invalid event type")
)

type Event interface {
	GetApplication() uint32
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
	GetEventType() uint64
}

type EventSerializerDeSerializer interface {
	EventSerializer
	EventDeserializer
}

type EventSerializer interface {
	SerializeEvent(Event) ([]byte, error)
}

type EventDeserializer interface {
	DeSerializeEvent(uint64, []byte) (Event, error)
}

type EventMemento struct {
	application uint32 `json:"__application"` // Application the target aggregate belongs to
	domain      uint32 `json:"__domain"`      // The type of aggregate (type is semantically equivalent to doman)
	id          uint64 `json:"__id"`          // Domain-unique identifier for the aggregate instance
	version     uint32 `json:"__version"`     // Derived from the number of events applied to the aggregate
	eventType   uint32 `json:"__etype"`       // Domain-unique identifier for the type of event message
}

func NewEvent(application uint32, domain uint32, id uint64, version uint32, eventType uint64) Event {
	return &EventMemento{
		application: application,
		domain:      domain,
		id:          id,
		version:     version,
		eventType:   uint32((eventType << 32) >> 32),
	}
}

func (event *EventMemento) GetApplication() uint32 {
	return event.application
}

func (event *EventMemento) GetDomain() uint32 {
	return event.domain
}

func (event *EventMemento) GetId() uint64 {
	return event.id
}

func (event *EventMemento) GetVersion() uint32 {
	return event.version
}

func (event *EventMemento) GetEventType() uint64 {
	return (uint64(event.domain) << 32) | uint64(event.eventType)
}

func (event *EventMemento) String() string {
	return fmt.Sprintf(" <E [ <A D[%d] ID[%d] V[%d] \\> -> E[%d] ] E\\> ",
		event.GetDomain(), event.GetId(), event.GetVersion(), event.eventType)
}
