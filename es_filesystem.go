package eventstore

/*
#include "eventstore.h"
*/
import "C"
import (
	//"errors"
	"fmt"
	"log"
	//"os"
	"time"
	//"unsafe"
)

func ignore_filesystemeventstore() {
	log.Printf(fmt.Sprintf(""))
	time.After(10)
}

type FileSystemES struct {
	connString string
	path       string
	//db_handle  *[0]byte
	domain    uint32
	kindStore map[uint32]KindPartitioner
	//eventStore map[uint32]*FileSystemEventStorePartition
}

type FileSystemESKindPartition struct {
	domain         uint32
	kind           *AggregateKind
	aggregateStore map[uint64]AggregatePartitioner
}

type FileSystemESAggregatePartition struct {
	domain uint32
	kind   uint32
	id     uint64
	events *EventSet
}

func NewFileSystemES(path string) EventStorer {
	//str := C.CString(path)
	//defer C.free(unsafe.Pointer(str))
	//db_handle, err := C.es_open(str)
	//if err != nil {
	//	log.Printf("Error opening database: %s", err)
	//}
	return &FileSystemES{
		connString: path,
		path:       path,
		//db_handle:  db_handle,
		domain:    MakeCRC([]byte(path)),
		kindStore: make(map[uint32]KindPartitioner),
	}
}

func (es *FileSystemES) Close() {
	//C.es_close(es.db_handle)
}

func (es *FileSystemES) Kind(kind *AggregateKind) KindPartitioner {
	partition, foundPartition := es.kindStore[kind.Hash()]
	if !foundPartition {
		partition = &FileSystemESKindPartition{
			domain:         es.domain,
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
			domain: kindPartition.domain,
			kind:   kindPartition.kind.Hash(),
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
	// Retrieve enough put commands to handle the new events
	//puts, err := C.es_allocate_puts(len(newEvents))
	// Populate the puts
	/*for _, i := range newEvents {
		puts[i].domain_id = aggregatePartition.domain,
		puts[i].kind_id = aggregatePartition.kind,
		puts[i].aggregate_id = aggregatePartition.id,1
	}*/

	newEventSet, err := aggregatePartition.events.Put(newEvents...)
	if err != nil {
		return nil, err
	}
	aggregatePartition.events = newEventSet
	return aggregatePartition.events, nil
}
