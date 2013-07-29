package eventstore_test

import (
	//"encoding/binary"
	//"bytes"
	goes "github.com/vizidrix/eventstore"
	//. "github.com/vizidrix/eventstore/test_utils"
	//"hash/crc32"
	"log"
	//"math/rand"
	"testing"
	//"time"
)

func ignore_memoryeventstore_test() { log.Println("") }

func Test_Should_return_empty_slice_for_new_id(t *testing.T) {
	EventStore_Should_return_empty_slice_for_new_id(t, "mem://")
}

func Test_Should_return_single_matching_event_for_existing_id(t *testing.T) {
	EventStore_Should_return_single_matching_event_for_existing_id(t, "mem://")
}

func Test_Should_return_middle_events_for_version_range(t *testing.T) {
	EventStore_Should_return_middle_events_for_version_range(t, "mem://")
}

func Test_Should_return_two_matching_events_for_existing_ids(t *testing.T) {
	EventStore_Should_return_two_matching_events_for_existing_ids(t, "mem://")
}

func Test_Should_not_panic_when_range_is_too_long(t *testing.T) {
	EventStore_Should_not_panic_when_range_is_too_long(t, "mem://")
}

func Test_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t *testing.T) {
	EventStore_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t, "mem://")
}

func Test_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t *testing.T) {
	EventStore_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t, "mem://")
}

func Test_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t *testing.T) {
	EventStore_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t, "mem://")
}

func Benchmark_Create_Serialize_DeSerialize_EventStoreEntry_10bytePayload(b *testing.B) {
	b.StopTimer()
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		entry := goes.NewEventStoreEntryFrom(1, data)
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
		entry := goes.NewEventStoreEntryFrom(1, data)
		temp := entry.ToBinary()
		goes.FromBinary(temp)
	}
}

func Benchmark_KindUri_to_AggregateRootUri(b *testing.B) {
	kind := goes.NewAggregateKind("namespace", "kind")

	for i := 0; i < b.N; i++ {
		kind.ToAggregateUri(10000)
	}
}

// 256 4k events / 1mb
// 4124 ns/ 1 op = 2424 op / ms = 2,424,000 op / s = 9468 mb / s

func Benchmark_MemoryEventStore_AppendOnly_10bytePayload(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 10, 1)
}

func Benchmark_MemoryEventStore_AppendOnly_4084bytePayload(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 4084, 1)
}

func Benchmark_MemoryEventStore_ReadOnly_10bytePayload(b *testing.B) {
	Run_ReadOnly(b, "mem://", "namespace", "kind", 10, 1)
}

func Benchmark_MemoryEventStore_ReadOnly_4084bytePayload(b *testing.B) {
	Run_ReadOnly(b, "mem://", "namespace", "kind", 4084, 1)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_10bytePayload(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 10, 1)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_4084bytePayload(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 4084, 1)
}

func Benchmark_MemoryEventStore_AppendOnly_20_10bytePayloads(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 10, 20)
}

func Benchmark_MemoryEventStore_AppendOnly_20_4084bytePayloads(b *testing.B) {
	Run_AppendOnly(b, "mem://", "namespace", "kind", 4084, 20)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_20_10bytePayloads(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 10, 20)
}

func Benchmark_MemoryEventStore_AppendAndReadAll_20_4084bytePayloads(b *testing.B) {
	Run_AppendAndReadAll(b, "mem://", "namespace", "kind", 4084, 20)
}
