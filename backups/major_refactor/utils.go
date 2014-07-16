package eventstore

import "C"
import (
	"hash/crc32"
	"log"
	"math/rand"
	"time"
)

//export DebugPrintf
func DebugPrintf(format *C.char) {
	log.Printf(C.GoString(format))
}

var table *crc32.Table = crc32.MakeTable(crc32.Castagnoli)

// https://github.com/basho/bitcask/blob/master/c_src/murmurhash.c
func MakeCRC(data []byte) uint32 {
	return crc32.Checksum(data, table)
}

func NewKey() int64 {
	return keyGen2()
}

var keyGen = func() func() int64 {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func() int64 {
		return rnd.Int63()
	}
}()

var keyGen2 = func() func() int64 {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	value := rnd.Int63()
	count := 0
	return func() int64 {
		if count == 6 {
			value = rnd.Int63()
			count = 0
		} else {
			value++
			count++
		}
		return value
	}
}()

const (
	MaxUint   = ^uint(0)
	MinUint   = 0
	MaxInt    = int(^uint(0) >> 1)
	MinInt    = -(MaxInt - 1)
	MaxUint16 = ^uint16(0)
	MinUint16 = 0
	MaxInt16  = int16(^uint16(0) >> 1)
	MinInt16  = -(MaxInt16 - 1)
	MaxUint32 = ^uint32(0)
	MinUint32 = 0
	MaxInt32  = int32(^uint32(0) >> 1)
	MinInt32  = -(MaxInt - 1)
	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int64(^uint64(0) >> 1)
	MinInt64  = -(MaxInt - 1)
)
