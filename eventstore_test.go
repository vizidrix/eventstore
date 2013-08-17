package eventstore_test

import (
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	//"hash/crc32"
	"log"
	"testing"
	//"time"
)

func ignore() { log.Println("") }

// go test -v github.com/vizidrix/eventstore -bench . -cpuprofile /go/prof.out -cpu 1,2,4,8

func Test_Should_put_single_event_into_write_store(t *testing.T) {
	writer, err := goes.NewESWriter("/go/esdata/")
	//log.Printf("Writer: % x", writer)
	defer writer.Close()
	if err != nil {
		log.Printf("Error opening ES Writer")
		t.Fail()
	}
	//max_size := 512 - 32

	data := make([][]byte, 4)
	data[0] = MakeByteSlice(50)
	data[1] = MakeByteSlice(100)
	data[2] = MakeByteSlice(200)
	data[3] = MakeByteSlice(400)

	var domain uint32 = 100
	var kind uint32 = 110
	var aggregate uint64 = 111

	//writer.AllocBatch(domain, kind, aggregate, 4)
	//writer.AllocBatch(domain, kind, aggregate, 4)

	batch, err := writer.AllocBatch(domain, kind, aggregate, 9, 4)
	if err != nil {
		log.Printf("Error allocating branch: %s", err)
	}

	for i := 0; i < 4; i++ {
		batch.Entries[i].EventType = uint16(i)
		batch.Entries[i].EventSize = uint16(len(data[i]))
		batch.Entries[i].CopyFrom(data[i])
		//copy(batch.Entries[i].GetEventData()[:], data[i])
	}
	batch.Publish()

	// get a set of put commands from db
	// wrap them in managed slices
	// fill them over here from wherever
	// publish commands
}

func Get_Event(eventType uint16, size int) goes.Event {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(i | 0xFF)
	}
	return goes.Event{
		EventType: eventType,
		Data:      data,
	}
}

type VisitorLogged struct {
	UnixNSTimeStamp int64
	IPv4Address     int64
	IPv6Header      int64
	IPv6Address     int64
	Referrer        string
}

func Test_Should(t *testing.T) {
	/*es, err := goes.Connect("fs://eventstore/")
	defer func() {
		if es != nil {
			es.Close()
		}
	}()

	if err != nil {
		t.Fail()
	}*/
	t.Fail()
}

func Test_Should_try_to_connect_to_MemoryEventstore_with_correct_path(t *testing.T) {
	// Arrange
	path := "mem://"

	// Act
	eventStore, err := goes.Connect(path)
	defer eventStore.Close()

	// Assert
	if err != nil {
		log.Printf("Connect failed: %s", err)
		t.Fail()
		return
	}
	switch es := (eventStore).(type) {
	case *goes.MemoryES:
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
	path := "fs://eventstore/"

	// Act
	eventStore, err := goes.Connect(path)
	defer eventStore.Close()

	// Assert
	if err != nil {
		log.Printf("Connect failed: %s", err)
		t.Fail()
		return
	}
	switch es := (eventStore).(type) {
	case *goes.FileSystemES:
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
		es, err := goes.Connect(path)
		defer func() { // Make sure to clean up if it passed
			if es != nil {
				es.Close()
			}
		}()
		// Assert
		if err == nil {
			log.Printf("Invalid path should have raised an err: %s", path)
			t.Fail()
		}
	}
}

func Test_Should_put_a_bunch_of_entries(t *testing.T) {
	connString := "mem://"
	eventCount := 100
	kind := goes.NewAggregateKind("namespace", "type")
	event := Get_Event(10, 10)
	i := 0
	index := 0
	for i = 0; i < 1; i++ {
		eventStore, _ := goes.Connect(connString)
		defer eventStore.Close()
		kindPartition := eventStore.Kind(kind)
		aggregatePartition := kindPartition.Id(uint64(i))

		for index = 0; index < eventCount; index++ {
			aggregatePartition.Put(event)

		}
	}
}

func EventStore_Should_return_empty_slice_for_new_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	defer eventStore.Close()
	kind := goes.NewAggregateKind("namespace", "type")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Id(1)

	// Act
	events, _ := aggregatePartition.Get()

	// Assert
	AreEqual(t, 0, events.Count(), "Shouldn't have received any events")
}

func EventStore_Should_return_single_matching_event_for_existing_id(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	defer eventStore.Close()
	kind := goes.NewAggregateKind("namespace", "type")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Id(1)
	event := Get_Event(1, 10)
	aggregatePartition.Put(event)

	// Act
	events, _ := aggregatePartition.Get()

	// Assert
	AreEqual(t, uint16(10), events.LengthAt(0), "Length should have been int32 10")
	AreEqual(t, uint16(1), events.EventTypeAt(0), "EvenType should have been set")
	AreAllEqual(t, event.Data, events.DataAt(0), "Data should have been set")
}

func EventStore_Should_return_middle_events_for_version_range(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	defer eventStore.Close()
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Id(1)
	for index := 0; index < 5; index++ {
		event := Get_Event(uint16(index), 10)
		aggregatePartition.Put(event)
	}

	// Act
	events, _ := aggregatePartition.GetSlice(2, 4)

	// Assert
	for index := 0; index < 2; index++ {
		AreEqual(t, uint16(10), events.LengthAt(index), "Length should have been int32 10")
		AreEqual(t, uint16(index+2), events.EventTypeAt(index), "EvenType should have been set")
	}
}

func EventStore_Should_return_two_matching_events_for_existing_ids(t *testing.T, connString string) {
	// Arrange
	eventStore, _ := goes.Connect(connString)
	defer eventStore.Close()
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Id(1)
	event1 := Get_Event(1, 10)
	event2 := Get_Event(2, 10)

	aggregatePartition.Put(event1, event2)

	// Act
	events, _ := aggregatePartition.Get()

	// Assert
	AreEqual(t, uint16(10), events.LengthAt(0), "Length should have been int32 10")
	AreEqual(t, uint16(1), events.EventTypeAt(0), "EvenType should have been set")
	AreAllEqual(t, event1.Data, events.DataAt(0), "Data should have been set")
	AreEqual(t, uint16(10), events.LengthAt(1), "Length should have been int32 10")
	AreEqual(t, uint16(2), events.EventTypeAt(1), "EvenType should have been set")
	AreAllEqual(t, event1.Data, events.DataAt(1), "Data should have been set")
}

func EventStore_Should_not_panic_when_range_is_too_long(t *testing.T, connString string) {
	eventStore, _ := goes.Connect(connString)
	defer eventStore.Close()
	kind := goes.NewAggregateKind("namespace", "kind")
	kindPartition := eventStore.Kind(kind)
	aggregatePartition := kindPartition.Id(1)
	event := Get_Event(1, goes.MAX_EVENT_SIZE)
	aggregatePartition.Put(event)
	aggregatePartition.GetSlice(0, 4)
}
