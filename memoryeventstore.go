package eventstore

import (
	//"errors"
	//"fmt"
	"log"
)

func ignore_memoryeventstore() { log.Printf("") }

type MemoryEventStore struct {
	datastore map[uint32]map[int64][]byte
}

func NewMemoryEventStore() EventStorer {
	return &MemoryEventStore{
		datastore: make(map[uint32]map[int64][]byte),
	}
}

func (es *MemoryEventStore) LoadRaw(uri *AggregateUri) []byte {
	partition, foundPartition := es.datastore[uri.Hash()]
	if !foundPartition {
		partition = make(map[int64][]byte)
		es.datastore[uri.Hash()] = partition
	}
	aggregate, foundAggregate := partition[uri.Id()]
	if !foundAggregate {
		aggregate = make([]byte, 0)
		partition[uri.Id()] = aggregate
	}
	return aggregate
}

func (es *MemoryEventStore) AppendRaw(uri *AggregateUri, entry []byte) {
	prevData := es.LoadRaw(uri)
	newData := make([]byte, len(prevData)+len(entry))
	copy(newData, prevData)
	copy(newData[len(prevData):], entry)
	es.datastore[uri.Hash()][uri.Id()] = newData
}

func (es *MemoryEventStore) LoadAll(uri *AggregateUri, entries chan<- *EventStoreEntry) error {
	return es.LoadIndexRange(uri, entries, 0, MaxUint64)
}

func (es *MemoryEventStore) LoadAllAsync(uri *AggregateUri, entries chan<- *EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	return es.LoadIndexRangeAsync(uri, entries, 0, MaxUint64)
}

func (es *MemoryEventStore) LoadIndexRange(uri *AggregateUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) error {
	index := uint64(0)
	data := es.LoadRaw(uri)
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
			entry := FromBinary(data[position : position+HEADER_SIZE+entryLength])
			entries <- entry
		}
		// Move the position cursor to the next event
		position = position + HEADER_SIZE + entryLength
	}
	return nil
}

func (es *MemoryEventStore) LoadIndexRangeAsync(uri *AggregateUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) (completeChan <-chan struct{}, errorChan <-chan error) {
	completed := make(chan struct{}, 1)
	errored := make(chan error)
	//go func() {
	defer func() {
		completed <- struct{}{}
	}()
	index := uint64(0)
	data := es.LoadRaw(uri)
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
			entry := FromBinary(data[position : position+HEADER_SIZE+entryLength])
			entries <- entry
		}
		// Move the position cursor to the next event
		position = position + HEADER_SIZE + entryLength
	}
	//}()
	return completed, errored
}

func (es *MemoryEventStore) Append(uri *AggregateUri, entries ...*EventStoreEntry) error {
	for _, entry := range entries {
		data := entry.ToBinary()

		es.AppendRaw(uri, data)
	}
	return nil
}

func (es *MemoryEventStore) AppendAsync(uri *AggregateUri, entries ...*EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	completed := make(chan struct{}, 1)
	errored := make(chan error)
	//go func() {
	defer func() {
		completed <- struct{}{}
	}()
	for _, entry := range entries {
		data := entry.ToBinary()

		es.AppendRaw(uri, data)
	}
	//}()
	return completed, errored
}
