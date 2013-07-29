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

func ignore_eventstore_benchmarks() { log.Printf("") }

func Get_EventStoreEntry(size uint16) *goes.EventStoreEntry {
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	return goes.NewEventStoreEntry(10, 1, 1, data)
}

func Run_Sync_AppendOnly(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry) {
	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
		<-appendComplete
	}
}

func Run_Sync_AppendOnly_Multiples(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry, count int) {
	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		for index := 0; index < count; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
	}
}

func Run_Sync_AppendAndReadAll(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry) {
	events := make(chan *goes.EventStoreEntry, b.N) //*10)

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
		<-appendComplete
		readComplete, _ := eventStore.LoadAll(uri, events)
		<-readComplete
	}
}

func Run_Sync_AppandAndReadAll_Multiples(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry, count int) {
	events := make(chan *goes.EventStoreEntry, b.N) //*10)

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		for index := 0; index < count; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
		readComplete, _ := eventStore.LoadAll(uri, events)
		for index := 0; index < count; index++ {
			<-events
		}
		<-readComplete
	}
}

func Run_Async_AppendAndReadAll(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry) {
	runs := b.N
	roundComplete := make(chan struct{}, runs)

	for i := 0; i < runs; i++ {
		id := int64(i)
		go func() {
			uri := kindUri.ToAggregateRootUri(id)
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
			events := make(chan *goes.EventStoreEntry, 1)
			readComplete, _ := eventStore.LoadAll(uri, events)
			<-readComplete
			roundComplete <- struct{}{}
		}()
	}
	for i := 0; i < runs; i++ {
		<-roundComplete
	}
}

func Run_Async_AppendAndReadAll_Multiples(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry, count int) {
	roundComplete := make(chan struct{}, b.N*1)

	for i := 0; i < b.N; i++ {
		id := int64(i)
		go func() {
			uri := kindUri.ToAggregateRootUri(id)
			for index := 0; index < count; index++ {
				appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
				<-appendComplete
			}
			events := make(chan *goes.EventStoreEntry, 1)
			readComplete, _ := eventStore.LoadAll(uri, events)
			for index := 0; index < count; index++ {
				<-events
			}
			<-readComplete
			roundComplete <- struct{}{}
		}()
	}
	for i := 0; i < b.N; i++ {
		<-roundComplete
	}
}
