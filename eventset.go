package eventstore

/*
#include <stdlib.h>
*/
import (
	//"C"
	"errors"
	//"bufio"
	//"bytes"
	"fmt"
	//"io"
	"log"
	"reflect"
	//"runtime"
	"unsafe"
	//"math/rand"
	//"os"
	//"os"
	//"encoding/binary"
	//"strings"
	"hash/crc32"
	//"time"
)

func eventset_ignore() {
	log.Printf(fmt.Sprintf(""))
}

var table *crc32.Table = crc32.MakeTable(crc32.Castagnoli)

func MakeCRC(data []byte) uint32 {
	return crc32.Checksum(data, table)
}

type Header struct {
	length    uint16
	eventType uint16
	crc       uint32
}

type Event struct {
	EventType uint16
	Data      []byte
}

type EventSet struct {
	headers []byte
	data    []byte
}

const (
	HEADER_SLICE_SIZE = 64   // Room for 8 appends before expand
	DATA_SLICE_SIZE   = 4096 //* 2048
)

func NewEmptyEventSet() *EventSet {
	return &EventSet{
		headers: make([]byte, 0, HEADER_SLICE_SIZE),
		data:    make([]byte, 0, DATA_SLICE_SIZE),
	}
}

func (set *EventSet) CheckSum() error {
	headers := UnsafeCastBytesToHeader(set.headers)
	position := 0
	for _, header := range headers {
		crc := MakeCRC(set.data[position : position+int(header.length)])
		if crc != header.crc {
			return errors.New("Data appears corrupted")
		}
		position += int(header.length)
	}
	return nil
}

func (set *EventSet) Get() ([]Event, error) {
	return set.GetSlice(0, MaxInt)
}

func (set *EventSet) GetSlice(startIndex int, endIndex int) ([]Event, error) {
	length := len(set.headers) / 8
	// Validate inputs
	if startIndex < 0 || startIndex >= endIndex || startIndex > length {
		return nil, errors.New("Either start or end index is out of range")
	}
	// No data so just return empty slice
	if length == 0 {
		return make([]Event, 0), nil
	}
	// End out of range so move it back
	if endIndex >= length {
		endIndex = length - 1
	}
	// Figure out how many records to return
	count := endIndex - startIndex + 1
	// Look at the headers
	headers := UnsafeCastBytesToHeader(set.headers)
	events := make([]Event, count)

	index := 0
	position := 0
	for i := 0; i <= endIndex; i++ {
		if i < startIndex {
			position += int(headers[i].length)
			continue // Skip events below start index
		}
		events[index].EventType = headers[i].eventType
		events[index].Data = set.data[position : position+int(headers[i].length)]
		position += int(headers[i].length)

		index++
	}
	return events, nil
}

func (set *EventSet) Put(events ...Event) (*EventSet, error) {
	oldCount := len(set.headers) / 8
	newCount := len(events)
	newSize, headers, err := set.expandAndCopyHeaders(oldCount, newCount, events...)
	currentSize := len(set.data)
	if err != nil {
		return nil, err
	}
	data := set.expandAndCopyData(currentSize, newSize, events...)
	return &EventSet{
		headers: UnsafeCastHeaderToBytes(headers),
		data:    data,
	}, nil
}

func (set *EventSet) expandAndCopyHeaders(oldCount int, newCount int, events ...Event) (int, []Header, error) {

	dataSize := len(set.headers)
	dataCap := cap(set.headers)
	requiredSize := dataSize + (newCount * 8)

	var headers []Header
	if requiredSize < dataCap {
		data := set.headers[0 : (oldCount+newCount)*8]
		//log.Printf("set.headers[ %d, %d ]: % x", len(set.headers), cap(set.headers), set.headers)
		headers = UnsafeCastBytesToHeader(data)
		//log.Printf("headers[ %d, %d ]: % x", len(headers), cap(headers), headers)
	} else {
		headers = make([]Header, requiredSize>>3, requiredSize>>2)
		if dataSize > 0 {
			copy(headers, UnsafeCastBytesToHeader(set.headers))
		}
	}

	newSize := 0
	maxSize := int(MaxUint16)
	for i := 0; i < len(events); i++ {
		size := len(events[i].Data)
		if size > maxSize {
			return 0, nil, errors.New("Event data too large")
		}
		newSize += size
		headers[oldCount+i].length = uint16(size)
		headers[oldCount+i].eventType = events[i].EventType
		headers[oldCount+i].crc = MakeCRC(events[i].Data)
	}
	return newSize, headers, nil
}

func (set *EventSet) expandAndCopyData(currentSize int, newSize int, events ...Event) []byte {

	dataSize := len(set.data)          // Correct actually consumed size
	dataCap := cap(set.data)           // Available in backing array
	requiredSize := dataSize + newSize // Number of bytes needed
	var data []byte
	if requiredSize < dataCap { // Simple expand into existing cap
		data = set.data[0:requiredSize]
	} else { // Magic expando sauce needed
		// Ensures that the cap is 16 byte alligned... 0x3F would be 8 byte alligned
		// Account for at least 2 similarly sized, 16 byte alligned adds
		requiredCap := (requiredSize | 0x7F) + ((newSize << 1) | 0x7F)
		//dataArray := new([1024 * 1024]byte)
		//data = dataArray[0:requiredSize]
		data = make([]byte, requiredSize, requiredCap)
		copy(data, set.data)
		/*
			if dataSize <= 64 { // At very small sizes it's fater to do a simply loop copy
				for i := range set.data {
					data[i] = set.data[i]
				}
			} else {
				if dataSize <= 1024 { // Below a certain threshold it's faster to do a simple byte copy
					copy(data, set.data)
				} else { // Working on optimized method of copy for 16 byte alligned data
					copy(data, set.data)
					//copy(UnsafeCastBytesToUint64(data), UnsafeCastBytesToUint64(set.data))
				}
			}
		*/
	}
	// NEED TO TEST VALIDITY OF BOUNDARIES DURING COPY!!  DOES ENDIAN MATTER???
	/*
		size := currentSize + newSize
		// Nearest 4k chunk plus the slice size
		capacity := (size | 0xFFF) + DATA_SLICE_SIZE // currentSize + (newSize << 8)
		//capacity := size | 0xF
		var data []byte
	*/
	/*if currentSize == 0 {
		capacity := size > 4096 ? size : 4096
		data = make([]byte, size, (size+1024) * 2)
	}*/
	/*
		if cap(set.data) < size { // Need to expand and copy
			data = make([]byte, size, capacity)
			copy(data, set.data)
		} else { // Use available capacity
			data = set.data[0:size]
		}
	*/
	for i := range events {
		copy(data[currentSize:], events[i].Data)
		currentSize += len(events[i].Data)
	}
	return data
}

func (set *EventSet) Count() int {
	return len(set.headers) / 8
}

func UnsafeCastBytesToHeader(source []byte) []Header {
	//reflect.Array()
	length := len(source) / 8 // Bytes / Header struct
	capacity := cap(source) / 8
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  capacity,
	}
	return *(*[]Header)(unsafe.Pointer(&result))
}

func UnsafeCastHeaderToBytes(source []Header) []byte {
	length := len(source) * 8 // Bytes / Header struct
	capacity := cap(source) * 8
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  capacity,
	}
	return *(*[]byte)(unsafe.Pointer(&result))
}

func UnsafeCastBytesToUint64(source []byte) []uint64 {
	length := len(source) / 8 // Bytes / uint64
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  length,
	}
	return *(*[]uint64)(unsafe.Pointer(&result))
}

func UnsafeCastUint64ToBytes(source []uint64) []byte {
	length := len(source) * 8 // Bytes / uint64
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  length,
	}
	return *(*[]byte)(unsafe.Pointer(&result))
}
