package eventstore

import (
	//"errors"
	//"fmt"
	"log"
)

func ignore_memoryeventstore() { log.Printf("") }

type MemoryEventStore struct {
	eventStore map[uint32]*MemoryEventStorePartition
}

type MemoryEventStorePartition struct {
	aggregateStore map[int64][]byte
}

func NewMemoryEventStore(connString string) EventStorer {
	return &MemoryEventStore{
		eventStore: make(map[uint32]*MemoryEventStorePartition),
	}
}

func (es *MemoryEventStore) RegisterKind(kind *AggregateKind) EventPartitioner {
	partition, foundPartition := es.eventStore[kind.Hash()]
	if !foundPartition {
		partition = &MemoryEventStorePartition{
			aggregateStore: make(map[int64][]byte),
		}
		es.eventStore[kind.Hash()] = partition
	}
	return partition
}

//func (partition *MemoryEventStorePartition) LoadAll(uri *AggregateUri, entries chan<- *EventStoreEntry) error {
func (partition *MemoryEventStorePartition) LoadAll(id int64, entries chan<- *EventStoreEntry) error {
	return partition.LoadIndexRange(id, entries, 0, MaxUint64)
}

func (partition *MemoryEventStorePartition) LoadIndexRange(id int64, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) error {
	index := uint64(0)
	data, foundAggregate := partition.aggregateStore[id]
	if !foundAggregate {
		return nil
	}

	totalLength := len(data)

	for position := 0; position < totalLength; index++ {
		// If the top bound is reached then abort the loop
		if index > endIndex {
			break
		}
		// Find the length of the current entry's data
		entryLength := int(UInt12(data[position : position+3]))
		// Only return entries inside the range
		if index >= startIndex {
			// Load and return the entry at this index
			entry := FromBinary(data[position:position+HEADER_SIZE], data[position+HEADER_SIZE:position+HEADER_SIZE+entryLength])
			entries <- entry
		}
		// Move the position cursor to the next event
		position = position + HEADER_SIZE + entryLength
	}
	return nil
}

func (partition *MemoryEventStorePartition) Append(id int64, entry *EventStoreEntry) error {
	aggregate, foundAggregate := partition.aggregateStore[id]
	header, body := entry.ToBinary()
	position := 0
	if foundAggregate {
		position = len(aggregate)
	} else {
		aggregate = make([]byte, 0)
	}
	newData := make([]byte, position+HEADER_SIZE+len(body))
	copy(newData, aggregate)
	copy(newData[position:], header)
	copy(newData[position+HEADER_SIZE:], body)

	partition.aggregateStore[id] = newData

	return nil
}
