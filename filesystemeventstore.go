package eventstore

import (
	"log"
)

type FileSystemEventStore struct {
}

func NewFileSystemEventStore() *FileSystemEventStore {
	return &FileSystemEventStore{}
}

func (es *FileSystemEventStore) LoadAll(uri *AggregateRootUri, entries chan<- *EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	log.Printf("FileSystemEventStore LoadAll")
	completed := make(chan struct{})
	errored := make(chan error)

	return completed, errored
}

func (es *FileSystemEventStore) LoadIndexRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) (completeChan <-chan struct{}, errorChan <-chan error) {
	log.Printf("FileSystemEventStore LoadIndexRange")
	completed := make(chan struct{})
	errored := make(chan error)

	return completed, errored
}

func (es *FileSystemEventStore) Append(uri *AggregateRootUri, entries ...*EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error) {
	log.Printf("FileSystemEventStore Append")
	completed := make(chan struct{})
	errored := make(chan error)

	return completed, errored
}
