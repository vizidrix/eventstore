package eventstore

import (
	"fmt"
	//"hash/crc32"
	"log"
)

func ignore_aggregatekind() { log.Printf("") }

type AggregateKind struct {
	namespace string
	kind      string
	hash      uint32
}

func NewAggregateKind(namespace string, kind string) *AggregateKind {
	//hash := crc32.Checksum([]byte(namespace+kind), crc32.MakeTable(crc32.Castagnoli))
	return &AggregateKind{
		namespace: namespace,
		kind:      kind,
		hash:      MakeCRC([]byte(namespace + kind)),
	}
}

func (kind *AggregateKind) Namespace() string {
	return kind.namespace
}

func (kind *AggregateKind) Kind() string {
	return kind.kind
}

func (kind *AggregateKind) Hash() uint32 {
	return kind.hash
}

func (kind *AggregateKind) KindPath() string {
	return fmt.Sprint("%s/%s",
		kind.Namespace,
		kind.Kind)
}

func (kind *AggregateKind) ToAggregateUri(id int64) *AggregateUri {
	return &AggregateUri{
		AggregateKind: *kind,
		id:            id,
	}
}
