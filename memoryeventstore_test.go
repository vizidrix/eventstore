package eventstore_test

import (
	//"encoding/binary"
	//"bytes"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	"hash/crc32"
	"log"
	"math/rand"
	"testing"
	"time"
)

func ignore_memoryeventstore_test() { log.Println("") }

func Test_Should_return_empty_slice_for_new_id(t *testing.T) {
	// Arrange
	eventStore, _ := goes.Connect("mem://")
	uri := goes.NewAggregateRootUri("namespace", "type", 1)
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
	uri := goes.NewAggregateRootUri("namespace", "type", 1)
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
	AreEqual(t, int32(10), actual.Length(), "Length should have been int32 10")
	AreEqual(t, uint16(1), actual.EventType(), "EvenType should have been set")
	AreEqual(t, uint32(1), actual.CRC(), "CRC should have been calculated")
	AreAllEqual(t, data, actual.Data(), "Data should have been set")
}

func Test_Should_return_middle_events_for_version_range(t *testing.T) {
	// Arrange
	eventstore, _ := goes.Connect("mem://")
	uri := goes.NewAggregateRootUri("namespace", "kind", 1)
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
	case <-time.After(1000 * time.Millisecond):
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

				AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
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
	uri := goes.NewAggregateRootUri("namespace", "kind", 1)
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

				AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
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
	AreNotEqual(t, int32(0), entry.CRC(), "CRC should have been calculated")
	crc := crc32.Checksum(entry.Data(), crc32.MakeTable(crc32.Castagnoli))
	AreEqual(t, crc, entry.CRC(), "CRC should be correct")
}

func Benchmark_Create_Serialize_DeSerialize_EventStoreEntry_10bytePayload(b *testing.B) {
	b.StopTimer()
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		entry := goes.NewEventStoreEntry(10, 1, 1, data)
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
		entry := goes.NewEventStoreEntry(4084, 1, 1, data)
		temp := entry.ToBinary()
		goes.FromBinary(temp)
	}
}

func Benchmark_KindUri_to_AggregateRootUri(b *testing.B) {
	kindUri := goes.NewAggregateKindUri("namespace", "kind")

	for i := 0; i < b.N; i++ {
		kindUri.ToAggregateRootUri(10000)
	}
}

func Benchmark_MemoryEventStore_Sync_RandomId_AppendOnly_10bytePayload(b *testing.B) {
	b.StopTimer()
	eventStore, _ := goes.Connect("mem://")
	kindUri := goes.NewAggregateKindUri("namespace", "kind")
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(10, 1, 1, data)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ids := make([]int64, b.N)
	for i := 0; i < b.N; i++ {
		ids[i] = rnd.Int63()
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uri := kindUri.ToAggregateRootUri(ids[i])
		complete, _ := eventStore.Append(uri, entry1)
		<-complete
	}
}

func Benchmark_MemoryEventStore_Sync_RandomId_AppendOnly_4084bytePayload(b *testing.B) {
	b.StopTimer()
	eventStore, _ := goes.Connect("mem://")
	kindUri := goes.NewAggregateKindUri("namespace", "kind")
	data := make([]byte, 4084)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(4084, 1, 1, data)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ids := make([]int64, b.N)
	for i := 0; i < b.N; i++ {
		ids[i] = rnd.Int63()
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uri := kindUri.ToAggregateRootUri(ids[i])
		complete, _ := eventStore.Append(uri, entry1)
		<-complete
	}
}
