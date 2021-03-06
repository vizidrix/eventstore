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

func Run_PutGet_Spread(b *testing.B, eventSize int, batchSize int, batchCount int) {
	runtime.GC()
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
		//go func() {
		//b.StopTimer()
		/*
			gcTimer++
			if gcTimer == gcTime {
				//b.StopTimer()
				runtime.GC()
				//b.StartTimer()
				gcTimer = 0
			}
		*/

		store := make([]*goes.EventSet, batchCount)
		for index := 0; index < batchCount; index++ {
			store[index] = goes.NewEmptyEventSet()
		}

		b.StartTimer()

		for index := 0; index < batchCount; index++ {
			for i := 0; i < batchCount; i++ {
				store[index], _ = store[index].Put(batch...)
			}
		}

		for index := 0; index < batchCount; index++ {
			store[index].GetSlice(0, 1)
		}

		b.StopTimer()
		//}()
		//count++
	}
	/*
		for count < b.N {
			time.Sleep(10 * time.Nanosecond)
		}
	*/
}

func Run_PutGet_Trim(b *testing.B, eventSize int) {
	//b.StopTimer()
	runtime.GC()
	event := goes.Event{
		EventType: 1,
		Data:      make([]byte, eventSize),
	}
	/*
		eventSets := make([]*goes.EventSet, b.N)
		for i := 0; i < b.N; i++ {
			eventSets[i] = goes.NewEmptyEventSet()
		}
	*/
	//eventSet := goes.NewEmptyEventSet()
	b.ResetTimer()
	b.StopTimer()
	for index := 0; index < b.N; index++ {
		//b.StopTimer()
		eventSet := goes.NewEmptyEventSet()

		//for i := 0; i < 10; i++ {
		////eventSets[index], _ = eventSets[index].Put(event)
		//}
		//for i := 0; i < 10; i++ {
		////eventSets[index].Get()
		//}

		//b.StartTimer()
		for i := 0; i < 20; i++ {
			eventSet, _ = eventSet.Put(event) //, event, event)
		}
		b.StartTimer()

		for i := 0; i < 100; i++ {
			eventSet.GetSlice(5, 15)
		}
		b.StopTimer()
	}
}

func Run_PutGet_Trim2(b *testing.B, eventSize int) {
	//b.StopTimer()
	//runtime.GC()
	event := goes.Event{
		EventType: 1,
		Data:      make([]byte, eventSize),
	}
	/*
		eventSets := make([]*goes.EventSet, b.N)
		for i := 0; i < b.N; i++ {
			eventSets[i] = goes.NewEmptyEventSet()
		}
	*/
	//eventSet := goes.NewEmptyEventSet()
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		eventSet := goes.NewEmptyEventSet()
		//for i := 0; i < 10; i++ {
		////eventSets[index], _ = eventSets[index].Put(event)
		//}
		//for i := 0; i < 10; i++ {
		////eventSets[index].Get()
		//}

		eventSet, _ = eventSet.Put(event)
		eventSet.GetSlice(0, 1)
	}
}

func Run_PutGet(b *testing.B, eventSize int, batchSize int, batchCount int) {
	runtime.GC()
	//eventSet := goes.NewEmptyEventSet()
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

	//b.ResetTimer()
	//b.StopTimer()
	for i := 0; i < b.N; i++ {
		/*
			gcTimer++
			if gcTimer == gcTime {
				//b.StopTimer()
				runtime.GC()
				//b.StartTimer()
				gcTimer = 0
			}
		*/

		eventSet := goes.NewEmptyEventSet()
		//b.StartTimer()

		for index := 0; index < batchCount; index++ {
			eventSet, _ = eventSet.Put(batch...)
			//eventSet.Put(batch...)
		}

		eventSet.GetSlice(0, 1)
		//b.StopTimer()
	}
}

func Run_PutGet2(b *testing.B, eventSize int, batchSize int, batchCount int) {
	runtime.GC()
	//eventSet := goes.NewEmptyEventSet()
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
	var err error
	for i := 0; i < b.N; i++ {
		gcTimer++
		if gcTimer == gcTime {
			//b.StopTimer()
			runtime.GC()
			//b.StartTimer()
			gcTimer = 0
		}

		eventSet := goes.NewEmptyEventSet()
		for index := 0; index < batchCount; index++ {
			//b.StartTimer()

			eventSet, err = eventSet.Put(batch...)
			if err != nil {
				log.Printf("Put failed: %s", err)
				b.Fail()
			}

			//b.StopTimer()
		}

		b.StartTimer()
		_, err := eventSet.GetSlice(0, 1)
		//b.Logf("Read: %d", len(events))
		b.StopTimer()
		if err != nil {
			b.Fail()
		}
		/*
			if err := eventSet.CheckSum(); err != nil {
				b.Fail()
			}
		*/

		/*
			if _, err := eventSet.Get(); err == nil {
				//if len(events) != batchSize*batchCount {
				//	log.Printf("Events[ %d, %d ]: %d, %d, %d", len(events), cap(events), i, batchSize, batchCount)
				//	b.Fail()
				//}
			} else {
				log.Printf("Put failed: %s", err)
				b.Fail()
			}
		*/

	}
}

func Benchmark_EventSet_PutGetTrim_10byte(b *testing.B) {
	Run_PutGet_Trim(b, 10)
}

func Benchmark_EventSet_PutGetTrim_1024byte(b *testing.B) {
	Run_PutGet_Trim(b, 1024)
}

func Benchmark_EventSet_PutGetTrim_4096byte(b *testing.B) {
	Run_PutGet_Trim(b, 4096)
}

func Benchmark_EventSet_PutGetTrim_16384byte(b *testing.B) {
	Run_PutGet_Trim(b, 16384)
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
