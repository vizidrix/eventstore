package eventstore

import (
	//"errors"
	"fmt"
	"log"
	"time"
)

func ignore_chaneventstore() {
	log.Printf(fmt.Sprintf(""))
	time.After(10)
}

type MemoryES struct {
	connString string
	//kindStore  map[uint32]*MemoryESKindPartition
	kindStore map[uint32]KindPartitioner
}

type MemoryESKindPartition struct {
	kind           *AggregateKind
	aggregateStore map[uint64]AggregatePartitioner
}

type MemoryESAggregatePartition struct {
	id     uint64
	events *EventSet
}

func NewMemoryES(connString string) EventStorer {
	return &MemoryES{
		connString: connString,
		kindStore:  make(map[uint32]KindPartitioner),
	}
}

func (es *MemoryES) Kind(kind *AggregateKind) KindPartitioner {
	partition, foundPartition := es.kindStore[kind.Hash()]
	if !foundPartition {
		partition = &MemoryESKindPartition{
			kind:           kind,
			aggregateStore: make(map[uint64]AggregatePartitioner),
		}
		es.kindStore[kind.Hash()] = partition
	}
	return partition
}

func (kindPartition *MemoryESKindPartition) Id(id uint64) AggregatePartitioner {
	partition, foundPartition := kindPartition.aggregateStore[id]
	if !foundPartition {
		partition = &MemoryESAggregatePartition{
			id:     id,
			events: NewEmptyEventSet(),
		}
		kindPartition.aggregateStore[id] = partition
	}
	return partition
}

func (aggregatePartition *MemoryESAggregatePartition) Get() (*EventSet, error) {
	//return aggregatePartition.events.Get()
	return aggregatePartition.events, nil
}

func (aggregatePartition *MemoryESAggregatePartition) GetSlice(startIndex int, endIndex int) (*EventSet, error) {
	return aggregatePartition.events.GetSlice(startIndex, endIndex)
}

func (aggregatePartition *MemoryESAggregatePartition) Put(newEvents ...Event) (*EventSet, error) {
	newEventSet, err := aggregatePartition.events.Put(newEvents...)
	if err != nil {
		return nil, err
	}
	aggregatePartition.events = newEventSet
	return aggregatePartition.events, nil
}
