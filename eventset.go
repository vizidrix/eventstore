package eventstore

import (
	"errors"
	//"bufio"
	//"bytes"
	"fmt"
	//"io"
	"log"
	//"math/rand"
	//"os"
	//"os"
	//"encoding/binary"
	//"strings"

	//"time"
	//"runtime"
)

func eventset_ignore() {
	log.Printf(fmt.Sprintf(""))
	//runtime.Sto
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
	headerData []byte // Raw header data
	eventData  []byte // Raw event data

	headPos int // Index of first event in set
	tailPos int // Index of last event in set
	//offsets []int // Precalculated map of event locations

	headers []Header // Typed view of header
	events  [][]byte
}

const (
	HEADER_SLICE_SIZE = 64   // Room for 8 appends before expand
	DATA_SLICE_SIZE   = 4096 //* 2048
	MAX_EVENT_SIZE    = int(MaxUint16)
)

func NewEmptyEventSet() *EventSet {
	headerData := make([]byte, 0, HEADER_SLICE_SIZE)
	return &EventSet{
		// Initialize primary raw data store
		headerData: headerData,
		eventData:  make([]byte, 0, DATA_SLICE_SIZE),
		headPos:    0,
		tailPos:    0,
		headers:    UnsafeCastBytesToHeader(headerData),
		events:     make([][]byte, 0, HEADER_SLICE_SIZE),
	}
}

func (set *EventSet) CheckSum() error {
	//headers := UnsafeCastBytesToHeader(set.headers)
	//position := 0
	for index, _ := range set.headers {
		//crc := MakeCRC(set.data[position : position+int(header.length)])
		crc := MakeCRC(set.events[index])
		if crc != set.headers[index].crc {
			return errors.New("Data appears corrupted")
		}
		//position += int(set.headers[index].length)
	}
	return nil
}

func (set *EventSet) Count() int {
	return len(set.headers) // >> 3
}

func (set *EventSet) LengthAt(index int) uint16 {
	return set.headers[index].length
}

func (set *EventSet) EventTypeAt(index int) uint16 {
	return set.headers[index].eventType
}

func (set *EventSet) CheckSumAt(index int) uint32 {
	return set.headers[index].crc
}

func (set *EventSet) DataAt(index int) []byte {
	return set.events[index]
}

/*
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
*/

func (set *EventSet) Get() (*EventSet, error) {
	return set.GetSlice(0, MaxInt)
}

func (set *EventSet) GetSlice(startIndex int, endIndex int) (*EventSet, error) {
	headerLength := len(set.headers)
	// Validate inputs
	if startIndex < 0 || endIndex < 0 || startIndex >= endIndex || startIndex > headerLength {
		return nil, errors.New("Either start or end index is out of range")
	}
	// No data so just return empty slice
	if headerLength == 0 {
		return NewEmptyEventSet(), nil
	}
	// End out of range so move it back
	if endIndex >= headerLength {
		endIndex = headerLength
	}

	headerData := set.headerData[startIndex<<3 : endIndex<<3]
	events := set.events[startIndex:endIndex]
	dataLength := cap(set.eventData)
	lowerBound := dataLength - cap(events[0])
	upperBound := dataLength - cap(events[len(events)-1])
	eventData := set.eventData[lowerBound:upperBound]
	headers := UnsafeCastBytesToHeader(headerData)

	return &EventSet{
		headerData: headerData,
		eventData:  eventData,
		headPos:    set.headPos - startIndex,
		tailPos:    set.tailPos - endIndex,
		headers:    headers,
		events:     events,
	}, nil

	//return eventSet, nil
}

func (set *EventSet) Put(newEvents ...Event) (*EventSet, error) {
	prevHeaderSize := len(set.headerData)
	prevHeaderCap := cap(set.headerData)
	prevEventSize := len(set.eventData)
	prevEventCap := cap(set.eventData)

	newCount := len(newEvents)

	reqHeaderSize := prevHeaderSize + (newCount << 3)
	prevHeaderCount := prevHeaderSize >> 3
	newHeaderSize := newCount << 3

	// Placeholders for the new immutable set
	var headerData []byte
	var eventData []byte
	var headers []Header
	var events [][]byte

	// Decide if we need expand or grow for headers
	if reqHeaderSize < prevHeaderCap {
		headerData = set.headerData[0 : prevHeaderSize+newHeaderSize]
		// Event overlay follows the same growth rules as headers
		events = set.events[0 : prevHeaderSize+newHeaderSize]
	} else { // Magic expando sauce needed for header (growth algo for fixed size enries)
		headerData = make([]byte, reqHeaderSize, reqHeaderSize<<1)
		copy(headerData, set.headerData)
		// Should only have to copy over pointers to the slices so it will be fast
		events = make([][]byte, reqHeaderSize, reqHeaderSize<<1)
		copy(events, set.events)
	}
	// Build the header overlay now that data has been resized
	headers = UnsafeCastBytesToHeader(headerData)

	newEventSize := 0
	// Copy info from events into the new header structures
	for i := range newEvents {
		eventSize := len(newEvents[i].Data)
		if eventSize > MAX_EVENT_SIZE {
			return nil, errors.New("Event data too large")
		}
		newEventSize += eventSize

		headers[prevHeaderCount+i].length = uint16(eventSize)
		headers[prevHeaderCount+i].eventType = newEvents[i].EventType
		headers[prevHeaderCount+i].crc = MakeCRC(newEvents[i].Data)
	}

	// Now we can figure out the size needed to store event data
	reqEventSize := prevEventSize + newEventSize

	// Decide if we need expand or grow for data
	if reqEventSize < prevEventCap {
		eventData = set.eventData[0:reqEventSize]
	} else { // Magic expando sauce needed (growth algo for variable sized entries)
		// Ensures that the cap is 16 byte alligned... 0x3F would be 8 byte alligned
		// Account for at least 2 similarly sized, 16 byte alligned adds
		reqEventCap := (reqEventSize | 0x7F) + ((newEventSize << 1) | 0x7F)
		eventData = make([]byte, reqEventSize, reqEventCap)
		copy(eventData, set.eventData)
	}

	// Now we can copy over the event data and build the event overlay
	for i := range newEvents {
		eventSize := len(newEvents[i].Data)
		if eventSize == 0 {
			continue
		}
		//UnsafeCopyBytes(data[prevDataSize:], events[i].Data)
		copy(eventData[prevEventSize:], newEvents[i].Data)
		// Build the overlay for this event
		events[i] = eventData[prevEventSize : prevEventSize+eventSize]
		// Track the moving position
		prevEventSize += eventSize
	}

	return &EventSet{
		headerData: headerData,
		eventData:  eventData,
		headPos:    set.headPos,
		tailPos:    set.tailPos + newCount,
		headers:    headers,
		events:     events,
	}, nil
}
