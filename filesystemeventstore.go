package eventstore

import (
	"log"
)

type FileSystemEventStore struct {
}

func NewFileSystemEventStore() *FileSystemEventStore {
	return &FileSystemEventStore{}
}

func (es *FileSystemEventStore) LoadAll(uri *AggregateRootUri, entries chan<- *EventStoreEntry) <-chan error {
	log.Printf("FileSystemEventStore LoadAll")
	errorChan := make(chan error)

	return errorChan
}

func (es *FileSystemEventStore) LoadTSRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startTS int32, endTS int32) <-chan error {
	log.Printf("FileSystemEventStore LoadTSRange")
	errorChan := make(chan error)

	return errorChan
}

func (es *FileSystemEventStore) LoadIndexRange(uri *AggregateRootUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) <-chan error {
	log.Printf("FileSystemEventStore LoadIndexRange")
	errorChan := make(chan error)

	return errorChan
}

func (es *FileSystemEventStore) Append(uri *AggregateRootUri, entries ...*EventStoreEntry) <-chan error {
	log.Printf("FileSystemEventStore Append")
	errorChan := make(chan error)

	return errorChan
}
