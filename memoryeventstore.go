package eventstore

import (
	"errors"
	"fmt"
	"log"
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
			log.Printf("Data is: % x", data[position:])
			entry, err := FromBinary(data[position:])
			if err != nil {
				errorChan <- err
			}
			//log.Printf("Appending: %d with len: %d", index, len(entries))
			//Append(entries, *entry)
			//entries[index] = *entry
			//data = data[entry.length:]
			position = position + 17 + int(entry.length)

			entries <- entry
		}
	}()
	return errorChan
	//return entries, nil
}

func (es *MemoryEventStore) LoadTSRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startTS int32, endTS int32) <-chan error {
	errorChan := make(chan error)

	return errorChan
}

func (es *MemoryEventStore) LoadIndexRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) <-chan error {
	errorChan := make(chan error)

	return errorChan
}

func (es *MemoryEventStore) Append(uri *AggregateRootUri, entries ...*EventStoreEntry) <-chan error {
	errorChan := make(chan error)
	go func() {
		for _, entry := range entries {
			data, err := entry.ToBinary()
			if err != nil {
				log.Printf("Error converting to binary: %s", entry)
				errorChan <- err
			}
			prevData := es.data[uri.RelativePath()]
			newData := make([]byte, len(prevData)+len(data))
			copy(newData, prevData)
			copy(newData[len(prevData):], data)
			es.data[uri.RelativePath()] = newData
		}
		return
	}()
	return errorChan
}
