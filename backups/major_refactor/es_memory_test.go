package eventstore_test

import (
	//goes "github.com/vizidrix/eventstore"
	"log"
	"testing"
)

func ignore_memoryeventstore_test() { log.Println("") }

// 256 4k events / 1mb
// 4124 ns/ 1 op = 2424 op / ms = 2,424,000 op / s = 9468 mb / s

/*******************/

const (
	MemoryUri = "mem://"
)

func Test_Memory_Should_return_empty_slice_for_new_id(t *testing.T) {
	EventStore_Should_return_empty_slice_for_new_id(t, MemoryUri)
}

func Test_Memory_Should_return_single_matching_event_for_existing_id(t *testing.T) {
	EventStore_Should_return_single_matching_event_for_existing_id(t, MemoryUri)
}

func Test_Memory_Should_return_middle_events_for_version_range(t *testing.T) {
	EventStore_Should_return_middle_events_for_version_range(t, MemoryUri)
}

func Test_Memory_Should_return_two_matching_events_for_existing_ids(t *testing.T) {
	EventStore_Should_return_two_matching_events_for_existing_ids(t, MemoryUri)
}

func Test_Memory_Should_not_panic_when_range_is_too_long(t *testing.T) {
	EventStore_Should_not_panic_when_range_is_too_long(t, MemoryUri)
}

/*******************/

/*
func Benchmark_MemoryES_PutGet_10bytePayload(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 10, 1)
}

func Benchmark_MemoryES_PutGet_10bytePayloads_x10(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 10, 10)
}

func Benchmark_MemoryES_PutGet_10bytePayloads_x100(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 10, 100)
}

func Benchmark_MemoryES_PutGet_10bytePayloads_x1000(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 10, 1000)
}

func Benchmark_MemoryES_PutGet_10bytePayloads_x10000(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 10, 10000)
}

func Benchmark_MemoryES_PutGet_1024bytePayload(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 1024, 1)
}

func Benchmark_MemoryES_PutGet_1024bytePayloads_x10(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 1024, 10)
}

func Benchmark_MemoryES_PutGet_1024bytePayloads_x100(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 1024, 100)
}

func Benchmark_MemoryES_PutGet_1024bytePayloads_x1000(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 1024, 1000)
}

func Benchmark_MemoryES_PutGet_1024bytePayloads_x10000(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 1024, 10000)
}

func Benchmark_MemoryES_PutGet_4087bytePayload(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 4087, 1)
}

func Benchmark_MemoryES_PutGet_4087bytePayloads_x10(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 4087, 10)
}

func Benchmark_MemoryES_PutGet_4087bytePayloads_x100(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 4087, 100)
}

func Benchmark_MemoryES_PutGet_4087bytePayloads_x1000(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 4087, 1000)
}

func Benchmark_MemoryES_PutGet_4087bytePayloads_x10000(b *testing.B) {
	Run_PutGet(b, MemoryUri, "namespace", "kind", 4087, 10000)
}
*/
