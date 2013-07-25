package eventstore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"
)

type EventStoreEntry struct {
	length        int32  // 4 bytes
	crc           int32  // 4 bytes
	unixTimeStamp int64  // 8 bytes
	eventType     byte   // 1 byte
	data          []byte // Max 4096 (3 bytes) bytes less header (17 bytes)
}

func NewEventStoreEntryFrom(eventType byte, data []byte) (*EventStoreEntry, error) {
	if eventType <= 0 {
		message := fmt.Sprintf("Invalid event type: %d", eventType)
		return nil, errors.New(message)
	}
	if len(data) == 0 {
		return nil, errors.New("Empty data slice")
	}
	return NewEventStoreEntry(
		int32(len(data)),
		0,
		time.Now().UnixNano(),
		eventType,
		data)
	/*
		return &EventStoreEntry{
			length:        int32(len(data)),
			crc:           0,
			unixTimeStamp: time.Now().UnixNano(),
			eventType:     eventType,
			data:          data,
		}, nil
	*/
}

func NewEventStoreEntry(length int32, crc int32, unixTimeStamp int64, eventType byte, data []byte) (*EventStoreEntry, error) {
	return &EventStoreEntry{
		length:        length,
		crc:           crc,
		unixTimeStamp: unixTimeStamp,
		eventType:     eventType,
		data:          data,
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
		entry.length,
		entry.crc,
		entry.unixTimeStamp,
		entry.eventType,
		entry.data)

	if err != nil {
		log.Printf("Error writing to buffer: %s", err)
		return nil, err
	}

	return buffer.Bytes(), nil
}

func FromBinary(data []byte) (*EventStoreEntry, error) {
	buffer := bytes.NewBuffer(data)
	entry := EventStoreEntry{}
	if err := binary.Read(buffer, binary.LittleEndian, &entry.length); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &entry.crc); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &entry.unixTimeStamp); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.LittleEndian, &entry.eventType); err != nil {
		return nil, err
	}
	entry.data = data[17 : 17+entry.length]
	data = data[17+entry.length:]

	return &entry, nil
}

// Length of the trailing data block
func (entry *EventStoreEntry) Length() int32 {
	return entry.length
}

// Checksum of the trailing data block
func (entry *EventStoreEntry) CRC() int32 {
	return entry.crc
}

// Timestamp when entry was written to the store
func (entry *EventStoreEntry) UnixTimeStamp() int64 {
	return entry.unixTimeStamp
}

// Identifier used by serializer to do it's magic
func (entry *EventStoreEntry) EventType() byte {
	return entry.eventType
}

// Byte slice to hold the binary data
func (entry *EventStoreEntry) Data() []byte {
	return entry.data
}
