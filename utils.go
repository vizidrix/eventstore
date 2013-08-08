package eventstore

import (
	"hash/crc32"
	"math/rand"
	"time"
)

var table *crc32.Table = crc32.MakeTable(crc32.Castagnoli)

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
