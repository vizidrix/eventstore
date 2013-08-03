package eventstore

import (
	//"errors"
	"fmt"
	"hash/crc32"
	"log"

	//"reflect"
	//"unsafe"

//"time"
)

func ignore_eventstore() { log.Printf(fmt.Sprintf("")) }

const (
	HEADER_SIZE    = 9
	MAX_TOTAL_SIZE = 4096
	MAX_EVENT_SIZE = MAX_TOTAL_SIZE - HEADER_SIZE
)

//type EventStoreHeader struct {
//	data []byte
//}
type EventStoreEntry struct {
	//header []byte
	//header *EventStoreHeader
	//body   []byte
	data []byte
}

/*
func NewEventStoreHeaderFrom(length uint16, eventType uint16, crc uint32) *EventStoreHeader {
	return NewEventStoreHeader([]byte{
		byte(length & 0x0F00 >> 8),
		byte(length & 0x00F0 >> 4),
		byte(length & 0x000F),
		byte((eventType & 0xFF00) >> 8),
		byte((eventType & 0x00FF)),
		byte((crc & 0xFF000000) >> 24),
		byte((crc & 0x00FF0000) >> 16),
		byte((crc & 0x0000FF00) >> 8),
		byte((crc & 0x000000FF)),
	})
}

func NewEventStoreHeader(data []byte) *EventStoreHeader {
	return &EventStoreHeader{
		data: data,
	}
}
*/

func UInt12(data []byte) uint16 {
	return uint16(data[0])<<8 | uint16(data[1])<<4 | uint16(data[2])
}

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

/*
const (
	BYTES_IN_INT32 = 4
)

func UnsafeCastInt32sToBytes(ints []int32) []byte {
	length := len(ints) * BYTES_IN_INT32
	hdr := reflect.SlicdeHeader{Data:uintptr(unsafe.Pointer(&ints[0])), Len: length, Cap: length}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}
*/

// Allows you to create EventStoreEntries directly, no value checking so be careful
func NewEventStoreEntry(length uint16, eventType uint16, crc uint32, body []byte) *EventStoreEntry {
	if length > MAX_EVENT_SIZE {
		panic("Invalid event length")
	}
	if length != uint16(len(body)) {
		panic("Reported length wrong")
	}

	/*
		return &EventStoreEntry{
			header: NewEventStoreHeaderFrom(length, eventType, crc),
			body:   body,
		}
	*/
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
	//header := reflect.SliceHeader{}

	/*
		return &EventStoreEntry{
			header: []byte{
				byte(length & 0x0F00 >> 8),
				byte(length & 0x00F0 >> 4),
				byte(length & 0x000F),
				byte((eventType & 0xFF00) >> 8),
				byte((eventType & 0x00FF)),
				byte((crc & 0xFF000000) >> 24),
				byte((crc & 0x00FF0000) >> 16),
				byte((crc & 0x0000FF00) >> 8),
				byte((crc & 0x000000FF)),
			},
			body: body,
		}
	*/
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
