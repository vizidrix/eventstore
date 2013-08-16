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

/*
type ESPutCommand struct {
	batch_id     uint64
	command_id   uint64
	crc          uint64
	domain_id    uint64
	kind_id      uint64
	aggregate_id uint64
	event_type   uint16
	event_size   uint16
	Event_data   [1024]byte
}
*/

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
	BatchSize   byte
	Entries     []ESBatchEntry
}

//func (command *ESPutCommand)

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

func (writer *ESWriter) AllocBatch(
	domain uint32,
	kind uint32,
	aggregate uint64,
	count byte) (*ESBatch, error) {

	//max_data_size := 512 - 32

	var ES_batch *C.ES_batch
	ES_batch, _ = C.es_alloc_batch(writer.db_handle,
		C.uint32_t(domain),
		C.uint32_t(kind),
		C.uint64_t(aggregate),
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

	log.Printf("[ESWriter]\tAllocated batch: % v", batch)

	log.Printf("Batch entry len: %d", len(batch.Entries))

	//log.Printf("Batch entry[0] val: % v", *batch.Entries[0].EventData)

	return &batch, nil
}

func (batch *ESBatch) Publish() {
	var ES_batch *[0]byte
	ES_batch = (*[0]byte)(unsafe.Pointer(batch))
	C.es_publish_batch(ES_batch)
}

func (entry *ESBatchEntry) GetEventData() *[MAX_EVENT_SIZE]byte {
	return (*[MAX_EVENT_SIZE]byte)((unsafe.Pointer)(entry.EventData))
}

/*
	max_size := 11 //512 - 32
	for i := 0; i < int(count); i++ {
		ptr := uintptr(unsafe.Pointer(&batch.Entries[i].EventData))
		entry_header := (*reflect.SliceHeader)((unsafe.Pointer(&batch.Entries[i])))
		entry_header.Cap = max_size
		entry_header.Len = max_size
		entry_header.Data = ptr
	}
*/
//header.Data = ptr

/*
	array, _ := C.es_alloc(writer.db_handle, 1)
	var slice []ESPutCommand
	ptr := uintptr(unsafe.Pointer(array))
	header := (*reflect.SliceHeader)((unsafe.Pointer(&slice)))
	header.Cap = 2
	header.Len = 2
	header.Data = ptr

	log.Printf("slice: %s", slice)
*/
//log.Printf("Batch_ptr: %s", batch_ptr)
//batch := (*ESBatch)(unsafe.Pointer(&batch_ptr))

/*
struct ES_batch_entry {
	char			command_id;
	uint16_t		event_type;
	uint16_t		event_size;
	char *			event_data;
};

struct ES_batch {
	uint64_t			batch_id;
	uint32_t			domain_id;
	uint32_t			kind_id;
	uint64_t			aggregate_id;
	ES_batch_entry * 	entries;
};
*/

/* was putting batch (of put command array)
	log.Printf("arrays size: %d", len(array))
	//var slice []ESPutCommand

	ptr := uintptr(unsafe.Pointer(array))
	header := (*reflect.SliceHeader)((unsafe.Pointer(&slice)))
	header.Cap = count
	header.Len = count
	header.Data = ptr
	//slice[0].Event_data[1] = 88
	log.Printf("slice: %s", slice)

	log.Printf("Len: %d", len(slice))

	return &slice[0], nil
}
*/
/*
func (writer *ESWriter) AllocSingle(domain string, kind string, id uint64) (*ESPutCommand, error) {
	//var p *C.ES_put_command
	//p, _ = C.es_alloc(writer.db_handle, 1)

	//var CArray *ESPutCommand
	array, _ := C.es_alloc(writer.db_handle, 1)
	var slice []ESPutCommand
	ptr := uintptr(unsafe.Pointer(array))
	header := (*reflect.SliceHeader)((unsafe.Pointer(&slice)))
	header.Cap = 2
	header.Len = 2
	header.Data = ptr

	log.Printf("slice: %s", slice)

	log.Printf("Len: %d", len(slice[0].Event_data))
	//p := ([2]C.ES_put_command)temp
	//log.Printf("p: %s", p)


		log.Printf("temp: %s\t\t%s", &temp, err)

		//t := (C.ES_put_command)

		ptr := uintptr(unsafe.Pointer(&temp)) //[0]))
		log.Printf("ptr: %s", ptr)

		//put_commands := *(*[]*ESPutCommand)(unsafe.Pointer(&reflect.SliceHeader{
		put_commands := *(*[]*ESPutCommand)(unsafe.Pointer(&reflect.SliceHeader{
			Data: ptr,
			Len:  2,
			Cap:  2,
		}))
		if err != nil {
			return nil, err
		}
		log.Printf("Puts: %s", put_commands)

		//	return put_commands[0], nil


	return nil, nil
}
*/

func es_write_ignore() {
	log.Println(fmt.Sprintf("", 10))
	log.Printf("", reflect.SliceHeader{}, errors.New("stuff"), strings.HasPrefix("s", "q"), unsafe.Pointer(nil))
}
