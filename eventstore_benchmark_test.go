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

func Get_EventStoreEntry(size int) *goes.EventStoreEntry {
	data := make([]byte, size)
	for index, _ := range data {
		data[index] = byte(index)
	}
	return goes.NewEventStoreEntry(uint16(size), 1, 1, data)
}

func Run_AppendOnlySync(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	//b.StopTimer()
	//eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	//partition := eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	i := 0
	index := 0
	b.ResetTimer()
	b.StopTimer()
	for i = 0; i < b.N; i++ {
		eventStore, _ := goes.Connect(connString)
		partition := eventStore.RegisterKind(kind)
		id := int64(i)
		//b.StartTimer()
		for index = 0; index < eventCount; index++ {

			b.StartTimer()
			partition.Append(id, eventStoreEntry)
			b.StopTimer()
		}
		//b.StopTimer()
	}
}

func Run_ReadOnlySync(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	//b.StopTimer()
	//b.StopTimer()
	//eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	//partition := eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	/*for i := 0; i < b.N; i++ {
		for index := 0; index < eventCount; index++ {
			partition.Append(int64(i), eventStoreEntry)
		}
	}*/
	i := 0
	//b.StartTimer()
	//b.StopTimer()b.StartTimer()
	//log.Printf("Doing iteration: %d", b.N)
	b.ResetTimer()
	b.StopTimer()
	for i = 0; i < b.N; i++ {
		eventStore, _ := goes.Connect(connString)
		partition := eventStore.RegisterKind(kind)
		for index := 0; index < eventCount; index++ {
			partition.Append(int64(i), eventStoreEntry)
		}
		//events := make(chan *goes.EventStoreEntry, eventCount)
		id := int64(i)
		b.StartTimer()
		//log.Printf("Doing iteration: %d", i)

		//b.StartTimer()
		partition.LoadAll(id)
		b.StopTimer()
	}
	//log.Printf("Done iteration: %d", b.N)
}

func Run_AppendAndReadAllSync(b *testing.B, connString string, namespace string, kindName string, eventSize int, eventCount int) {
	//b.StopTimer()
	//eventStore, _ := goes.Connect(connString)
	kind := goes.NewAggregateKind(namespace, kindName)
	//partition := eventStore.RegisterKind(kind)
	eventStoreEntry := Get_EventStoreEntry(eventSize)
	i := 0
	index := 0
	b.ResetTimer()
	b.StopTimer()
	//b.StartTimer()
	for i = 0; i < b.N; i++ {
		eventStore, _ := goes.Connect(connString)
		partition := eventStore.RegisterKind(kind)

		//events := make(chan *goes.EventStoreEntry, eventCount)
		id := int64(i)
		b.StartTimer()
		for index = 0; index < eventCount; index++ {
			partition.Append(id, eventStoreEntry)
		}
		partition.LoadAll(id)
		b.StopTimer()
	}
}

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
