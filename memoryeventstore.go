package eventstore

import (
	"log"
)

type MemoryEventStore struct {
}

func (es *MemoryEventStore) GetById(uri *AggregateRootUri) ([]EventStoreEntry, error) {
	log.Printf("MemoryEventStore GetById")
	return nil, nil
}

func (es *MemoryEventStore) GetByTSRange(uri *AggregateRootUri, startTS int32, endTS int32) ([]EventStoreEntry, error) {
	log.Printf("MemoryEventStore GetByTSRange")
	return nil, nil
}

func (es *MemoryEventStore) GetByIndexRange(uri *AggregateRootUri, startIndex uint64, endIndex uint64) ([]EventStoreEntry, error) {
	log.Printf("MemoryEventStore GetByIndexRange")
	return nil, nil
}

func (es *MemoryEventStore) Append(uri *AggregateRootUri, entries ...EventStoreEntry) error {
	log.Printf("MemoryEventStore Append")
	return nil
}
