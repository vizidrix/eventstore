package eventstore

import (
	//"errors"
	//"fmt"
	"log"
)

func ignore_filesystemeventstore() { log.Printf("") }

type FileSystemEventStore struct {
	eventStore map[uint32]*FileSystemEventStorePartition
}

type FileSystemEventStorePartition struct {
	aggregateStore map[int64][]byte
}

func NewFileSystemEventStore(connString string) EventStorer {
	return &FileSystemEventStore{
		eventStore: make(map[uint32]*FileSystemEventStorePartition),
	}
}

func (es *FileSystemEventStore) RegisterKind(kind *AggregateKind) EventPartitioner {
	partition, foundPartition := es.eventStore[kind.Hash()]
	if !foundPartition {
		partition = &FileSystemEventStorePartition{
			aggregateStore: make(map[int64][]byte),
		}
		es.eventStore[kind.Hash()] = partition
	}
	return partition
}

func (partition *FileSystemEventStorePartition) LoadAll(id int64) ([]*EventStoreEntry, error) {
	return nil, nil
}

func (partition *FileSystemEventStorePartition) LoadIndexRange(id int64, startIndex uint64, endIndex uint64) ([]*EventStoreEntry, error) {
	return nil, nil
}

func (partition *FileSystemEventStorePartition) LoadAll2(id int64, entries chan<- *EventStoreEntry) error {
	return partition.LoadIndexRange2(id, entries, 0, MaxUint64)
}

func (partition *FileSystemEventStorePartition) LoadIndexRange2(id int64, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) error {
	index := uint64(0)
	data, foundAggregate := partition.aggregateStore[id]
	if !foundAggregate {
		return nil
	}

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
			entry := FromBinary(data[position:position+HEADER_SIZE], data[position+HEADER_SIZE:position+HEADER_SIZE+entryLength])
			entries <- entry
		}
		// Move the position cursor to the next event
		position = position + HEADER_SIZE + entryLength
	}
	return nil
}

func (partition *FileSystemEventStorePartition) Append(id int64, entry *EventStoreEntry) error {
	aggregate, foundAggregate := partition.aggregateStore[id]
	position := 0
	if foundAggregate {
		position = len(aggregate)
	} else {
		//aggregate = make([]byte, 0, PARTITION_BUFFER)
		aggregate = make([]byte, 0)
	}
	header, body := entry.ToBinary()
	//log.Printf("Check for cap: %d -(%d+%d+%d) < %d", cap(aggregate), position, HEADER_SIZE, len(body), PARTITION_BUFFER)
	// Check for room in the capacity and expand the aggregate if needed
	/*if cap(aggregate)-(position+HEADER_SIZE+len(body)) < PARTITION_BUFFER {
		newData := make([]byte, position, cap(aggregate)+PARTITION_BUFFER)
		copy(newData[0:position], aggregate)
		aggregate = newData
	}*/
	newData := make([]byte, position+HEADER_SIZE+len(body))
	copy(newData, aggregate)
	copy(newData[position:], header)
	copy(newData[position+HEADER_SIZE:], body)

	partition.aggregateStore[id] = newData

	return nil
}
