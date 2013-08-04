package eventstore_test

import (
	//goes "github.com/vizidrix/eventstore"
	//. "github.com/vizidrix/eventstore/test_utils"
	//"hash/crc32"
	"log"
	//"testing"
	//"time"
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

func EventStore_Should_return_empty_slice_for_new_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "type")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)

	// Act
	events, _ := aggregatePartition.Get()

	// Assert
	AreEqual(t, 0, len(events), "Shouldn't have received any events")
}

func Test_Should_put_a_bunch_of_entries(t *testing.T) {
	connString := "mem://"
	eventCount := 100
	kind := goes.NewAggregateKind("namespace", "type")
	eventStoreEntry := Get_EventStoreEntry(10)
	i := 0
	index := 0
	for i = 0; i < 1; i++ {
		eventStore, _ := goes.Connect(connString)
		kindPartition := eventStore.Kind(kind)
		aggregatePartition := kindPartition.Aggregate(int64(i))

		for index = 0; index < eventCount; index++ {
			aggregatePartition.Put(eventStoreEntry)

		}
	}
}

func EventStore_Should_return_single_matching_event_for_existing_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "type")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	data := Get_EventStoreEntry(10).Data()
	entry := goes.NewEventStoreEntry(10, 1, 1, data)
	aggregatePartition.Put(entry)

	// Act
	events, _ := aggregatePartition.Get()

	// Assert
	AreEqual(t, uint16(10), events[0].Length(), "Length should have been int32 10")
	AreEqual(t, uint16(1), events[0].EventType(), "EvenType should have been set")
	AreEqual(t, uint32(1), events[0].CRC(), "CRC should have been calculated")
	AreAllEqual(t, data, events[0].Data(), "Data should have been set")
}

func EventStore_Should_return_middle_events_for_version_range(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	data := Get_EventStoreEntry(10).Data()
	for index := 0; index < 5; index++ {
		entry := goes.NewEventStoreEntry(10, uint16(index), uint32(index), data)
		aggregatePartition.Put(entry)
	}

	// Act
	events, _ := aggregatePartition.GetSlice(2, 3)

	// Assert
	for index := 0; index < 2; index++ {
		AreEqual(t, uint16(10), events[index].Length(), "Length should have been int32 10")
		AreEqual(t, uint16(index+2), events[index].EventType(), "EvenType should have been set")
		AreEqual(t, uint32(index+2), events[index].CRC(), "CRC should have been calculated")
		AreAllEqual(t, data, events[index].Data(), "Data should have been set")
	}
}

func EventStore_Should_return_two_matching_events_for_existing_ids(t *testing.T, connString string) {
	// Arrange
	log.Printf("\n\nReturn two\n\n")
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	data := Get_EventStoreEntry(10).Data()
	entry1 := goes.NewEventStoreEntry(10, 0, 0, data)
	entry2 := goes.NewEventStoreEntry(10, 1, 1, data)

	aggregatePartition.Put(entry1)
	aggregatePartition.Put(entry2)

	// Act
	events, _ := aggregatePartition.Get()

	// Assert
	for i := 0; i < 2; i++ {
		AreEqual(t, uint16(10), events[i].Length(), "Length should have been int32 10")
		AreEqual(t, uint16(i), events[i].EventType(), "EvenType should have been set")
		AreEqual(t, uint32(i), events[i].CRC(), "CRC should have been calculated")
		AreAllEqual(t, data, events[i].Data(), "Data should have been set")
	}
}

func EventStore_Should_not_panic_when_range_is_too_long(t *testing.T, connString string) {
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	data := Get_EventStoreEntry(goes.MAX_EVENT_SIZE).Data()
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE, 1, 1, data)
	aggregatePartition.Put(entry1)
	aggregatePartition.GetSlice(0, 4)
}

func EventStore_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t *testing.T, connString string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	data := Get_EventStoreEntry(goes.MAX_EVENT_SIZE + 1).Data()
	entry1 := goes.NewEventStoreEntry(goes.MAX_EVENT_SIZE+1, 1, 1, data)
	aggregatePartition.Put(entry1)
}

func EventStore_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t *testing.T, connString string) {
	defer func() {
		if r := recover(); r == nil {
			log.Printf("Should have raised a panic")
			t.Fail()
		}
	}()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	data := Get_EventStoreEntry(3082).Data()
	entry1 := goes.NewEventStoreEntry(3083, 1, 1, data) // <- invalid length!
	aggregatePartition.Put(entry1)
}

func EventStore_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t *testing.T, connString string) {
	count := 2

	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Aggregate(1)
	entry1 := Get_EventStoreEntry(10)

	roundComplete := make(chan struct{})

	for i := 0; i < count; i++ {
		go func() {
			aggregatePartition.Put(entry1)
			aggregatePartition.Get()
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

*/
