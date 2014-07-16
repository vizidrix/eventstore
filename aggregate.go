package eventstore

import (
	"errors"
	"fmt"
)

type Aggregate interface {
	GetApplication() uint32
	GetDomain() uint32
	GetId() uint64
	GetVersion() uint32
}

var (
	ErrInvalidAggregate = errors.New("invalid aggregate type")
)

type AggregateMemento struct {
	application uint32 `json:"__application"`	// Application this aggregate belongs to
	domain      uint32 `json:"__domain"`  		// The type of aggregate (type is semantically equivalent to doman)
	id          uint64 `json:"__id"`      		// Domain-unique identifier for the aggregate instance
	version     uint32 `json:"__version"` 		// Derived from the number of events applied to the aggregate
}

func NewAggregate(application uint32, domain uint32, id uint64, version uint32) Aggregate {
	return &AggregateMemento{
		application: application,
		domain:      domain,
		id:          id,
		version:     version,
	}
}

func (aggregate *AggregateMemento) GetApplication() uint32 {
	return aggregate.application
}

func (aggregate *AggregateMemento) GetDomain() uint32 {
	return aggregate.domain
}

func (aggregate *AggregateMemento) GetId() uint64 {
	return aggregate.id
}

func (aggregate *AggregateMemento) GetVersion() uint32 {
	return aggregate.version
}

func (aggregate AggregateMemento) String() string {
	return fmt.Sprintf("<A A[%d] D[%d] ID[%d] V[%d] \\>", aggregate.GetApplication(), aggregate.GetDomain(), aggregate.GetId(), aggregate.GetVersion())
}
