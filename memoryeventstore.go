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

type MemoryEventStore struct {
	connString string
	kindStore  map[uint32]*MemoryEventStoreKindPartition
}

type MemoryEventStoreKindPartition struct {
	kind           *AggregateKind
	aggregateStore map[int64]*MemoryEventStoreAggregatePartition
}

type MemoryEventStoreAggregatePartition struct {
	id     int64
	events []*EventStoreEntry
}

func NewMemoryEventStore(connString string) EventStorer {
	return &MemoryEventStore{
		connString: connString,
		kindStore:  make(map[uint32]*MemoryEventStoreKindPartition),
	}
}

func (es *MemoryEventStore) Kind(kind *AggregateKind) KindPartitioner {
	partition, foundPartition := es.kindStore[kind.Hash()]
	if !foundPartition {
		partition = &MemoryEventStoreKindPartition{
			kind:           kind,
			aggregateStore: make(map[int64]*MemoryEventStoreAggregatePartition), //[][]byte),
		}
		es.kindStore[kind.Hash()] = partition
	}
	return partition
}

func (kindPartition *MemoryEventStoreKindPartition) Aggregate(id int64) AggregatePartitioner {
	partition, foundPartition := kindPartition.aggregateStore[id]
	if !foundPartition {
		partition = &MemoryEventStoreAggregatePartition{
			id:     id,
			events: make([]*EventStoreEntry, 0, 32),
		}
		kindPartition.aggregateStore[id] = partition
	}
	return partition
}

func (aggregatePartition *MemoryEventStoreAggregatePartition) LoadAll() ([]*EventStoreEntry, error) {
	return aggregatePartition.LoadIndexRange(0, MaxInt)
}

func (partition *MemoryEventStoreAggregatePartition) LoadIndexRange(startIndex int, endIndex int) ([]*EventStoreEntry, error) {
	if startIndex < 0 {
		return nil, errors.New("Invalid startIndex")
	}
	if endIndex <= startIndex {
		return nil, errors.New("End index should be greater than start index")
	}
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
}

func (partition *MemoryEventStoreAggregatePartition) Append(entry *EventStoreEntry) error {
	partition.events = append(partition.events, entry)
	return nil
}
