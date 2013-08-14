package eventstore

/*
#include "eventstore.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"strings"
	"unsafe"
)

func event_store_ignore() { log.Println(fmt.Sprintf("", 10)) }

// http://graphics.stanford.edu/~seander/bithacks.html
//func PowerOf2(value uint64) bool {
//	return value && !(value & (value - 1))
//}

const (
	MaxUint   = ^uint(0)
	MinUint   = 0
	MaxInt    = int(^uint(0) >> 1)
	MinInt    = -(MaxInt - 1)
	MaxUint16 = ^uint16(0)
	MinUint16 = 0
	MaxInt16  = int16(^uint16(0) >> 1)
	MinInt16  = -(MaxInt16 - 1)
	MaxUint32 = ^uint32(0)
	MinUint32 = 0
	MaxInt32  = int32(^uint32(0) >> 1)
	MinInt32  = -(MaxInt - 1)
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int64(^uint64(0) >> 1)
	MinInt64  = -(MaxInt - 1)
)

type EventReader interface {
	Get() (*EventSet, error)
	GetSlice(startIndex int, endIndex int) (*EventSet, error)
}

type EventWriter interface {
	Put(newEvents ...Event) (*EventSet, error)
}

type EventStorer interface {
	Kind(kind *AggregateKind) KindPartitioner
}

type KindPartitioner interface {
	Id(id uint64) AggregatePartitioner
}

type AggregatePartitioner interface {
	EventReader
	EventWriter
}

func Connect(connString string) (EventStorer, error) {
	if strings.HasPrefix(connString, "fs://") {
		es := EventStore{}
		//es.Connect("/go/vizidrix/src/github.com/vizidrix/eventstore/data/")
		es.Connect("/go/esdata/")
		return NewFileSystemES(connString), nil
	} else if strings.HasPrefix(connString, "mem://") {
		return NewMemoryES(connString), nil
	} else {
		return nil, errors.New("Unable to find delimiter in connection string")
	}
}

type EventStore struct {
}

func (es *EventStore) Connect(path string) {
	str := C.CString(path)
	defer C.free(unsafe.Pointer(str))
	_, err := C.es_open(str)
	if err != nil {
		log.Printf("Error opening database: %s", err)
	}
	//var handle C.int
	//handle, err := C.godb_open_file(str, 0666)
}

//export DebugPrintf
func DebugPrintf(format *C.char) {
	log.Printf(C.GoString(format))
}
