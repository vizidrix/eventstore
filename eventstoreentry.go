package eventstore

import (
	//"errors"
	"fmt"
	//"hash/crc32"
	"log"

	//"reflect"
	//"unsafe"

//"time"
)

func ignore_eventstore() { log.Printf(fmt.Sprintf("")) }

/*
const (
	BYTES_IN_INT32 = 4
)

split header into separate structure and store as []byte using this to load []Header
method to take []byte and []Header and yield []Entry

func UnsafeCastInt32sToBytes(ints []int32) []byte {
	length := len(ints) * BYTES_IN_INT32
	hdr := reflect.SlicdeHeader{Data:uintptr(unsafe.Pointer(&ints[0])), Len: length, Cap: length}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

Header Index - Per Aggregate with Header Record (8 byte) per event
[ 2 byte		| 2 byte		| 4 byte 	]
[ Length		| EventType		| CRC		]

Max Len: 		65535
Max Event Type: 65535

Need unsafe cast from []Header -> []int64 and vica-versa

Stride data in 64bit chunks by Blocks() count

Actual data is Length() long padded to fit chunks

*/

/*
func PackHeaderIntoBytes(target []byte, length uint16, eventType uint16, crc uint32) error {
	if len(target) < 8 {
		return errors.New("Slice is to small to pack header")
	}
	target[0] = byte(length & 0xFF00 >> 8)
	target[1] = byte(length & 0x00FF)
	target[2] = byte((eventType & 0xFF00) >> 8)
	target[3] = byte((eventType & 0x00FF))
	target[4] = byte((crc & 0xFF000000) >> 24)
	target[5] = byte((crc & 0x00FF0000) >> 16)
	target[6] = byte((crc & 0x0000FF00) >> 8)
	target[7] = byte((crc & 0x000000FF))

	return nil
}

func GetHeaderLength(target []byte) (uint16, error) {
	if len(target) < 8 {
		return 0, errors.New("Slice is to small to be an event header")
	}
	return uint16(target[0])<<8 | uint16(target[1]), nil
}

func GetHeaderBlocks(target []byte) (uint16, error) {
	if len(target) < 8 {
		return 0, errors.New("Slice is to small to be an event header")
	}
	length, err := GetHeaderLength(target)
	if err != nil {
		return 0, err
	}
	if length|0x3f == 0 { // No remainder
		return length / 64, nil
	} else { // Account for the final padded block for remainder
		return (length / 64) + 1, nil
	}
}

func GetHeaderEventType(target []byte) (uint16, error) {
	if len(target) < 8 {
		return 0, errors.New("Slice is to small to be an event header")
	}
	return uint16(target[2])<<8 | uint16(target[3]), nil
}

func GetHeaderCRC(target []byte) (uint32, error) {
	if len(target) < 8 {
		return 0, errors.New("Slice is to small to be an event header")
	}
	return uint32(target[4])<<24 | uint32(target[5])<<16 | uint32(target[6])<<8 | uint32(target[7]), nil
}


type Header struct {
	data []byte
}

func NewHeader(target []byte) (*Header, error) {
	if len(target) != 8 {
		return nil, errors.New("Invalid header length")
	}
	return &Header{
		data: target,
	}, nil
}

func (header *Header) Length() uint16 {
	result, _ := GetHeaderLength(header.data)
	return result
}

func (header *Header) Blocks() uint16 {
	result, _ := GetHeaderBlocks(header.data)
	return result
}

func (header *Header) EventType() uint16 {
	result, _ := GetHeaderEventType(header.data)
	return result
}

func (header *Header) CRC() uint32 {
	result, _ := GetHeaderCRC(header.data)
	return result
}

const (
	HEADER_SIZE    = 9
	MAX_TOTAL_SIZE = 4096
	MAX_EVENT_SIZE = MAX_TOTAL_SIZE - HEADER_SIZE
)

//
type EventStoreEntry struct {
	data []byte
}

func MakeCRC(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}

/*
func UInt12(data []byte) uint16 {
	return uint16(data[0])<<8 | uint16(data[1])<<4 | uint16(data[2])
}
*/
/*
// Length of the trailing data block
func (entry *EventStoreEntry) Length() uint16 {
	return UInt12(entry.data[0:3])
}

// Identifier used by serializer to do it's magic
func (entry *EventStoreEntry) EventType() uint16 {
	return uint16(entry.data[3])<<8 | uint16(entry.data[4])
}

// Checksum of the trailing data block
func (entry *EventStoreEntry) CRC() uint32 {
	return uint32(entry.data[5])<<24 | uint32(entry.data[6])<<16 | uint32(entry.data[7])<<8 | uint32(entry.data[8])
}

func NewEventStoreEntryFrom(eventType uint16, body []byte) *EventStoreEntry {
	return NewEventStoreEntry(uint16(len(body)), eventType, crc32.Checksum(body, crc32.MakeTable(crc32.Castagnoli)), body)
}

// Allows you to create EventStoreEntries directly, no value checking so be careful
func NewEventStoreEntry(length uint16, eventType uint16, crc uint32, body []byte) *EventStoreEntry {
	if length > MAX_EVENT_SIZE {
		panic("Invalid event length")
	}
	if length != uint16(len(body)) {
		panic("Reported length wrong")
	}

	totalLength := HEADER_SIZE + len(body)
	data := make([]byte, totalLength)
	data[0] = byte(length & 0x0F00 >> 8)
	data[1] = byte(length & 0x00F0 >> 4)
	data[2] = byte(length & 0x000F)
	data[3] = byte((eventType & 0xFF00) >> 8)
	data[4] = byte((eventType & 0x00FF))
	data[5] = byte((crc & 0xFF000000) >> 24)
	data[6] = byte((crc & 0x00FF0000) >> 16)
	data[7] = byte((crc & 0x0000FF00) >> 8)
	data[8] = byte((crc & 0x000000FF))
	for i := 0; i < len(body); i++ {
		data[i+HEADER_SIZE] = body[i]
	}
	return &EventStoreEntry{
		data: data,
	}
}

func (entry *EventStoreEntry) ToBinary() []byte {
	return entry.data
}

func FromBinary(data []byte) *EventStoreEntry {
	return &EventStoreEntry{
		data: data,
	}
}

// Byte slice to hold the binary data
func (entry *EventStoreEntry) Data() []byte {
	return entry.data[HEADER_SIZE:]
}

*/
