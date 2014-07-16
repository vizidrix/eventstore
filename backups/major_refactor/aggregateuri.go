package eventstore

import (
	"fmt"
)

type AggregateUri struct {
	AggregateKind
	id int64
}

func NewAggregateUri(namespace string, kind string, id int64) *AggregateUri {
	return &AggregateUri{
		AggregateKind: *NewAggregateKind(namespace, kind),
		id:            id,
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
