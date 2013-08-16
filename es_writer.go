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

type ESWriter struct {
	path      string
	db_handle *[0]byte
}

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

func (writer *ESWriter) Next() (*ESPutCommand, error) {
	array, _ := C.es_alloc(writer.db_handle, 1)
	log.Printf("arrays size: %d", len(array))
	var slice []ESPutCommand
	ptr := uintptr(unsafe.Pointer(array))
	header := (*reflect.SliceHeader)((unsafe.Pointer(&slice)))
	header.Cap = 3
	header.Len = 3
	header.Data = ptr
	slice[0].Event_data[1] = 88
	log.Printf("slice: %s", slice)

	log.Printf("Len: %d", len(slice[0].Event_data))

	return &slice[0], nil
}

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
