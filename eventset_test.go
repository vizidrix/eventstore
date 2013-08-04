package eventstore_test

import (
	"fmt"
	goes "github.com/vizidrix/eventstore"
	. "github.com/vizidrix/eventstore/test_utils"
	"testing"
)

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
	AreEqual(t, 0, len(events), "Should have returned an empty slice")
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
	events, err := eventSet.GetSlice(2, 3)

	// Assert
	IsNil(t, err, "Should have allowed read slice from populated set")
	AreEqual(t, 2, len(events), "Should have returned two results")
	AreEqual(t, uint16(3), events[0].EventType, "Should have returned the second event")
	AreEqual(t, 14, len(events[0].Data), "Should have returned 14 bytes of data")
	AreEqual(t, uint16(4), events[1].EventType, "Should have returned the third event")
	AreEqual(t, 16, len(events[1].Data), "Should have returned 16 bytes of data")

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
	AreEqual(t, 3, len(events), "Should have returned three results")
}

func Test_Should_checksum_successfully_for_valid_data(t *testing.T) {
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
	AreEqual(t, 6, len(events), "Should have brought back all records")
}
