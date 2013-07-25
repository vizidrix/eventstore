package eventstore

import (
	"errors"
	"fmt"
	//"log"
)

type MemoryEventStore struct {
	data map[string][]byte
}

func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		data: make(map[string][]byte),
	}
}

func (es *MemoryEventStore) LoadAll(uri *AggregateRootUri, entries chan<- *EventStoreEntry) <-chan error {
	//entries := make([]EventStoreEntry, 3)
	errorChan := make(chan error)
	go func() {
		index := 0
		data, found := es.data[uri.RelativePath()]
		if !found {
			errorChan <- errors.New(fmt.Sprintf("Item not found: %s", uri.RelativePath()))
		}
		for position := 0; position < len(data); index++ {
			//log.Printf("Data is: % x", data[position:])
			entry := FromBinary(data[position:])

			//log.Printf("Appending: %d with len: %d", index, len(entries))
			//Append(entries, *entry)
			//entries[index] = *entry
			//data = data[entry.length:]
			position = position + header_size + int(entry.Length())

			entries <- entry
		}
	}()
	return errorChan
	//return entries, nil
}

func (es *MemoryEventStore) LoadIndexRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) <-chan error {
	errorChan := make(chan error)

	return errorChan
}

func (es *MemoryEventStore) Append(uri *AggregateRootUri, entries ...*EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	completeC := make(chan struct{})
	errorC := make(chan error)
	go func() {
		for _, entry := range entries {
			data := entry.ToBinary()

			prevData := es.data[uri.RelativePath()]
			newData := make([]byte, len(prevData)+len(data))
			copy(newData, prevData)
			copy(newData[len(prevData):], data)
			es.data[uri.RelativePath()] = newData
		}
		completeC <- struct{}{}
	}()
	return completeC, errorC
}
