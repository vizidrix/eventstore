package eventstore

import (
	"fmt"
	//"hash/crc32"
)

type AggregateUri struct {
	AggregateKind
	//namespace string
	//kind      string
	//hash      uint32
	id int64
}

func NewAggregateUri(namespace string, kind string, id int64) *AggregateUri {
	return &AggregateUri{
		AggregateKind: *NewAggregateKind(namespace, kind),
		//namespace: namespace,
		//kind:      kind,
		//hash:      hash,
		//kindHash:  crc32.Checksum([]byte(namespace+kind), crc32.MakeTable(crc32.Castagnoli)),
		id: id,
	}
}

func (uri *AggregateUri) Namespace() string {
	return uri.namespace
}

func (uri *AggregateUri) Kind() string {
	return uri.kind
}

func (uri *AggregateUri) Hash() uint32 {
	return uri.hash
}

func (uri *AggregateUri) Id() int64 {
	return uri.id
}

func (uri *AggregateUri) RelativePath() string {
	return fmt.Sprint("%s/%s/%d",
		uri.Namespace,
		uri.Kind,
		uri.Id)
}
