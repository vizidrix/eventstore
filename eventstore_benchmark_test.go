package eventstore_test

import (
	"runtime"
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

func ignore_eventstore_benchmarks() {
	log.Printf("")
	runtime.GOMAXPROCS(10)
}

func Get_EventStoreEntry(size int) *goes.EventStoreEntry {
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	return goes.NewEventStoreEntry(10, 1, 1, data)
}

func Run_AppendOnly(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateUri(id)
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
	}
}

func Run_ReadOnly(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateUri(id)
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateUri(id)
		events := make(chan *goes.EventStoreEntry, 1) //b.N)
		readComplete, _ := eventStore.LoadAll(uri, events)
		<-readComplete
	}
}

func Run_AppendAndReadAll(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateUri(id)
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
		events := make(chan *goes.EventStoreEntry, eventCount) //1) //b.N)
		readComplete, _ := eventStore.LoadAll(uri, events)
		<-readComplete
	}
}

func Run_AppandAndReadAll_Multiples(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKind, eventStoreEntry *goes.EventStoreEntry, count int) {
	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateUri(id)
		for index := 0; index < count; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
		events := make(chan *goes.EventStoreEntry, b.N)
		readComplete, _ := eventStore.LoadAll(uri, events)
		for index := 0; index < count; index++ {
			<-events
		}
		<-readComplete
	}
}
