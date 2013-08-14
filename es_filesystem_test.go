package eventstore_test

import (
	"log"
	"testing"
)

func ignore_filesystemeventstore_test() { log.Println("") }

const (
	FileSystemUri = "fs://"
)

func Test_FileSystem_Should_return_empty_slice_for_new_id(t *testing.T) {
	EventStore_Should_return_empty_slice_for_new_id(t, FileSystemUri)
}

func Test_FileSystem_Should_return_single_matching_event_for_existing_id(t *testing.T) {
	EventStore_Should_return_single_matching_event_for_existing_id(t, FileSystemUri)
}

func Test_FileSystem_Should_return_middle_events_for_version_range(t *testing.T) {
	EventStore_Should_return_middle_events_for_version_range(t, FileSystemUri)
}

func Test_FileSystem_Should_return_two_matching_events_for_existing_ids(t *testing.T) {
	EventStore_Should_return_two_matching_events_for_existing_ids(t, FileSystemUri)
}

func Test_FileSystem_Should_not_panic_when_range_is_too_long(t *testing.T) {
	EventStore_Should_not_panic_when_range_is_too_long(t, FileSystemUri)
}
