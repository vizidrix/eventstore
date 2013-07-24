package eventstore

import (
	"fmt"
)

type AggregateRootUri struct {
	namespace string
	kind      string
	id        int64
}

func NewAggregateRootUri(namespace string, kind string, id int64) *AggregateRootUri {
	return &AggregateRootUri{
		namespace: namespace,
		kind:      kind,
		id:        id,
	}
}

func (uri *AggregateRootUri) Namespace() string {
	return uri.namespace
}

func (uri *AggregateRootUri) Kind() string {
	return uri.kind
}

func (uri *AggregateRootUri) Id() int64 {
	return uri.id
}

func (uri *AggregateRootUri) RelativePath() string {
	return fmt.Sprint("%s/%s/%d",
		uri.Namespace,
		uri.Kind,
		uri.Id)
}

func (uri *AggregateRootUri) ToAggregateKindUri() *AggregateKindUri {
	return NewAggregateKindUri(uri.namespace, uri.kind)
}
