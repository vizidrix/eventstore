package eventstore

import (
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
	//return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
	return crc32.Checksum(data, table)
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

func (set *EventSet) expandAndCopyHeaders(oldCount int, newCount int, events ...Event) (int, []Header, error) {

	// Switch to range copy and introduce padded allocations and cap checks
	headers := make([]Header, oldCount+newCount)

	if oldCount > 0 {
		copy(headers, UnsafeCastBytesToHeader(set.headers))
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
	size := currentSize + newSize
	//capacity := size | 0xF

	data := make([]byte, size, size)
	copy(data, set.data)
	for i := range events {
		copy(data[currentSize:], events[i].Data)
		currentSize += len(events[i].Data)
	}
	return data
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

/*
func (set *EventSet) PutV2(events ...Event) (*EventSet, error) {
	newCount := len(events)
	oldCount := len(set.headers) / 8

	// Switch to range copy and introduce padded allocations and cap checks
	headers := make([]Header, oldCount+newCount)

	index := 0

	if oldCount > 0 {
		oldHeaders := UnsafeCastBytesToHeader(set.headers)
		oldCount := len(oldHeaders)

		for index = 0; index < oldCount; index++ {
			headers[index] = oldHeaders[index]
		}
	}

	newSize := 0
	// Populate the header for each event
	for i := 0; i < newCount; i++ {
		size := len(events[i].Data)
		// Enforce 2 byte max length in header
		if size > int(MaxUint16) {
			return nil, errors.New("Event data too large")
		}
		newSize += size
		headers[oldCount+i].length = uint16(len(events[i].Data))
		headers[oldCount+i].eventType = events[i].EventType
		headers[oldCount+i].crc = MakeCRC(events[i].Data)
	}

	currentSize := len(set.data)
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
*/

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

func UnsafeCastBytesToHeader(source []byte) []Header {
	length := len(source) / 8 // Bytes / Header struct
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  length,
	}
	return *(*[]Header)(unsafe.Pointer(&result))
}

func UnsafeCastHeaderToBytes(source []Header) []byte {
	length := len(source) * 8 // Bytes / Header struct
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  length,
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

/*
func UnsafeCastInt32sToBytes(ints []int32) []byte {
	length := len(ints) * 4
	hdr := reflect.SliceHeader{Data:uintptr(unsafe.Pointer(&ints[0])), Len: length, Cap: length}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}


*/
