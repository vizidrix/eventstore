package eventstore

import (
	"errors"
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
	kindStore  map[uint32]*MemoryESKindPartition
}

type MemoryESKindPartition struct {
	kind           *AggregateKind
	aggregateStore map[int64]*MemoryESAggregatePartition
}

type MemoryESAggregatePartition struct {
	id     int64
	events *EventSet
}

func NewMemoryES(connString string) EventStorer {
	return &MemoryES{
		connString: connString,
		kindStore:  make(map[uint32]*MemoryESKindPartition),
	}
}

func (es *MemoryES) Kind(kind *AggregateKind) KindPartitioner {
	partition, foundPartition := es.kindStore[kind.Hash()]
	if !foundPartition {
		partition = &MemoryESKindPartition{
			kind:           kind,
			aggregateStore: make(map[int64]*MemoryESAggregatePartition), //[][]byte),
		}
		es.kindStore[kind.Hash()] = partition
	}
	return partition
}

func (kindPartition *MemoryESKindPartition) Aggregate(id int64) AggregatePartitioner {
	partition, foundPartition := kindPartition.aggregateStore[id]
	if !foundPartition {
		partition = &MemoryESAggregatePartition{
			id:     id,
			events: NewEmptyEventSet(),
			//events: make([]*EventStoreEntry, 0, 32),
		}
		kindPartition.aggregateStore[id] = partition
	}
	return partition
}

func (aggregatePartition *MemoryESAggregatePartition) Get() (*EventSet, error) {
	return aggregatePartition.GetSlice(0, MaxInt)
}

func (partition *MemoryESAggregatePartition) GetSlice(startIndex int, endIndex int) (*EventSet, error) {
	if startIndex < 0 {
		return nil, errors.New("Invalid startIndex")
	}
	if endIndex <= startIndex {
		return nil, errors.New("End index should be greater than start index")
	}
	return nil, nil
	/*
		// Move end index into range
		if endIndex > len(partition.events) {
			endIndex = len(partition.events) - 1
		}
		// If start point is beyond length bail out
		if startIndex > endIndex || startIndex > len(partition.events) {
			return make([]*EventStoreEntry, 0), nil
		}

		results := make([]*EventStoreEntry, endIndex-startIndex+1)

		copy(results, partition.events[startIndex:endIndex+1])

		return results, nil
	*/
}

func (partition *MemoryESAggregatePartition) Put(eventType uint16, data []byte) error {
	//partition.events = append(partition.events, entry)
	return nil
}
