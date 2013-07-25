package eventstore_test

import (
	//"encoding/binary"
	//"bytes"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
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
	errors := eventStore.LoadAll(uri, events)

	// Assert
	select {
	case event := <-events:
		{
			log.Printf("Shouldn't have received any events: %s", event)
			t.Fail()
		}
	case err := <-errors:
		{
			if err == nil { // Returns Item not found error
				//log.Printf("Shouldn't have received any errors: %s", err)
				t.Fail()
			}
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
		}
	}
	/*if len(events) != 0 {
		log.Printf("Event list should have been empty")
		t.Fail()
	}*/
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
	eventStore.Append(uri, entry)

	// Act
	errors := eventStore.LoadAll(uri, events)

	// Assert
	select {
	case event := <-events:
		{
			AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
			AreEqual(t, byte(1), event.EventType(), "EvenType should have been set")
			AreEqual(t, int32(1), event.CRC(), "CRC should have been calculated")
			AreAllEqual(t, data, event.Data(), "Data should have been set")
		}
	case err := <-errors:
		{
			log.Printf("Shouldn't have received any errors: %s", err)
			t.Fail()
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
		}
	}

}

func Test_Should_return_two_matching_events_for_existing_ids(t *testing.T) {
	// Arrange
	eventStore, _ := goes.Connect("mem://")
	uri := goes.NewAggregateRootUri("namespace", "type", 1)
	events := make(chan *goes.EventStoreEntry, 1)
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	entry1 := goes.NewEventStoreEntry(10, 1, 1, data)
	entry2 := goes.NewEventStoreEntry(10, 2, 2, data)
	eventStore.Append(uri, entry1)
	eventStore.Append(uri, entry2)

	// Act
	errors := eventStore.LoadAll(uri, events)

	// Assert
	select {
	case event := <-events:
		{
			AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
			AreEqual(t, byte(1), event.EventType(), "EvenType should have been set")
			AreEqual(t, int32(1), event.CRC(), "CRC should have been calculated")
			AreAllEqual(t, data, event.Data(), "Data should have been set")

			select {
			case event := <-events:
				{
					AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
					AreEqual(t, byte(2), event.EventType(), "EvenType should have been set")
					AreEqual(t, int32(2), event.CRC(), "CRC should have been calculated")
					AreAllEqual(t, data, event.Data(), "Data should have been set")
				}
			case err := <-errors:
				{
					log.Printf("Shouldn't have received any errors: %s", err)
					t.Fail()
				}
			case <-time.After(1 * time.Millisecond):
				{
					log.Printf("Shouldn't have timed out")
					t.Fail()
				}
			}
		}
	case err := <-errors:
		{
			log.Printf("Shouldn't have received any errors: %s", err)
			t.Fail()
		}
	case <-time.After(1 * time.Millisecond):
		{
			log.Printf("Shouldn't have timed out")
			t.Fail()
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
	kindUri := goes.NewAggregateKindUri("namespace", "type")
	//rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		kindUri.ToAggregateRootUri(10000)
	}
}

func Benchmark_MemoryEventStore_Sync_RandomId_AppendOnly_10bytePayload(b *testing.B) {
	b.StopTimer()
	eventStore, _ := goes.Connect("mem://")
	kindUri := goes.NewAggregateKindUri("namespace", "type")
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
		complete, errors := eventStore.Append(uri, entry1)
		select {
		case <-complete:
			{
			}
		case err := <-errors:
			{
				log.Printf("Error: %s", err)
			}
		}
	}
}

func Benchmark_MemoryEventStore_Sync_RandomId_AppendOnly_4084bytePayload(b *testing.B) {
	b.StopTimer()
	eventStore, _ := goes.Connect("mem://")
	kindUri := goes.NewAggregateKindUri("namespace", "type")
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
		complete, errors := eventStore.Append(uri, entry1)
		select {
		case <-complete:
			{
			}
		case err := <-errors:
			{
				log.Printf("Error: %s", err)
			}
		}
	}
}
