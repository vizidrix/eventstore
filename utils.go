package eventstore

import (
	"hash/crc32"
)

var table *crc32.Table = crc32.MakeTable(crc32.Castagnoli)

func MakeCRC(data []byte) uint32 {
	return crc32.Checksum(data, table)
}
