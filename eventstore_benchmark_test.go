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
		agg := kind.ToAggregateUri(10000)
		agg.Hash()
	}
}

func Get_EventStoreEntry(size int) *goes.EventStoreEntry {
	data := make([]byte, 10)
	for index, _ := range data {
		data[index] = byte(index)
	}
	return goes.NewEventStoreEntry(10, 1, 1, data)
}

func Run_AppendOnlySync(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kindUri.ToAggregateUri(int64(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			eventStore.Append(uris[i], eventStoreEntry)
		}
	}
}

func Run_ReadOnlySync(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kindUri.ToAggregateUri(int64(i))
	}

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			eventStore.Append(uris[i], eventStoreEntry)
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		//ix := i
		events := make(chan *goes.EventStoreEntry, eventCount)
		eventStore.LoadAll(uris[i], events)
	}
}

func Run_AppendAndReadAllSync(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kindUri.ToAggregateUri(int64(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			eventStore.Append(uris[i], eventStoreEntry)
		}
		events := make(chan *goes.EventStoreEntry, eventCount)
		eventStore.LoadAll(uris[i], events)
	}
}

func Run_AppendOnlyAsync(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kindUri.ToAggregateUri(int64(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.AppendAsync(uris[i], eventStoreEntry)
			<-appendComplete
		}
	}
}

func Run_ReadOnlyAsync(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kindUri.ToAggregateUri(int64(i))
	}

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.AppendAsync(uris[i], eventStoreEntry)
			<-appendComplete
		}
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		events := make(chan *goes.EventStoreEntry, eventCount)
		readComplete, _ := eventStore.LoadAllAsync(uris[i], events)
		<-readComplete
	}
}

func Run_AppendAndReadAllAsync(b *testing.B, connString string, namespace string, kind string, eventSize int, eventCount int) {
	b.StopTimer()
	eventStore, _ := goes.Connect(connString)
	kindUri := goes.NewAggregateKind(namespace, kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	uris := make([]*goes.AggregateUri, b.N)
	for i := 0; i < b.N; i++ {
		uris[i] = kindUri.ToAggregateUri(int64(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			appendComplete, _ := eventStore.AppendAsync(uris[i], eventStoreEntry)
			<-appendComplete
		}
		events := make(chan *goes.EventStoreEntry, eventCount)
		readComplete, _ := eventStore.LoadAllAsync(uris[i], events)
		<-readComplete
	}
}
