package eventstore_test

import (
	"log"
	"testing"
)

func ignore_fragmentfilesystemeventstore_test() {
	var temp testing.T
	log.Println("%s", temp)

}

// 256 4k events / 1mb
// 4124 ns/ 1 op = 2424 op / ms = 2,424,000 op / s = 9468 mb / s

/*******************/

const (
	FragmentFileSystemUri = "ffs://frageventstore/"
)

/*
func Test_FragmentFileSystem_Sync_Should_return_empty_slice_for_new_id(t *testing.T) {
	EventStoreSync_Should_return_empty_slice_for_new_id(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_return_single_matching_event_for_existing_id(t *testing.T) {
	EventStoreSync_Should_return_single_matching_event_for_existing_id(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_return_middle_events_for_version_range(t *testing.T) {
	EventStoreSync_Should_return_middle_events_for_version_range(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_return_two_matching_events_for_existing_ids(t *testing.T) {
	EventStoreSync_Should_return_two_matching_events_for_existing_ids(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_not_panic_when_range_is_too_long(t *testing.T) {
	EventStoreSync_Should_not_panic_when_range_is_too_long(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t *testing.T) {
	EventStoreSync_Should_panic_when_event_length_greater_than_max_in_unchecked_ctor(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t *testing.T) {
	EventStoreSync_Should_panic_when_reported_event_length_greater_than_actual_in_unchecked_ctor(t, FragmentFileSystemUri)
}

func Test_FragmentFileSystem_Sync_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t *testing.T) {
	EventStoreSync_Should_fail_if_write_index_is_not_unique_when_expected_to_be(t, FragmentFileSystemUri)
}

/*******************

func Benchmark_Sync_FragmentFileSystemEventStore_AppendOnly_10bytePayload(b *testing.B) {
	Run_AppendOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 10, 1)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendOnly_4087bytePayload(b *testing.B) {
	Run_AppendOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 1)
}

func Benchmark_Sync_FragmentFileSystemEventStore_ReadOnly_10bytePayload(b *testing.B) {
	Run_ReadOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 10, 1)
}

func Benchmark_Sync_FragmentFileSystemventStore_ReadOnly_4087bytePayload(b *testing.B) {
	Run_ReadOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 1)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendAndReadAll_10bytePayload(b *testing.B) {
	Run_AppendAndReadAllSync(b, FragmentFileSystemUri, "namespace", "kind", 10, 1)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendAndReadAll_4087bytePayload(b *testing.B) {
	Run_AppendAndReadAllSync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 1)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendOnly_20_10bytePayloads(b *testing.B) {
	Run_AppendOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 10, 20)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendOnly_20_4087bytePayloads(b *testing.B) {
	Run_AppendOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 20)
}

func Benchmark_Sync_FragmentFileSystemEventStore_ReadOnly_20_10bytePayloads(b *testing.B) {
	Run_ReadOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 10, 20)
}

func Benchmark_Sync_FragmentFileSystemEventStore_ReadOnly_20_4087bytePayloads(b *testing.B) {
	Run_ReadOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 20)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendAndReadAll_20_10bytePayloads(b *testing.B) {
	Run_AppendAndReadAllSync(b, FragmentFileSystemUri, "namespace", "kind", 10, 20)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendAndReadAll_20_4087bytePayloads(b *testing.B) {
	Run_AppendAndReadAllSync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 20)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendOnly_100_10bytePayloads(b *testing.B) {
	Run_AppendOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 10, 100)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendOnly_100_4087bytePayloads(b *testing.B) {
	Run_AppendOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 100)
}

func Benchmark_Sync_FragmentFileSystemEventStore_ReadOnly_100_10bytePayloads(b *testing.B) {
	Run_ReadOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 10, 100)
}

func Benchmark_Sync_FragmentFileSystemEventStore_ReadOnly_100_4087bytePayloads(b *testing.B) {
	Run_ReadOnlySync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 100)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendAndReadAll_100_10bytePayloads(b *testing.B) {
	Run_AppendAndReadAllSync(b, FragmentFileSystemUri, "namespace", "kind", 10, 100)
}

func Benchmark_Sync_FragmentFileSystemEventStore_AppendAndReadAll_100_4087bytePayloads(b *testing.B) {
	Run_AppendAndReadAllSync(b, FragmentFileSystemUri, "namespace", "kind", 4087, 100)
}

/*******************/
