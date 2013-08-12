package eventstore_test

import (
	goes "github.com/vizidrix/eventstore"
	"log"
	"testing"
)

func Test_Stuff(t *testing.T) {
	result, err := goes.Stuff(10)
	log.Printf("Result: %d - err: %s", result, err)

	t.Fail()

}
