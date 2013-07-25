package eventstore_test

import (
	//"encoding/binary"
	goes "github.com/vizidrix/eventstore"
	"log"
	"testing"
)

func ignore() { log.Println("") }

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

func Test_Should_produce_correct_RelativePath_for_NamespaceUri(t *testing.T) {
	// Arrange
}

/*
func (binary *[]byte) FromBinary() (*MyEvent, error) {
	return nil, nil
}
*/

/*
func Test_Should_create_containing_folder_on_connect(t *testing.T) {
	// Arrange
	//options := map[string]string{
	//	"path": "/eventstore/",
	//}

	// Act
	id := eventstore.NewKey()
	es, _ := eventstore.Connect("/eventstore/")

	// Assert
	domain, _ := es.Domain("namespace")
	kind, err := domain.Kind("person")
	aggregate, err := kind.Aggregate(id)

	log.Printf(
		"%s \n\t %s \n\t\t %s \n\t\t\t %s \n\t *Error: %s",
		es, domain, kind, aggregate, err)

	event := &MyEvent{
		Value: "stuff",
		Index: 10,
	}

	aggregate.Append(event)

	//events := aggregate.LoadAll()

	//log.Printf("Read %d events", len(events))

	t.Fail()
}
*/
/*
func Test_Should_create_aggregate_by_uri(t *testing.T) {
	// Create a new id for the aggregate
	id := eventstore.NewKey()

	// Create the aggregate uri
	uri := fmt.Sprintf("/namespace/person/%d", id)

	es, _ := eventstore.Connect("/eventstore/")

	aggregate, err := es.Aggregate(uri)

	event := &MyEvent{
		Value: "stuff",
		Index: 10,
	}

	aggregate.Append(event)
}
*/

/*
func Benchmark_NewKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		goes.NewKey()
	}
}

func Benchmark_AppendEvent(b *testing.B) {
	//eventStore, err := goes.Connect("/eventstore/")
	//domain, _ := es.Domain("namespace")
	//kind, err := domain.Kind("person")

	for i := 0; i < b.N; i++ {
		//id := goes.NewKey()

		//aggregate, err := kind.Aggregate(id)
	}
}
*/
