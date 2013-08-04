package eventstore

import (
	"errors"
	//"bufio"
	//"bytes"
	"fmt"
	//"io"
	"log"
	"reflect"
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

func MakeCRC(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}

type Header struct {
	length    uint16
	eventType uint16
	crc       uint32
}

type EventSet struct {
	headers []byte
	data    []byte
}

type Event struct {
	EventType uint16
	Data      []byte
}

func NewEmptyEventSet() *EventSet {
	return &EventSet{
		headers: make([]byte, 0),
		data:    make([]byte, 0),
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
	newCount := len(events)
	currentSize := len(set.data)
	newSize := 0
	index := 0
	oldCount := 0
	var headers []Header
	if len(set.headers) == 0 {
		headers = make([]Header, oldCount+newCount)
	} else {
		// Copy over the existing headers
		oldHeaders := UnsafeCastBytesToHeader(set.headers)
		oldCount := len(oldHeaders)
		headers = make([]Header, oldCount+newCount)

		for index = 0; index < oldCount; index++ {
			headers[index] = oldHeaders[index]
		}
	}
	// Populate the header for each event
	for i := 0; i < newCount; i++ {
		size := len(events[i].Data)
		// Enforce 2 byte max length in header
		if size > int(MaxUint16) {
			return nil, errors.New("Event data too large")
		}
		newSize += size
		headers[i].length = uint16(len(events[i].Data))
		headers[i].eventType = events[i].EventType
		headers[i].crc = MakeCRC(events[i].Data)
	}

	data := make([]byte, currentSize+newSize)

	// Fill from existing data
	for index = 0; index < currentSize; index++ {
		data[index] = set.data[index]
	}
	// Fill from new event data set(s)
	for i := 0; i < newCount; i++ {
		for j := 0; j < len(events[i].Data); j++ {
			data[index] = events[i].Data[j]
			index++
		}
	}

	return &EventSet{
		headers: UnsafeCastHeaderToBytes(headers),
		data:    data,
	}, nil
}

func (set *EventSet) Count() int {
	return len(set.headers) / 8
}

/*
func (set *EventSet) positionOf(index int) int {

}

func (set *EventSet) LengthOf(index int) uint16 {
	position := index * 8
	return UnpackUint16(set.headers[position : position+1])
}

func (set *EventSet) Slice(startIndex int, endIndex int) *EventSet {

	return &EventSet{
		headers: set.headers[startIndex:endIndex],
		//data:
	}
}
*/

/*
func UnpackUint16(byte []bytes) uint16 {
	return uint16(bytes[0])<<8 | uint16(bytes[1])
}

func GetLengthSlice(data []uint64) []uint16 {
	return MaskAndShiftInt64ToUint16(data, 0xFF000000, 24)
}

func GetEventTypeSlice(data []uint64) []uint16 {
	return MaskAndShiftInt64ToUint16(data, mask, shift)
}
func MaskAndShiftInt64ToUint16(data []uint64, mask uint64, shift int) []uint16 {
	results := make([]uint64, len(data))
	for item, index := range data {
		results[index] = item & mask >> shift
	}
}
*/

func UnsafeCastBytesToHeader(data []byte) []Header {
	length := len(data) / 8 // Bytes / Header struct
	headers := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&data[0])),
		Len:  length,
		Cap:  length,
	}
	return *(*[]Header)(unsafe.Pointer(&headers))
}

func UnsafeCastHeaderToBytes(headers []Header) []byte {
	length := len(headers) * 8 // Bytes / Header struct
	data := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&headers[0])),
		Len:  length,
		Cap:  length,
	}
	return *(*[]byte)(unsafe.Pointer(&data))
}

/*
func UnsafeCastBytesToUint64(data []byte) []uint64 {
	length := len(data) / 8 // Bytes / int64
	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&data[0])),
		Len:  lenth,
		Cap:  length,
	}
	return *(*[]uint64)(unsafe.Pointer(&header))
}
*/
/*
func UnsafeCastInt32sToBytes(ints []int32) []byte {
	length := len(ints) * 4
	hdr := reflect.SliceHeader{Data:uintptr(unsafe.Pointer(&ints[0])), Len: length, Cap: length}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}
*/
