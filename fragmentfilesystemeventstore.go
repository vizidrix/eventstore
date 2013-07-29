package eventstore

import (
	"log"
)

type FragmentFileSystemEventStore struct {
}

func NewFragmentFileSystemEventStore() EventStorer {
	return &FileSystemEventStore{}
}

func (es *FragmentFileSystemEventStore) LoadAll(uri *AggregateUri, entries chan<- *EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	log.Printf("FileSystemEventStore LoadAll")
	completed := make(chan struct{})
	errored := make(chan error)

	return completed, errored
}

func (es *FragmentFileSystemEventStore) LoadIndexRange(uri *AggregateUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) (completeChan <-chan struct{}, errorChan <-chan error) {
	completed := make(chan struct{}, 1)
	errored := make(chan error)

	return completed, errored
}

func (es *FragmentFileSystemEventStore) Append(uri *AggregateUri, entries ...*EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	log.Printf("FileSystemEventStore Append")
	completed := make(chan struct{})
	errored := make(chan error)

	return completed, errored
}
