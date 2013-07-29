package eventstore

import (
	//"errors"
	//"fmt"
	"log"
)

func ignore_memoryeventstore() { log.Printf("") }

type MemoryEventStore struct {
	datastore map[string]map[string]map[int64][]byte
}

func NewMemoryEventStore() EventStorer {
	return &MemoryEventStore{
		datastore: make(map[string]map[string]map[int64][]byte),
	}
}

func (es *MemoryEventStore) LoadRaw(uri *AggregateRootUri) []byte {
	namespace, foundNamespace := es.datastore[uri.Namespace()]
	if !foundNamespace {
		namespace = make(map[string]map[int64][]byte)
		es.datastore[uri.namespace] = namespace
	}
	kind, foundKind := namespace[uri.Kind()]
	if !foundKind {
		kind = make(map[int64][]byte)
		namespace[uri.Kind()] = kind
	}
	aggregate, foundAggregate := kind[uri.Id()]
	if !foundAggregate {
		aggregate = make([]byte, 0)
		kind[uri.Id()] = aggregate
	}
	return aggregate
}

func (es *MemoryEventStore) AppendRaw(uri *AggregateRootUri, entry []byte) {
	prevData := es.LoadRaw(uri)
	newData := make([]byte, len(prevData)+len(entry))
	copy(newData, prevData)
	copy(newData[len(prevData):], entry)
	es.datastore[uri.namespace][uri.kind][uri.id] = newData
}

func (es *MemoryEventStore) LoadAll(uri *AggregateRootUri, entries chan<- *EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	return es.LoadIndexRange(uri, entries, 0, MaxUint64)
}

func (es *MemoryEventStore) LoadIndexRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) (completeChan <-chan struct{}, errorChan <-chan error) {
	completed := make(chan struct{}, 1)
	errored := make(chan error)
	go func() {
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
	}()
	return completed, errored
}

func (es *MemoryEventStore) Append(uri *AggregateRootUri, entries ...*EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	completed := make(chan struct{})
	errored := make(chan error)
	go func() {
		for _, entry := range entries {
			data := entry.ToBinary()

			es.AppendRaw(uri, data)
		}
		completed <- struct{}{}
	}()
	return completed, errored
}
