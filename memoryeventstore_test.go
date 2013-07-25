package eventstore_test

import (
	//"encoding/binary"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	"log"
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
	entry, _ := goes.NewEventStoreEntry(10, 1, 1, 1, data)
	eventStore.Append(uri, entry)

	// Act
	errors := eventStore.LoadAll(uri, events)

	// Assert
	select {
	case event := <-events:
		{
			//AreEqual(t, 1, len(events), "Event list should have had appended value")

			AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
			AreEqual(t, int32(1), event.CRC(), "CRC should have been calculated")
			AreEqual(t, int64(1), event.UnixTimeStamp(), "TimeStamp should have been set")
			AreEqual(t, byte(1), event.EventType(), "EvenType should have been set")
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
	entry1, _ := goes.NewEventStoreEntry(10, 1, 1, 1, data)
	entry2, _ := goes.NewEventStoreEntry(10, 2, 2, 2, data)
	eventStore.Append(uri, entry1)
	eventStore.Append(uri, entry2)

	// Act
	errors := eventStore.LoadAll(uri, events)

	// Assert
	select {
	case event := <-events:
		{
			//AreEqual(t, 1, len(events), "Event list should have had appended value")

			AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
			AreEqual(t, int32(1), event.CRC(), "CRC should have been calculated")
			AreEqual(t, int64(1), event.UnixTimeStamp(), "TimeStamp should have been set")
			AreEqual(t, byte(1), event.EventType(), "EvenType should have been set")
			AreAllEqual(t, data, event.Data(), "Data should have been set")

			select {
			case event := <-events:
				{
					//AreEqual(t, 1, len(events), "Event list should have had appended value")

					AreEqual(t, int32(10), event.Length(), "Length should have been int32 10")
					AreEqual(t, int32(2), event.CRC(), "CRC should have been calculated")
					AreEqual(t, int64(2), event.UnixTimeStamp(), "TimeStamp should have been set")
					AreEqual(t, byte(2), event.EventType(), "EvenType should have been set")
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

	/*
		AreEqual(t, 2, len(events), "Event list should have had appended value")

		AreEqual(t, int32(10), events[0].Length(), "Length should have been int32 10")
		AreEqual(t, int32(1), events[0].CRC(), "CRC should have been calculated")
		AreEqual(t, int64(1), events[0].UnixTimeStamp(), "TimeStamp should have been set")
		AreEqual(t, byte(1), events[0].EventType(), "EvenType should have been set")
		AreAllEqual(t, data, events[0].Data(), "Data should have been set")

		if len(events) < 2 {
			return
		}
		AreEqual(t, int32(10), events[1].Length(), "Length should have been int32 10")
		AreEqual(t, int32(2), events[1].CRC(), "CRC should have been calculated")
		AreEqual(t, int64(2), events[1].UnixTimeStamp(), "TimeStamp should have been set")
		AreEqual(t, byte(2), events[1].EventType(), "EvenType should have been set")
		AreAllEqual(t, data, events[1].Data(), "Data should have been set")
	*/
}
