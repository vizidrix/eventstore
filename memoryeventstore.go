package eventstore

import (
	"errors"
	//"fmt"
	"log"
)

func ignore_memoryeventstore() { log.Printf("") }

type MemoryEventStore struct {
	eventStore map[uint32]*MemoryEventStorePartition
}

type MemoryEventStorePartition struct {
	aggregateStore map[int64][][]byte
}

func NewMemoryEventStore(connString string) EventStorer {
	return &MemoryEventStore{
		eventStore: make(map[uint32]*MemoryEventStorePartition),
	}
}

func (es *MemoryEventStore) RegisterKind(kind *AggregateKind) EventPartitioner {
	partition, foundPartition := es.eventStore[kind.Hash()]
	if !foundPartition {
		partition = &MemoryEventStorePartition{
			aggregateStore: make(map[int64][][]byte),
		}
		es.eventStore[kind.Hash()] = partition
	}
	return partition
}

func (partition *MemoryEventStorePartition) LoadAll(id int64) ([]*EventStoreEntry, error) {
	return partition.LoadIndexRange(id, 0, MaxUint64)
}

// Tests for start / end range checks
func (partition *MemoryEventStorePartition) LoadIndexRange(id int64, startIndex uint64, endIndex uint64) ([]*EventStoreEntry, error) {
	if startIndex < 0 {
		return nil, errors.New("Invalid startIndex")
	}
	if endIndex <= startIndex {
		return nil, errors.New("End index should be greater than start index")
	}
	entryStore, foundEntries := partition.aggregateStore[id]
	// Return an empty slice if the id isn't found
	if !foundEntries {
		return make([]*EventStoreEntry, 0), nil
	}
	// Move end index into range
	if endIndex > uint64(len(entryStore)) {
		endIndex = uint64(len(entryStore)) - 1
	}
	// If start point is beyond length bail out
	if startIndex > endIndex {
		return make([]*EventStoreEntry, 0), nil
	}
	count := endIndex - startIndex + 1
	resultEntries := entryStore[startIndex:]
	result := make([]*EventStoreEntry, len(resultEntries))
	for i := 0; i < int(count); i++ {
		result[i] = FromBinary(resultEntries[i][0:HEADER_SIZE], resultEntries[i][HEADER_SIZE:])
	}

	return result, nil
}

const (
	ENTRYSTORE_INCREMENT = 32
)

func (partition *MemoryEventStorePartition) Append(id int64, entry *EventStoreEntry) error {
	newEntry := make([]byte, HEADER_SIZE+entry.Header().Length())
	copy(newEntry[0:HEADER_SIZE], entry.Header().data)
	copy(newEntry[HEADER_SIZE:], entry.Data())

	entryStore, foundEntries := partition.aggregateStore[id]
	if !foundEntries {
		// Allocate slice of slices with extra
		entryStore = make([][]byte, 1, ENTRYSTORE_INCREMENT)

		entryStore[0] = newEntry
		//partition.aggregateStore[id] = entryStore
	} else {
		position := len(entryStore)
		entryStore = entryStore[0 : position+1]
		entryStore[position] = newEntry

		// If there is no more room to append we have to grow
		if len(entryStore) >= cap(entryStore) {
			//log.Printf("Growing entry store: { %d -> %d }", len(entryStore), cap(entryStore)*2)
			temp := make([][]byte, len(entryStore), cap(entryStore)+ENTRYSTORE_INCREMENT)
			copy(temp[0:len(entryStore)], entryStore)
			entryStore = temp
			//log.Printf("Extended entry store: %d / %d", len(entryStore), cap(entryStore))
		}
	}

	// Allocate and populate a new slice to hold the entry blob
	/*newEntry := make([]byte, HEADER_SIZE+entry.Header().Length())
	copy(newEntry[0:HEADER_SIZE], entry.Header().data)
	copy(newEntry[HEADER_SIZE:], entry.Data())*/

	//log.Printf("Created entry: % v", newEntry)
	// Extend by one for the new item

	//position := len(entryStore)
	//entryStore = entryStore[0 : position+1]
	//entryStore[position] = newEntry

	//log.Printf("From entry store: % v", entryStore[position])
	partition.aggregateStore[id] = entryStore

	return nil
}
