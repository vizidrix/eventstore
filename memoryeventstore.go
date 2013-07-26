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

func NewMemoryEventStore() *MemoryEventStore {
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
	/*completed := make(chan struct{})
	errored := make(chan error)
	go func() {
		index := 0
		data := es.LoadRaw(uri)
		totalLength := len(data)

		for position := 0; position < totalLength; index++ {
			//event[0] = byte(data[position:position+1] & 0x0F00)
			//event[1] = byte(data[position + 1:position+2] & 0x00F0)
			//event[2] = byte(data[position+2:position:+3] & 0x000F)
			// Load and return the entry at this index
			length := Int24(data[position : position+3])
			log.Printf("Length: %d", length)
			entry := FromBinary(data[position : position+int(length)])

			// Move the position cursor to the next event
			position = position + header_size + int(entry.Length())

			entries <- entry
		}
		//log.Printf("LoadAll Index: %d", index)
		completed <- struct{}{}
	}()
	return completed, errored*/
}

func (es *MemoryEventStore) LoadIndexRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) (completeChan <-chan struct{}, errorChan <-chan error) {
	completed := make(chan struct{})
	errored := make(chan error)
	go func() {
		index := 0
		data := es.LoadRaw(uri)
		totalLength := len(data)

		log.Printf("Loading range:  between {%d, %d} with total len: %d", startIndex, endIndex, totalLength)
		for position := 0; position < totalLength; index++ {
			log.Printf("Entry [%d] between {%d, %d}", index, startIndex, endIndex)
			// If the top bound is reached then abort the loop
			if uint64(index) > endIndex {
				log.Printf("Entry out of top range at: %d", index)
				break
			}
			// Find the length of the current entry's data
			length := Int24(data[position : position+3])
			log.Printf("Entry length: %d", length)
			// Only return entries inside the range
			if uint64(index) >= startIndex {
				log.Printf("Entry in range: %d", index)
				// Load and return the entry at this index
				entry := FromBinary(data[position : position+header_size+int(length)])

				log.Printf("Sending entry: % x", entry)
				entries <- entry
			}
			// Move the position cursor to the next event
			position = position + header_size + int(length)
		}
		completed <- struct{}{}
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
