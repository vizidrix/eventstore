package eventstore

/*
#cgo LDFLAGS: -L/go/vizidrix/src/github.com/vizidrix/ringbuffer/
#include "eventstore.h"
*/
import "C"
import (
	"errors"
	"fmt"
	//"github.com/vizidrix/ringbuffer"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

func event_store_ignore() {
	log.Println(fmt.Sprintf("", 10))
	log.Printf("", reflect.SliceHeader{}, errors.New("stuff"), strings.HasPrefix("s", "q"), unsafe.Pointer(nil))
}

/*
func NewESWriter(path string) (*ESWriter, error) {
	str := C.CString(path)
	defer C.free(unsafe.Pointer(str))
	db_handle, err := C.es_open_write(str)
	if err != nil {
		log.Printf("Error opening write database: %s", err)
		return nil, err
	}
	return &ESWriter{
		path:      path,
		db_handle: db_handle,
	}, nil
}

func (writer *ESWriter) Close() {
	C.es_close_write(writer.db_handle)
}
*/

// http://graphics.stanford.edu/~seander/bithacks.html
//func PowerOf2(value uint64) bool {
//	return value && !(value & (value - 1))
//}

type EventReader interface {
	Get() (*EventSet, error)
	GetSlice(startIndex int, endIndex int) (*EventSet, error)
}

type EventWriter interface {
	Put(newEvents ...Event) (*EventSet, error)
}

type EventStorer interface {
	Kind(kind *AggregateKind) KindPartitioner
	Close()
}

type KindPartitioner interface {
	Id(id uint64) AggregatePartitioner
}

type AggregatePartitioner interface {
	EventReader
	EventWriter
}

func Connect(connString string) (EventStorer, error) {
	//var temp C.rb_buffer
	//C.rb_init_buffer(&temp, 10, 10)
	if strings.HasPrefix(connString, "fs://") {
		//es := EventStore{}
		//es.Connect("/go/vizidrix/src/github.com/vizidrix/eventstore/data/")
		//es := OpenEventStore("/go/esdata/")
		connString = "/go/esdata/"
		//return es, nil
		return NewFileSystemES(connString), nil
	} else if strings.HasPrefix(connString, "mem://") {
		return NewMemoryES(connString), nil
	} else {
		return nil, errors.New("Unable to find delimiter in connection string")
	}
}

type EventStore struct {
}
