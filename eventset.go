package eventstore

import (
	"errors"
	"fmt"
	"log"
)

func eventset_ignore() {
	log.Printf(fmt.Sprintf(""))
	//runtime.Sto
}

/*

Header Index - Per Aggregate with Header Record (8 byte) per event
[ 2 byte		| 2 byte		| 4 byte 	]
[ Length		| EventType		| CRC		]

Max Len: 		65535
Max Event Type: 65535

Stride data in 64bit chunks by Blocks() count?

Actual data is Length() long padded to fit chunks?

*/

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

	headers []Header // Typed view of header
	events  [][]byte // Map of index to event data
}

const (
	HEADER_SLICE_SIZE = 64 // Room for 8 appends before expand
	// ** Huge tuning variable
	// Smaller # will result in general speed improvement for items below the threshold size
	// Larger # will result in faster first append for larger events
	// Set this to the largest average event size?
	DATA_SLICE_SIZE = 1024
	MAX_EVENT_SIZE  = int(MaxUint16)
)

func NewEmptyEventSet() *EventSet {
	headerData := new([HEADER_SLICE_SIZE]byte)[0:0]
	return &EventSet{
		// Initialize primary raw data store
		headerData: headerData,
		eventData:  new([DATA_SLICE_SIZE]byte)[0:0],
		headPos:    0,
		tailPos:    0,
		headers:    UnsafeCastBytesToHeader(headerData),
		events:     new([HEADER_SLICE_SIZE][]byte)[0:0],
	}
}

func (set *EventSet) CheckSum() error {
	for index, _ := range set.headers {
		crc := MakeCRC(set.events[index])
		if crc != set.headers[index].crc {
			return errors.New("Data appears corrupted")
		}
	}
	return nil
}

func (set *EventSet) Count() int {
	return len(set.headers)
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

	// If the slice size is the whole set then just return it
	if startIndex == 0 && endIndex == headerLength {
		return set, nil
	}

	dataLength := cap(set.eventData)
	lowerBound := dataLength - cap(set.events[startIndex])
	upperBound := dataLength - cap(set.events[endIndex])
	return &EventSet{
		headerData: set.headerData[startIndex<<3 : endIndex<<3],
		eventData:  set.eventData[lowerBound:upperBound],
		headPos:    set.headPos - startIndex,
		tailPos:    set.tailPos - endIndex,
		headers:    set.headers[startIndex:endIndex],
		events:     set.events[startIndex:endIndex],
	}, nil
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
	if reqHeaderSize <= prevHeaderCap {
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
	if reqEventSize <= prevEventCap {
		eventData = set.eventData[0:reqEventSize]
	} else { // Magic expando sauce needed (growth algo for variable sized entries)
		// Ensures that the cap is byte alligned... 0x3F would be 8 byte alligned
		desiredBuffer := ((reqEventSize * EVENT_GROWTH_MULTIPLIER) >> 7) | 0x3F
		if MAX_EVENT_BUFFER >= 0 && MAX_EVENT_BUFFER < desiredBuffer {
			desiredBuffer = MAX_EVENT_BUFFER
		}
		if MIN_EVENT_BUFFER > desiredBuffer {
			desiredBuffer = MIN_EVENT_BUFFER
		}
		eventData = make([]byte, reqEventSize, reqEventSize+desiredBuffer)
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

const (
	MIN_EVENT_BUFFER        = 0   //64
	MAX_EVENT_BUFFER        = -1  // Limits the growth of the event buffers. -1 turns it off.
	EVENT_GROWTH_MULTIPLIER = 256 // (BUFFER * N) / 128 so 160 = 125%, 192 = 150%
)
