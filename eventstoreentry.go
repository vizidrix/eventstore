package eventstore

import (
	//"errors"
	//"fmt"
	"hash/crc32"
	"log"

//"time"
)

func ignore_eventstore() { log.Printf("") }

const (
	header_size = 9
)

type EventStoreEntry struct {
	event []byte
}

func NewEventStoreEntryFrom(eventType uint16, data []byte) *EventStoreEntry {
	crc := crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
	return NewEventStoreEntry(int32(len(data)), eventType, crc, data)
}

func NewEventStoreEntry(length int32, eventType uint16, crc uint32, data []byte) *EventStoreEntry {
	event := make([]byte, header_size+length)
	event[0] = byte(length & 0x0F00)
	event[1] = byte(length & 0x00F0)
	event[2] = byte(length & 0x000F)
	event[3] = byte((eventType & 0xFF00) >> 8)
	event[4] = byte((eventType & 0x00FF))
	event[5] = byte((crc & 0xFF000000) >> 24)
	event[6] = byte((crc & 0x00FF0000) >> 16)
	event[7] = byte((crc & 0x0000FF00) >> 8)
	event[8] = byte((crc & 0x000000FF))
	copy(event[9:], data[:length])
	return &EventStoreEntry{
		event: event,
	}
}

func (entry *EventStoreEntry) ToBinary() []byte {
	return entry.event
}

func FromBinary(event []byte) *EventStoreEntry {
	return &EventStoreEntry{
		event: event,
	}
}

func Int24(data []byte) int32 {
	return int32(data[0])<<2 | int32(data[1])<<1 | int32(data[2])
}

// Length of the trailing data block
func (entry *EventStoreEntry) Length() int32 {
	return Int24(entry.event[0:3])
	//return int32(entry.event[0])<<2 | int32(entry.event[1])<<1 | int32(entry.event[2])
}

// Identifier used by serializer to do it's magic
func (entry *EventStoreEntry) EventType() uint16 {
	return uint16(entry.event[3])<<8 | uint16(entry.event[4])
}

// Checksum of the trailing data block
func (entry *EventStoreEntry) CRC() uint32 {
	return uint32(entry.event[5])<<24 | uint32(entry.event[6])<<16 | uint32(entry.event[7])<<8 | uint32(entry.event[8])
}

// Byte slice to hold the binary data
func (entry *EventStoreEntry) Data() []byte {
	return entry.event[header_size : header_size+entry.Length()]
}
