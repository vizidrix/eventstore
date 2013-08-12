package eventstore

import (
	//"errors"
	"fmt"
	"log"
	//"os"
	"time"
)

func ignore_filesystemeventstore() {
	log.Printf(fmt.Sprintf(""))
	time.After(10)
}

type FileSystemES struct {
	connString string
	path       string
	kindStore  map[uint32]KindPartitioner
	//eventStore map[uint32]*FileSystemEventStorePartition
}

type FileSystemESKindPartition struct {
	kind           *AggregateKind
	aggregateStore map[uint64]AggregatePartitioner
}

type FileSystemESAggregatePartition struct {
	id     uint64
	events *EventSet
}

func NewFileSystemES(path string) EventStorer {

	return &FileSystemES{
		connString: path,
		path:       path,
		kindStore:  make(map[uint32]KindPartitioner),
	}
}

func (es *FileSystemES) Kind(kind *AggregateKind) KindPartitioner {
	partition, foundPartition := es.kindStore[kind.Hash()]
	if !foundPartition {
		partition = &FileSystemESKindPartition{
			kind:           kind,
			aggregateStore: make(map[uint64]AggregatePartitioner),
		}
		es.kindStore[kind.Hash()] = partition
	}
	return partition
}

func (kindPartition *FileSystemESKindPartition) Id(id uint64) AggregatePartitioner {
	partition, foundPartition := kindPartition.aggregateStore[id]
	if !foundPartition {
		partition = &FileSystemESAggregatePartition{
			id:     id,
			events: NewEmptyEventSet(),
		}
		kindPartition.aggregateStore[id] = partition
	}
	return partition
}

func (aggregatePartition *FileSystemESAggregatePartition) Get() (*EventSet, error) {
	return aggregatePartition.events, nil
}

func (aggregatePartition *FileSystemESAggregatePartition) GetSlice(startIndex int, endIndex int) (*EventSet, error) {
	return aggregatePartition.events.GetSlice(startIndex, endIndex)
}

func (aggregatePartition *FileSystemESAggregatePartition) Put(newEvents ...Event) (*EventSet, error) {
	newEventSet, err := aggregatePartition.events.Put(newEvents...)
	if err != nil {
		return nil, err
	}
	aggregatePartition.events = newEventSet
	return aggregatePartition.events, nil
}
