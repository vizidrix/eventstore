package eventstore_test

import (
	//"encoding/binary"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	"hash/crc32"
	"log"
	"testing"
	"time"
)

func ignore() { log.Println("") }

// go test -v github.com/vizidrix/eventstore -bench . -cpuprofile /go/prof.out -cpu 1,2,4,8

type VisitorLogged struct {
	UnixNSTimeStamp int64
	IPv4Address     int64
	IPv6Header      int64
	IPv6Address     int64
	Referrer        string
}

/*
func (event *VisitorLogged) ToBinary() ([]byte, error) {
	buffer := make([]byte, 4+4+4+4+len(event.Referrer))
	index := 0
	buffer[index] = binary.Write(buffer, binary.BigEndian, event.UnixNSTimeStamp)

	index++
	for char := range event.Value {
		buffer[index] = byte(char)
		index++
	}
	return buffer, nil
}
*/

func Test_Should_produce_correct_CRC_for_event_entry(t *testing.T) {
	// Arrange
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}

	// Act
	entry := goes.NewEventStoreEntryFrom(5, data)

	// Assert
	AreNotEqual(t, int16(0), entry.CRC(), "CRC should have been calculated")
	crc := crc32.Checksum(entry.Data(), crc32.MakeTable(crc32.Castagnoli))
	AreEqual(t, crc, entry.CRC(), "CRC should be correct")
}

func Test_Should_try_to_connect_to_MemoryEventstore_with_correct_path(t *testing.T) {
	// Arrange
	//path := "/eventstore/"
	path := "mem://"

	// Act
	eventStore, err := goes.Connect(path)

	// Assert
	if err != nil {
		log.Printf("Connect failed: %s", err)
		t.Fail()
		return
	}
	switch es := (eventStore).(type) {
	case *goes.MemoryEventStore:
		{
		}
	default:
		{
			log.Printf("Wrong event store type created: %s", es)
			t.Fail()
		}
	}
}

/*
func Test_Should_try_to_connect_to_FileSystemEventStore_with_correct_path(t *testing.T) {
	// Arrange
	path := "fs://eventstore/"

	// Act
	eventStore, err := goes.Connect(path)

	// Assert
	if err != nil {
		log.Printf("Connect failed: %s", err)
		t.Fail()
		return
	}
	switch es := (eventStore).(type) {
	case *goes.FileSystemEventStore:
		{
		}
	default:
		{
			log.Printf("Wrong event store type created: %s", es)
			t.Fail()
		}
	}
}
*/
func Test_Should_try_to_connect_to_FragmentFileSystemEventStore_with_correct_path(t *testing.T) {
	// Arrange
	path := "ffs://eventstore/"

	// Act
	eventStore, err := goes.Connect(path)

	// Assert
	if err != nil {
		log.Printf("Connect failed: %s", err)
		t.Fail()
		return
	}
	switch es := (eventStore).(type) {
	case *goes.FragmentFileSystemEventStore:
		{
		}
	default:
		{
			log.Printf("Wrong event store type created: %s", es)
			t.Fail()
		}
	}
}

func Test_Should_return_error_if_connstring_is_invalid(t *testing.T) {
	// Arrange
	paths := []string{
		"invalid://stuff",
		"a://bleh",
		"fs2://real/",
		"amem://stuff",
		"mem//stuff",
		"fs:/stuff",
		"£¢://broke",
		"http://stuff",
		"tcp://stuff",
		"fs",
		"mem",
		"http",
		"tcp",
	}

	// Act
	for _, path := range paths {
		if _, err := goes.Connect(path);
		// Assert
		err == nil {
			log.Printf("Invalid path should have raised an err: %s", path)
			t.Fail()
		}
	}
}

func EventStoreSync_Should_return_empty_slice_for_new_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "type")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 1)

	// Act
	eventStore.LoadAll(uri, events)

	// Assert
	select {
	case event := <-events:
		{
			log.Printf("Shouldn't have received any events: %s", event)
			t.Fail()
		}
	case <-time.After(1 * time.Microsecond):
		{
			// Shouldn't have anything on events channel
		}
	}
}

func EventStoreSync_Should_return_single_matching_event_for_existing_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "type")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 1)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry := goes.NewEventStoreEntry(10, 1, 1, data)
	eventStore.Append(uri, entry)

	// Act
	eventStore.LoadAll(uri, events)

	actual := <-events

	// Assert
	AreEqual(t, uint16(10), actual.Length(), "Length should have been int32 10")
	AreEqual(t, uint16(1), actual.EventType(), "EvenType should have been set")
	AreEqual(t, uint32(1), actual.CRC(), "CRC should have been calculated")
	AreAllEqual(t, data, actual.Data(), "Data should have been set")
}

func EventStoreSync_Should_return_middle_events_for_version_range(t *testing.T, connString string) {
	// Arrange
	eventstore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 5)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	for index := 0; index < 5; index++ {
		entry := goes.NewEventStoreEntry(10, uint16(index), uint32(index), data)
		eventstore.Append(uri, entry)
	}

	// Act
	eventstore.LoadIndexRange(uri, events, 2, 3)

	// Assert
	for index := 2; index < 4; index++ {
		//log.Printf("Index: %d", index)
		select {
		case event := <-events:
			{
				//log.Printf("Event received: % x", event)

				AreEqual(t, uint16(10), event.Length(), "Length should have been int32 10")
				AreEqual(t, uint16(index), event.EventType(), "EvenType should have been set")
				AreEqual(t, uint32(index), event.CRC(), "CRC should have been calculated")
				AreAllEqual(t, data, event.Data(), "Data should have been set")
			}
		case <-time.After(1 * time.Millisecond):
			{
				log.Printf("Shouldn't have timed out")
				t.Fail()
				return
			}
		}
	}
}

func EventStoreSync_Should_return_two_matching_events_for_existing_ids(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 2)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(10, 0, 0, data)
	entry2 := goes.NewEventStoreEntry(10, 1, 1, data)

	eventStore.Append(uri, entry1)
	eventStore.Append(uri, entry2)

	// Act
	eventStore.LoadAll(uri, events)

	// Assert
	for index := 0; index < 2; index++ {
		select {
		case event := <-events:
			{
				AreEqual(t, uint16(10), event.Length(), "Length should have been int32 10")
				AreEqual(t, uint16(index), event.EventType(), "EvenType should have been set")
				AreEqual(t, uint32(index), event.CRC(), "CRC should have been calculated")
				AreAllEqual(t, data, event.Data(), "Data should have been set")
			}
		case <-time.After(1 * time.Millisecond):
			{
				log.Printf("Shouldn't have timed out")
				t.Fail()
				return
			}
		}
	}
}

func EventStoreSync_Should_not_panic_when_range_is_too_long(t *testing.T, connString string) {
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, goes.MAX_EVENT_SIZE)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE, 1, 1, data)
	events := make(chan *goes.EventStoreEntry, 1)
	uri := kind.ToAggregateUri(1)
	eventStore.Append(uri, entry1)
	eventStore.LoadIndexRange(uri, events, 0, 4)
}

func EventStoreSync_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t *testing.T, connString string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, goes.MAX_EVENT_SIZE+1)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE+1, 1, 1, data)
	uri := kind.ToAggregateUri(1)
	eventStore.Append(uri, entry1)
}

func EventStoreSync_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t *testing.T, connString string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, 3082) // <- set to less than length
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(3083, 1, 1, data) // <- invalid length!
	uri := kind.ToAggregateUri(1)
	eventStore.Append(uri, entry1)
}

func EventStoreSync_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t *testing.T, connString string) {
	count := 2

	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	entry1 := Get_EventStoreEntry(10)

	roundComplete := make(chan struct{})

	for i := 0; i < count; i++ {
		index := i
		go func() {
			uri := kind.ToAggregateUri(int64(index))
			eventStore.Append(uri, entry1)
			events := make(chan *goes.EventStoreEntry, 2)
			eventStore.LoadAll(uri, events)
			roundComplete <- struct{}{}
		}()
	}
	for i := 0; i < count; i++ {
		select {
		case <-roundComplete:
			{
			}
		case <-time.After(10 * time.Millisecond):
			{
				log.Printf("Shouldn't have timed out")
				t.Fail()
			}
		}
	}
}

func EventStoreAsync_Should_return_empty_slice_for_new_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "type")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 1)

	// Act
	completed, errored := eventStore.LoadAllAsync(uri, events)

	// Assert
	select {
	case <-completed:
		{
			// Wait for complete to trigger
		}
	case err := <-errored:
		{
			log.Printf("Shouldn't have received any errors: %s", err)
			t.Fail()
			return
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
			return
		}
	}

	select {
	case event := <-events:
		{
			log.Printf("Shouldn't have received any events: %s", event)
			t.Fail()
		}
	case <-time.After(1 * time.Microsecond):
		{
			// Shouldn't have anything on events channel
		}
	}
}

func EventStoreAsync_Should_return_single_matching_event_for_existing_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "type")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 1)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry := goes.NewEventStoreEntry(10, 1, 1, data)
	appendCompleted, _ := eventStore.AppendAsync(uri, entry)
	<-appendCompleted

	// Act
	completed, errored := eventStore.LoadAllAsync(uri, events)

	select {
	case <-completed:
		{
		}
	case err := <-errored:
		{
			log.Printf("Shouldn't have received any errors: %s", err)
			t.Fail()
			return
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
			return
		}
	}
	actual := <-events

	// Assert
	AreEqual(t, uint16(10), actual.Length(), "Length should have been int32 10")
	AreEqual(t, uint16(1), actual.EventType(), "EvenType should have been set")
	AreEqual(t, uint32(1), actual.CRC(), "CRC should have been calculated")
	AreAllEqual(t, data, actual.Data(), "Data should have been set")
}

func EventStoreAsync_Should_return_middle_events_for_version_range(t *testing.T, connString string) {
	// Arrange
	eventstore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 5)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	for index := 0; index < 5; index++ {
		entry := goes.NewEventStoreEntry(10, uint16(index), uint32(index), data)
		appendCompleted, _ := eventstore.AppendAsync(uri, entry)
		<-appendCompleted
	}

	// Act
	completed, errored := eventstore.LoadIndexRangeAsync(uri, events, 2, 3)

	select {
	case <-completed:
		{
		}
	case err := <-errored:
		{
			log.Printf("Shouldn't have received any errors: %s", err)
			t.Fail()
			return
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
			return
		}
	}

	// Assert
	for index := 2; index < 4; index++ {
		//log.Printf("Index: %d", index)
		select {
		case event := <-events:
			{
				//log.Printf("Event received: % x", event)

				AreEqual(t, uint16(10), event.Length(), "Length should have been int32 10")
				AreEqual(t, uint16(index), event.EventType(), "EvenType should have been set")
				AreEqual(t, uint32(index), event.CRC(), "CRC should have been calculated")
				AreAllEqual(t, data, event.Data(), "Data should have been set")
			}
		case <-time.After(1 * time.Millisecond):
			{
				log.Printf("Shouldn't have timed out")
				t.Fail()
				return
			}
		}
	}
}

func EventStoreAsync_Should_return_two_matching_events_for_existing_ids(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 2)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(10, 0, 0, data)
	entry2 := goes.NewEventStoreEntry(10, 1, 1, data)

	append1Completed, _ := eventStore.AppendAsync(uri, entry1)
	append2Completed, _ := eventStore.AppendAsync(uri, entry2)
	<-append1Completed
	<-append2Completed

	// Act
	completed, errored := eventStore.LoadAllAsync(uri, events)

	select {
	case <-completed:
		{
		}
	case err := <-errored:
		{
			log.Printf("Shouldn't have received any errors: %s", err)
			t.Fail()
			return
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
			return
		}
	}

	// Assert
	for index := 0; index < 2; index++ {
		select {
		case event := <-events:
			{
				AreEqual(t, uint16(10), event.Length(), "Length should have been int32 10")
				AreEqual(t, uint16(index), event.EventType(), "EvenType should have been set")
				AreEqual(t, uint32(index), event.CRC(), "CRC should have been calculated")
				AreAllEqual(t, data, event.Data(), "Data should have been set")
			}
		case <-time.After(1 * time.Millisecond):
			{
				log.Printf("Shouldn't have timed out")
				t.Fail()
				return
			}
		}
	}
}

func EventStoreAsync_Should_not_panic_when_range_is_too_long(t *testing.T, connString string) {
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, goes.MAX_EVENT_SIZE)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE, 1, 1, data)
	events := make(chan *goes.EventStoreEntry, 1)
	uri := kind.ToAggregateUri(1)
	complete, _ := eventStore.AppendAsync(uri, entry1)
	<-complete
	readComplete, _ := eventStore.LoadIndexRangeAsync(uri, events, 0, 4)
	<-readComplete
}

func EventStoreAsync_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t *testing.T, connString string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, goes.MAX_EVENT_SIZE+1)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE+1, 1, 1, data)
	uri := kind.ToAggregateUri(1)
	complete, _ := eventStore.AppendAsync(uri, entry1)
	<-complete
}

func EventStoreAsync_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t *testing.T, connString string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, 3082) // <- set to less than length
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(3083, 1, 1, data) // <- invalid length!
	uri := kind.ToAggregateUri(1)
	complete, _ := eventStore.AppendAsync(uri, entry1)
	<-complete
}

func EventStoreAsync_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t *testing.T, connString string) {
	count := 2

	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	entry1 := Get_EventStoreEntry(10)

	roundComplete := make(chan struct{})

	for i := 0; i < count; i++ {
		index := i
		go func() {
			uri := kind.ToAggregateUri(int64(index))
			appendComplete, _ := eventStore.AppendAsync(uri, entry1)
			<-appendComplete
			events := make(chan *goes.EventStoreEntry, 2)
			readComplete, _ := eventStore.LoadAllAsync(uri, events)
			<-events
			<-readComplete
			roundComplete <- struct{}{}
		}()
	}
	for i := 0; i < count; i++ {
		select {
		case <-roundComplete:
			{
			}
		case <-time.After(10 * time.Millisecond):
			{
				log.Printf("Shouldn't have timed out")
				t.Fail()
			}
		}
	}
}
