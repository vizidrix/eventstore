package eventstore_test

import (
	eventstore "github.com/vizidrix/eventstore"
	"log"
	"testing"
)

func ignore() { log.Println("") }

type MyEvent struct {
	Value string
	Index byte
}

func (event *MyEvent) ToBinary() ([]byte, error) {
	buffer := make([]byte, len(event.Value)+4)
	log.Printf("ToBinary: Made buffer %d", len(buffer))
	index := 0
	buffer[index] = event.Index
	index++
	for char := range event.Value {
		buffer[index] = byte(char)
		index++
	}
	log.Printf("ToBinary: Pop buffer %d", buffer)
	return buffer, nil
}

func Test_Should_create_containing_folder_on_connect(t *testing.T) {
	// Arrange
	/*options := map[string]string{
		"path": "/eventstore/",
	}*/

	// Act
	id := eventstore.NewKey()
	es, _ := eventstore.Connect("/eventstore/")

	// Assert
	domain, _ := es.Domain("WearShare")
	kind, err := domain.Kind("Person")
	aggregate, err := kind.Aggregate(id)

	log.Printf(
		"%s \n\t %s \n\t\t %s \n\t\t\t %s \n\t *Error: %s",
		es, domain, kind, aggregate, err)

	event := &MyEvent{
		Value: "stuff",
		Index: 10,
	}

	aggregate.Append(event)

	t.Fail()
}

func BenchmarkRandomDataMaker2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		eventstore.NewKey()
	}
}
