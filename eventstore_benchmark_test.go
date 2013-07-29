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
	//runtime.GOMAXPROCS(2)
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKindUri(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
	}
}

func Run_ReadOnly(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	//runtime.GOMAXPROCS(2)
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKindUri(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		events := make(chan *goes.EventStoreEntry, 1) //b.N)
		readComplete, _ := eventStore.LoadAll(uri, events)
		<-readComplete
	}
}

func Run_AppendAndReadAll(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry) {
	//runtime.GOMAXPROCS(2)
	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
		<-appendComplete
		events := make(chan *goes.EventStoreEntry, 1) //b.N)
		readComplete, _ := eventStore.LoadAll(uri, events)
		<-readComplete
	}
}

func Run_AppandAndReadAll_Multiples(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry, count int) {
	//runtime.GOMAXPROCS(2)

	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
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

/*
func Run_Async_AppendOnly(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry) {
	appendCompletes := 0
	for i := 0; i < b.N; i++ {
		id := int64(i)
		go func() {
			uri := kindUri.ToAggregateRootUri(id)
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
			appendCompletes++
		}()
	}
	for appendCompletes < b.N {
		//log.Printf("Append completes: %d", appendCompletes)
	}
}

func Run_Async_AppendOnly_Multiples(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry, count int) {
	for i := 0; i < b.N; i++ {
		id := int64(i)
		uri := kindUri.ToAggregateRootUri(id)
		for index := 0; index < count; index++ {
			appendComplete, _ := eventStore.Append(uri, eventStoreEntry)
			<-appendComplete
		}
	}
}

func Run_Async_AppendAndReadAll(b *testing.B, eventStore goes.EventStorer, kindUri *goes.AggregateKindUri, eventStoreEntry *goes.EventStoreEntry) {
	runs := b.N
	roundComplete := make(chan struct{}, runs)

	runCount := 0
	roundCount := 0

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
			runCount++
		}()
	}
	for i := 0; i < runs; i++ {
		<-roundComplete
		roundCount++
	}

	log.Printf("RUN: %d - ROUND: %d", runCount, roundCount)
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
*/
