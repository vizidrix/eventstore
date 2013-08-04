package eventstore_test

import (
	goes "github.com/vizidrix/eventstore"
	"log"
	"runtime"
	"testing"
	"time"
)

func ignore_eventset_benchmarks() {
	log.Printf("")
	//runtime.GOMAXPROCS(10)
	time.Sleep(10)
}

//const gcTime = 10000

//var gcTimer int

type KeyGen func() uint64

func Run_PutGet(b *testing.B, eventSize int, batchSize int, batchCount int) {
	runtime.GC()
	eventSet := goes.NewEmptyEventSet()
	eventData := make([]byte, eventSize)
	for i := 0; i < eventSize; i++ {
		eventData[i] = byte(i | 0xFF)
	}
	batch := make([]goes.Event, batchSize)
	for i := 0; i < batchSize; i++ {
		batch[i] = goes.Event{
			EventType: uint16(i),
			Data:      eventData,
		}
	}
	b.ResetTimer()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		gcTimer++
		if gcTimer == gcTime {
			runtime.GC()
			gcTimer = 0
		}

		b.StartTimer()

		for index := 0; index < batchCount; index++ {
			//b.StartTimer()
			eventSet.Put(batch...)
			//b.StopTimer()

			//b.StartTimer()
			eventSet.Get()
			//b.StopTimer()
		}

		b.StopTimer()

	}
}

func Benchmark_EventSet_PutGet_10byte_b01xc00001(b *testing.B) {
	Run_PutGet(b, 10, 1, 1)
}

func Benchmark_EventSet_PutGet_10byte_b01xc00010(b *testing.B) {
	Run_PutGet(b, 10, 1, 10)
}

func Benchmark_EventSet_PutGet_10byte_b01xc00100(b *testing.B) {
	Run_PutGet(b, 10, 1, 100)
}

func Benchmark_EventSet_PutGet_10byte_b01xc01000(b *testing.B) {
	Run_PutGet(b, 10, 1, 1000)
}

func Benchmark_EventSet_PutGet_10byte_b01xc10000(b *testing.B) {
	Run_PutGet(b, 10, 1, 10000)
}

func Benchmark_EventSet_PutGet_10byte_b10xc00001(b *testing.B) {
	Run_PutGet(b, 10, 10, 1)
}

func Benchmark_EventSet_PutGet_10byte_b10xc00010(b *testing.B) {
	Run_PutGet(b, 10, 10, 10)
}

func Benchmark_EventSet_PutGet_10byte_b10xc00100(b *testing.B) {
	Run_PutGet(b, 10, 10, 100)
}

func Benchmark_EventSet_PutGet_10byte_b10xc01000(b *testing.B) {
	Run_PutGet(b, 10, 10, 1000)
}

func Benchmark_EventSet_PutGet_10byte_b10xc10000(b *testing.B) {
	Run_PutGet(b, 10, 10, 10000)
}

func Benchmark_EventSet_PutGet_1024byte_b01xc00001(b *testing.B) {
	Run_PutGet(b, 1024, 1, 1)
}

func Benchmark_EventSet_PutGet_1024byte_b01xc00010(b *testing.B) {
	Run_PutGet(b, 1024, 1, 10)
}

func Benchmark_EventSet_PutGet_1024byte_b01xc00100(b *testing.B) {
	Run_PutGet(b, 1024, 1, 100)
}

func Benchmark_EventSet_PutGet_1024byte_b01xc01000(b *testing.B) {
	Run_PutGet(b, 1024, 1, 1000)
}

func Benchmark_EventSet_PutGet_1024byte_b01xc10000(b *testing.B) {
	Run_PutGet(b, 1024, 1, 10000)
}

func Benchmark_EventSet_PutGet_1024byte_b10xc00001(b *testing.B) {
	Run_PutGet(b, 1024, 10, 1)
}

func Benchmark_EventSet_PutGet_1024byte_b10xc00010(b *testing.B) {
	Run_PutGet(b, 1024, 10, 10)
}

func Benchmark_EventSet_PutGet_1024byte_b10xc00100(b *testing.B) {
	Run_PutGet(b, 1024, 10, 100)
}

func Benchmark_EventSet_PutGet_1024byte_b10xc01000(b *testing.B) {
	Run_PutGet(b, 1024, 10, 1000)
}

func Benchmark_EventSet_PutGet_1024byte_b10xc10000(b *testing.B) {
	Run_PutGet(b, 1024, 10, 10000)
}

func Benchmark_EventSet_PutGet_4096byte_b01xc00001(b *testing.B) {
	Run_PutGet(b, 4096, 1, 1)
}

func Benchmark_EventSet_PutGet_4096byte_b01xc00010(b *testing.B) {
	Run_PutGet(b, 4096, 1, 10)
}

func Benchmark_EventSet_PutGet_4096byte_b01xc00100(b *testing.B) {
	Run_PutGet(b, 4096, 1, 100)
}

func Benchmark_EventSet_PutGet_4096byte_b01xc01000(b *testing.B) {
	Run_PutGet(b, 4096, 1, 1000)
}

func Benchmark_EventSet_PutGet_4096byte_b01xc10000(b *testing.B) {
	Run_PutGet(b, 4096, 1, 10000)
}

func Benchmark_EventSet_PutGet_4096byte_b10xc00001(b *testing.B) {
	Run_PutGet(b, 4096, 10, 1)
}

func Benchmark_EventSet_PutGet_4096byte_b10xc00010(b *testing.B) {
	Run_PutGet(b, 4096, 10, 10)
}

func Benchmark_EventSet_PutGet_4096byte_b10xc00100(b *testing.B) {
	Run_PutGet(b, 4096, 10, 100)
}

func Benchmark_EventSet_PutGet_4096byte_b10xc01000(b *testing.B) {
	Run_PutGet(b, 4096, 10, 1000)
}

func Benchmark_EventSet_PutGet_4096byte_b10xc10000(b *testing.B) {
	Run_PutGet(b, 4096, 10, 10000)
}
