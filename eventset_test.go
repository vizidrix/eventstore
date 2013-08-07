package eventstore_test

import (
	"fmt"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	"log"
	"testing"
)

func ignore_eventset_test() {
	log.Printf(fmt.Sprintf(""))
}

func MakeByteSlice(size int) []byte {
	result := make([]byte, size)
	for i := 0; i < size; i++ {
		result[i] = byte(i % 255)
	}
	return result
}

func Test_Should_put_slice_of_events_in_batch(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()

	// Act
	eventSet, err := eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
		{2, MakeByteSlice(12)},
		{3, MakeByteSlice(14)},
	}...)

	// Assert
	IsNil(t, err, "Should have put events")
	AreEqual(t, 3, eventSet.Count(), "Should have put all 3 values")
}

func Test_Should_return_err_for_event_over_max_size(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()

	// Act
	eventSet, err := eventSet.Put(goes.Event{1, MakeByteSlice(int(goes.MaxUint16) + 1)})

	// Assert
	IsNotNil(t, err, "Should have failed to put event")
}

func Test_Should_return_err_for_invalid_range(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	ranges := []struct {
		Start int
		Stop  int
	}{
		{-1, 0},
		{4, 1},
		{2, 1},
		{1, 1},
	}

	// Act
	eventSet2, _ := eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
		{2, MakeByteSlice(12)},
		{3, MakeByteSlice(14)},
	}...)

	// Assert
	for _, r := range ranges {
		_, err := eventSet2.GetSlice(r.Start, r.Stop)
		IsNotNil(t, err, fmt.Sprintf("Should have failed for [%d, %d]", r.Start, r.Stop))
	}

}

func Test_Should_return_empty_event_slice_for_empty_set(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()

	// Act
	events, err := eventSet.Get()

	// Assert
	IsNil(t, err, "Should have allowed read from empty set")
	AreEqual(t, 0, events.Count(), "Should have returned an empty slice")
}

func Test_Should_return_correct_sub_range_from_set(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	eventSet, err := eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
		{2, MakeByteSlice(12)},
		{3, MakeByteSlice(14)},
		{4, MakeByteSlice(16)},
		{5, MakeByteSlice(18)},
	}...)

	// Act
	events, err := eventSet.GetSlice(2, 4)

	// Assert
	IsNil(t, err, "Should have allowed read slice from populated set")
	AreEqual(t, 2, events.Count(), "Should have returned two results")
	AreEqual(t, uint16(3), events.EventTypeAt(0), "Should have returned the second event")
	AreEqual(t, 14, len(events.DataAt(0)), "Should have returned 14 bytes of data")
	AreEqual(t, uint16(4), events.EventTypeAt(1), "Should have returned the third event")
	AreEqual(t, 16, len(events.DataAt(1)), "Should have returned 16 bytes of data")

}

func Test_Should_limit_end_index_to_max_available_from_set(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	eventSet, err := eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
		{2, MakeByteSlice(12)},
		{3, MakeByteSlice(14)},
		{4, MakeByteSlice(16)},
		{5, MakeByteSlice(18)},
	}...)

	// Act
	events, err := eventSet.GetSlice(2, 40)

	// Assert
	IsNil(t, err, "Should have allowed read slice from populated set")
	AreEqual(t, 3, events.Count(), "Should have returned three results")
}

func Test_Should_checksum_successfully_for_empty_dataset(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()

	// Act
	result := eventSet.CheckSum()

	// Assert
	IsNil(t, result, "Should have passed check sum")
}

func Test_Should_checksum_successfully_for_valid_dataset(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	eventSet, _ = eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
		{2, MakeByteSlice(12)},
		{3, MakeByteSlice(14)},
	}...)

	// Act
	err := eventSet.CheckSum()

	// Assert
	IsNil(t, err, "Should not have failed check sum")
}

func Test_Should_return_error_if_end_index_less_than_zero(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()

	// Act
	_, err := eventSet.GetSlice(0, -1)

	// Assert
	IsNotNil(t, err, "Should have failed due to end index bounds")
}

func Test_Should_retrieve_from_multiple_puts(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	eventSet, _ = eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
		{2, MakeByteSlice(12)},
		{3, MakeByteSlice(14)},
	}...)

	// Act
	eventSet, _ = eventSet.Put([]goes.Event{
		{4, MakeByteSlice(16)},
		{5, MakeByteSlice(18)},
		{6, MakeByteSlice(20)},
	}...)

	// Act
	events, err := eventSet.Get()

	// Assert
	IsNil(t, err, "Should not have failed Get")
	AreEqual(t, 6, events.Count(), "Should have brought back all records")
	AreEqual(t, uint16(2), events.EventTypeAt(1), "Should have retained initial headers")
	AreEqual(t, uint16(5), events.EventTypeAt(4), "Should have appended additional headers")
}

func Test_Should_expand_headers(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	eventSet, _ = eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10)},
	}...)

	// Act
	batch := new([goes.HEADER_SLICE_SIZE]goes.Event)
	for i := 0; i < goes.HEADER_SLICE_SIZE; i++ {
		batch[i] = goes.Event{uint16(i), MakeByteSlice(i)}
	}
	eventSet, _ = eventSet.Put(batch[:]...)

	// Act
	events, err := eventSet.Get()

	// Assert
	IsNil(t, err, "Should not have failed Get")
	AreEqual(t, len(batch)+1, events.Count(), "Should have brought back all records")
	AreEqual(t, uint16(0), events.EventTypeAt(1), "Should have retained initial headers")
	AreEqual(t, uint16(3), events.EventTypeAt(4), "Should have appended additional headers")
}

func Test_Should_return_err_if_event_is_too_large(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()

	// Act
	eventSet, err := eventSet.Put([]goes.Event{
		{1, MakeByteSlice(10000)},
		{2, MakeByteSlice(12000)},
		{3, MakeByteSlice(140000)}, // make it 1400000
	}...)
	//t.Fail()
	// Assert
	IsNotNil(t, err, "Should have throw event too large error")
}

func Test_Should_be_able_to_resize_header_slice(t *testing.T) {
	// Arrange
	eventSet := goes.NewEmptyEventSet()
	events := make([]goes.Event, (goes.HEADER_SLICE_SIZE/8)+1)
	for i := range events {
		events[i] = goes.Event{uint16(i), MakeByteSlice(i)}
	}

	// Act
	eventSet, _ = eventSet.Put(events...)

	// Assert
	IsNotNil(t, eventSet, "Should have returned valid event set")
}
