package eventstore

import (
	"log"
)

type FileSystemEventStore struct {
}

func (es *FileSystemEventStore) GetById(uri *AggregateRootUri) ([]EventStoreEntry, error) {
	log.Printf("FileSystemEventStore GetById")
	return nil, nil
}

func (es *FileSystemEventStore) GetByTSRange(uri *AggregateRootUri, startTS int32, endTS int32) ([]EventStoreEntry, error) {
	log.Printf("FileSystemEventStore GetByTSRange")
	return nil, nil
}

func (es *FileSystemEventStore) GetByIndexRange(uri *AggregateRootUri, startIndex uint64, endIndex uint64) ([]EventStoreEntry, error) {
	log.Printf("FileSystemEventStore GetByIndexRange")
	return nil, nil
}

func (es *FileSystemEventStore) Append(uri *AggregateRootUri, entries ...EventStoreEntry) error {
	log.Printf("FileSystemEventStore Append")
	return nil
}
