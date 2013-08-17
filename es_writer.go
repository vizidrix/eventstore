package eventstore

/*
#include "eventstore.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

const (
	MAX_EVENT_SIZE = 512 - 32
)

type ESWriter struct {
	path      string
	db_handle *[0]byte
}

type ESBatchEntry struct {
	CommandId byte
	EventType uint16
	EventSize uint16
	EventData *[]byte
}

type ESBatch struct {
	BatchId     uint64
	DomainId    uint32
	KindId      uint32
	AggregateId uint64
	BufferSize  byte
	BatchSize   byte
	Entries     []ESBatchEntry
}

func NewESWriter(path string) (*ESWriter, error) {
	str := C.CString(path)
	defer C.free(unsafe.Pointer(str))
	db_handle, err := C.es_open_writer(str)
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
	C.es_close_writer(writer.db_handle)
}

func (writer *ESWriter) AllocBatch(
	domain uint32,
	kind uint32,
	aggregate uint64,
	size byte,
	count byte) (*ESBatch, error) {

	var ES_batch *C.ES_batch
	ES_batch, _ = C.es_alloc_batch(writer.db_handle,
		C.uint32_t(domain),
		C.uint32_t(kind),
		C.uint64_t(aggregate),
		C.char(size),
		C.char(count))

	batch := *(*ESBatch)(unsafe.Pointer(ES_batch))
	entries_header := (*reflect.SliceHeader)((unsafe.Pointer(&batch.Entries)))
	entries_header.Cap = int(count)
	entries_header.Len = int(count)

	//for i := 0; i < int(count); i++ {
	//ptr := (unsafe.Pointer)(batch.Entries[i].EventData)
	//var data *[MAX_EVENT_SIZE]byte = (*[MAX_EVENT_SIZE]byte)(ptr)

	//log.Printf("Data[%d]: % v", i, batch.Entries[i])
	//log.Printf("Ptr: %s", ptr)
	//log.Printf("Data: % v", data)
	//}

	//log.Printf("[ESWriter]\tAllocated batch: % v", batch)

	//log.Printf("Batch entry len: %d", len(batch.Entries))

	return &batch, nil
}

func (batch *ESBatch) Publish() {
	//var ES_batch *[0]byte
	//ES_batch = (*[0]byte)(unsafe.Pointer(batch))
	//C.es_publish_batch(ES_batch)
	C.es_publish_batch((*[0]byte)(unsafe.Pointer(batch)))
}

func (entry *ESBatchEntry) GetEventData() *[MAX_EVENT_SIZE]byte {
	return (*[MAX_EVENT_SIZE]byte)((unsafe.Pointer)(entry.EventData))
}

func (entry *ESBatchEntry) CopyFrom(src []byte) {
	copy((*[MAX_EVENT_SIZE]byte)((unsafe.Pointer)(entry.EventData))[:], src)
}

func es_write_ignore() {
	log.Println(fmt.Sprintf("", 10))
	log.Printf("", reflect.SliceHeader{}, errors.New("stuff"), strings.HasPrefix("s", "q"), unsafe.Pointer(nil))
}
