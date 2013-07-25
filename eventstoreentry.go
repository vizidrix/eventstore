package eventstore

import (
//"bytes"
//"encoding/binary"
//"errors"
//"fmt"
//"log"
//"time"
)

const (
	header_size = 8
)

type EventStoreEntry struct {
	event []byte
}

func NewEventStoreEntryFrom(eventType byte, data []byte) *EventStoreEntry {
	crc := int32(len(data))
	return NewEventStoreEntry(int32(len(data)), eventType, crc, data)
}

func NewEventStoreEntry(length int32, eventType byte, crc int32, data []byte) *EventStoreEntry {
	event := make([]byte, header_size+length)
	event[0] = byte(length & 0x000F)
	event[1] = byte(length & 0x00F0)
	event[2] = byte(length & 0x0F00)
	event[3] = eventType
	event[4] = byte(crc & 0xF000)
	event[5] = byte(crc & 0x0F00)
	event[6] = byte(crc & 0x00F0)
	event[7] = byte(crc & 0x000F)
	copy(event[8:], data)
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

// Length of the trailing data block
func (entry *EventStoreEntry) Length() int32 {
	//result := entry.data[2]<<2 & entry.data[1]<<1 & entry.data[0]
	return int32(entry.event[2]) << 2 & int32(entry.event[1]) << 1 & int32(entry.event[0])
	//return entry.length_type >> 4
}

// Identifier used by serializer to do it's magic
func (entry *EventStoreEntry) EventType() byte {
	return entry.event[3]
}

// Checksum of the trailing data block
func (entry *EventStoreEntry) CRC() int32 {
	return int32(entry.event[4]) << 3 & int32(entry.event[5]) << 2 & int32(entry.event[6]) << 1 & int32(entry.event[7])
}

// Byte slice to hold the binary data
func (entry *EventStoreEntry) Data() []byte {
	return entry.event[header_size:]
}

/*
type EventStoreEntry struct {
	length_type int32  // 3 bytes for len, 1 for type
	crc         int32  // 4 bytes
	data        []byte // Max 4096 (3 bytes) bytes less header (17 bytes)
}

func NewEventStoreEntryFrom(eventType byte, data []byte) (*EventStoreEntry, error) {
	if eventType <= 0 || 255 <= eventType {
		message := fmt.Sprintf("Invalid event type: %d", eventType)
		return nil, errors.New(message)
	}
	//if len(data) == 0 {
	//	return nil, errors.New("Empty data slice")
	//}
	return NewEventStoreEntry(
		int32(len(data)),
		eventType,
		0,
		data)
}

func NewEventStoreEntry(length int32, eventType byte, crc int32, data []byte) (*EventStoreEntry, error) {
	return &EventStoreEntry{
		length_type: length<<4&0xFFF0 | int32(eventType),
		crc:         crc,
		data:        data,
	}, nil
}

func WriteBinaryBatch(buffer *bytes.Buffer, data ...interface{}) error {
	for _, item := range data {
		if err := binary.Write(buffer, binary.LittleEndian, item); err != nil {
			return err
		}
	}
	return nil
}

func (entry *EventStoreEntry) ToBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := WriteBinaryBatch(buffer,
		entry.length_type,
		entry.crc,
		entry.data)

	if err != nil {
		log.Printf("Error writing to buffer: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

const (
	header_size = 8
)

func FromBinary(data []byte) (*EventStoreEntry, error) {
	buffer := bytes.NewBuffer(data)
	entry := EventStoreEntry{}
	if err := binary.Read(buffer, binary.LittleEndian, &entry.length_type); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &entry.crc); err != nil {
		return nil, err
	}
	entry.data = data[header_size : header_size+entry.Length()]
	data = data[header_size+entry.Length():]

	return &entry, nil
}

// Length of the trailing data block
func (entry *EventStoreEntry) Length() int32 {
	return entry.length_type >> 4
}

// Identifier used by serializer to do it's magic
func (entry *EventStoreEntry) EventType() byte {
	return byte(entry.length_type & 0x000F)
}

// Checksum of the trailing data block
func (entry *EventStoreEntry) CRC() int32 {
	return entry.crc
}

// Byte slice to hold the binary data
func (entry *EventStoreEntry) Data() []byte {
	return entry.data
}
*/
