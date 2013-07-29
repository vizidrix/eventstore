package eventstore_test

import (
	//"encoding/binary"
	//"bytes"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	"hash/crc32"
	"log"
	//"math/rand"
	"testing"
	"time"
)

func ignore_memoryeventstore_test() { log.Println("") }

func Test_Should_return_empty_slice_for_new_id(t *testing.T) {
	// Arrange
	eventStore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "type")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 1)

	// Act
	completed, errored := eventStore.LoadAll(uri, events)

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

func Test_Should_return_single_matching_event_for_existing_id(t *testing.T) {
	// Arrange
	eventStore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "type")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 1)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry := goes.NewEventStoreEntry(10, 1, 1, data)
	appendCompleted, _ := eventStore.Append(uri, entry)
	<-appendCompleted

	// Act
	completed, errored := eventStore.LoadAll(uri, events)

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

func Test_Should_return_middle_events_for_version_range(t *testing.T) {
	// Arrange
	eventstore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "kind")
	uri := kind.ToAggregateUri(1)
	events := make(chan *goes.EventStoreEntry, 5)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	for index := 0; index < 5; index++ {
		entry := goes.NewEventStoreEntry(10, uint16(index), uint32(index), data)
		appendCompleted, _ := eventstore.Append(uri, entry)
		<-appendCompleted
	}

	// Act
	completed, errored := eventstore.LoadIndexRange(uri, events, 2, 3)

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

func Test_Should_return_two_matching_events_for_existing_ids(t *testing.T) {
	// Arrange
	eventStore, _ := goes.Connect("mem://")
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
	completed, errored := eventStore.LoadAll(uri, events)

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

func Test_Should_not_panic_when_range_is_too_long(t *testing.T) {
	eventStore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, goes.MAX_EVENT_SIZE)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE, 1, 1, data)
	events := make(chan *goes.EventStoreEntry, 1)
	uri := kind.ToAggregateUri(1)
	complete, _ := eventStore.Append(uri, entry1)
	<-complete
	readComplete, _ := eventStore.LoadIndexRange(uri, events, 0, 4)
	<-readComplete
}

func Test_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, goes.MAX_EVENT_SIZE+1)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE+1, 1, 1, data)
	uri := kind.ToAggregateUri(1)
	complete, _ := eventStore.Append(uri, entry1)
	<-complete
}

func Test_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, 3082) // <- set to less than length
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(3083, 1, 1, data) // <- invalid length!
	uri := kind.ToAggregateUri(1)
	complete, _ := eventStore.Append(uri, entry1)
	<-complete
}

func Test_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t *testing.T) {
	count := 2

	eventStore, _ := goes.Connect("mem://")
	kind := goes.NewAggregateKind("namespace", "kind")
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(10, 1, 1, data)

	roundComplete := make(chan struct{})

	for i := 0; i < count; i++ {
		index := i
		go func() {
			uri := kind.ToAggregateUri(int64(index))
			appendComplete, _ := eventStore.Append(uri, entry1)
			<-appendComplete
			events := make(chan *goes.EventStoreEntry)
			readComplete, _ := eventStore.LoadAll(uri, events)
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

func Benchmark_Create_Serialize_DeSerialize_EventStoreEntry_10bytePayload(b *testing.B) {
	b.StopTimer()
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		entry := goes.NewEventStoreEntryFrom(1, data)
		temp := entry.ToBinary()
		goes.FromBinary(temp)
	}
}

func Benchmark_Create_Serialize_DeSerialize_EventStoreEntry_4084bytePayload(b *testing.B) {
	b.StopTimer()
	data := make([]byte, 4084)
	for index, _ := range data {
		data[index] = byte(index)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		entry := goes.NewEventStoreEntryFrom(1, data)
		temp := entry.ToBinary()
		goes.FromBinary(temp)
	}
}

func Benchmark_KindUri_to_AggregateRootUri(b *testing.B) {
	kind := goes.NewAggregateKind("namespace", "kind")

	for i := 0; i < b.N; i++ {
		kind.ToAggregateUri(10000)
	}
}

func Benchmark_MemoryEventStore_AppendOnly_10bytePayload(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 10, 1)
}

func Benchmark_MemoryEventStore_AppendOnly_4084bytePayload(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 4084, 1)
}

func Benchmark_MemoryEventStore_ReadOnly_10bytePayload(b *testing.B) {
	Run_ReadOnly(b, "mem://", "namespace", "kind", 10, 1)
}

func Benchmark_MemoryEventStore_ReadOnly_4084bytePayload(b *testing.B) {
	Run_ReadOnly(b, "mem://", "namespace", "kind", 4084, 1)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_10bytePayload(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 10, 1)
}

// 256 4k events / 1mb
// 4124 ns/ 1 op = 2424 op / ms = 2,424,000 op / s = 9468 mb / s
func Benchmark_MemoryEventStore_AppendAndReadAll_4084bytePayload(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 4084, 1)
}

func Benchmark_MemoryEventStore_AppendOnly_20_10bytePayloads(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 10, 20)
}

func Benchmark_MemoryEventStore_AppendOnly_20_4084bytePayloads(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 4084, 20)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_20_10bytePayloads(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 10, 20)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_20_4084bytePayloads(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 4084, 20)
}
