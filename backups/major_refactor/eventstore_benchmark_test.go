package eventstore_test

import (
	"runtime"
	//"encoding/binary"
	//"bytes"
	//goes "github.com/vizidrix/eventstore"
	//. "github.com/vizidrix/eventstore/test_utils"
	//"hash/crc32"
	"log"
	//"math/rand"
	//"testing"
	"time"
)

func ignore_eventstore_benchmarks() {
	log.Printf("")
	runtime.GOMAXPROCS(10)
	time.Sleep(10)
}

const gcTime = 10000

var gcTimer int

//func Run_ES_PutGet(b *testing.B, )
/*
gcTimer++
if gcTimer == gcTime {
	runtime.GC()
	gcTimer = 0
}
*/

/*********

func Get_EventStoreEntry(size int) *goes.EventStoreEntry {
	data := make([]byte, size)
	for index, _ := range data {
		data[index] = byte(index)
	}
	return goes.NewEventStoreEntry(uint16(size), 1, 1, data)
}

func Run_PutOnly(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	//log.Printf("Count: %d", b.N)
	b.StopTimer()
	//eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	//partition := eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	i := 0
	index := 0
	runtime.GC()
	//b.ResetTimer()
	//b.StopTimer()
	for i = 0; i < b.N; i++ {
		eventStore, _ := goes.Connect(connString)
		kindPartition := eventStore.Kind(kind)
		aggregatePartition := kindPartition.Aggregate(int64(i))
		//id := int64(i)
		gcTimer++
		if gcTimer == gcTime {
			runtime.GC()
			gcTimer = 0
			//log.Printf("GC Flush: %d", i)
		}
		b.StartTimer()
		for index = 0; index < eventCount; index++ {
			aggregatePartition.Put(eventStoreEntry)
		}
		b.StopTimer()
	}
}

func Run_GetOnly(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	kind := goes.NewAggregateKind(namespace, kindName)
	eventStoreEntry := Get_EventStoreEntry(eventSize)

	i := 0
	runtime.GC()
	b.ResetTimer()
	b.StopTimer()
	for i = 0; i < b.N; i++ {
		eventStore, _ := goes.Connect(connString)
		kindPartition := eventStore.Kind(kind)
		aggregatePartition := kindPartition.Aggregate(int64(i))
		for index := 0; index < eventCount; index++ {
			aggregatePartition.Put(eventStoreEntry)
		}
		gcTimer++
		if gcTimer == gcTime {
			runtime.GC()
			gcTimer = 0
		}
		b.StartTimer()
		aggregatePartition.Get()
		b.StopTimer()
	}
	//log.Printf("Done iteration: %d", b.N)
}

func Run_PutGet(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	//b.StopTimer()
	//eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	//partition := eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	i := 0
	index := 0
	runtime.GC()
	b.ResetTimer()
	b.StopTimer()
	//b.StartTimer()
	for i = 0; i < b.N; i++ {
		eventStore, _ := goes.Connect(connString)
		//partition := eventStore.RegisterKind(kind)
		kindPartition := eventStore.Kind(kind)
		aggregatePartition := kindPartition.Aggregate(int64(i))

		//events := make(chan *goes.EventStoreEntry, eventCount)
		//id := int64(i)
		gcTimer++
		if gcTimer == gcTime {
			runtime.GC()
			gcTimer = 0
		}
		b.StartTimer()
		for index = 0; index < eventCount; index++ {
			aggregatePartition.Put(eventStoreEntry)
		}
		aggregatePartition.Get()
		b.StopTimer()
	}
}


*********/

/*
func Run_AppendOnlyAsync(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kind.ToAggregateUri(int64(i))
	}

	for i := 0; i < b.N; i++ {
		index := 0
		b.StartTimer()
		for index = 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.AppendAsync(uris[i], eventStoreEntry)
			<-appendComplete
		}
		b.StopTimer()
	}
}

func Run_ReadOnlyAsync(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kind.ToAggregateUri(int64(i))
	}

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.AppendAsync(uris[i], eventStoreEntry)
			<-appendComplete
		}
	}

	for i := 0; i < b.N; i++ {
		events := make(chan *goes.EventStoreEntry, eventCount)
		b.StartTimer()
		readComplete, _ := eventStore.LoadAllAsync(uris[i], events)
		<-readComplete
		b.StopTimer()
	}
}

func Run_AppendAndReadAllAsync(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kind.ToAggregateUri(int64(i))
	}

	for i := 0; i < b.N; i++ {
		index := 0
		events := make(chan *goes.EventStoreEntry, eventCount)
		b.StartTimer()
		for index = 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.AppendAsync(uris[i], eventStoreEntry)
			<-appendComplete
		}
		readComplete, _ := eventStore.LoadAllAsync(uris[i], events)
		<-readComplete
		b.StopTimer()
	}
}
*/

/*
package main

import (
	"runtime"
	"testing"
)
*/

/*
func Benchmark_Create_Serialize_DeSerialize_EventStoreEntry_10bytePayload(b *testing.B) {
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		entry := goes.NewEventStoreEntryFrom(1, data)
		header, body := entry.ToBinary()
		goes.FromBinary(header, body)
	}
}

func Benchmark_Create_Serialize_DeSerialize_EventStoreEntry_4087bytePayload(b *testing.B) {
	data := make([]byte, 4087)
	for index, _ := range data {
		data[index] = byte(index)
	}
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		entry := goes.NewEventStoreEntryFrom(1, data)
		header, body := entry.ToBinary()
		goes.FromBinary(header, body)
	}
}

func Benchmark_Create_EventStoreEntry_10bytePayload(b *testing.B) {
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		goes.NewEventStoreEntryFrom(1, data)
	}
}

func Benchmark_Create_EventStoreEntry_4087bytePayload(b *testing.B) {
	data := make([]byte, 4087)
	for index, _ := range data {
		data[index] = byte(index)
	}
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		goes.NewEventStoreEntryFrom(1, data)
	}
}

func Benchmark_CreateRaw_EventStoreEntry_10bytePayload(b *testing.B) {
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		goes.NewEventStoreEntry(10, 1, 1, data)
	}
}

func Benchmark_CreateRaw_EventStoreEntry_4087bytePayload(b *testing.B) {
	data := make([]byte, 4087)
	for index, _ := range data {
		data[index] = byte(index)
	}
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		goes.NewEventStoreEntry(4087, 1, 1, data)
	}
}

func Benchmark_KindUri_to_AggregateRootUri(b *testing.B) {
	kind := goes.NewAggregateKind("namespace", "kind")
	i := 0
	b.ResetTimer()

	for i = 0; i < b.N; i++ {
		agg := kind.ToAggregateUri(10000)
		agg.Hash()
	}
}
*/
