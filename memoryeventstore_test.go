package eventstore_test

import (
	//"encoding/binary"
	goes "github.com/vizidrix/eventstore"
	"log"
	"testing"
)

func ignore_memoryeventstore_test() { log.Println("") }

func Test_Should_return_nil_for_new_id(t *testing.T) {
	// Arrange
	eventStore, _ := goes.Connect("mem://")

	uri := goes.NewAggregateRootUri("namespace", "type", 1)

	log.Printf("Got es: %s and uri %s", eventStore, uri)

	t.Fail()
}
