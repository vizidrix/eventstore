package eventstore

import (
	"fmt"
	//"log"
)

type AggregateKindUri struct {
	namespace string
	kind      string
}

func NewAggregateKindUri(namespace string, kind string) *AggregateKindUri {
	return &AggregateKindUri{
		namespace: namespace,
		kind:      kind,
	}
}

func (uri *AggregateKindUri) Namespace() string {
	return uri.namespace
}

func (uri *AggregateKindUri) Kind() string {
	return uri.kind
}

func (uri *AggregateKindUri) RelativePath() string {
	return fmt.Sprint("%s/%s",
		uri.Namespace,
		uri.Kind)
}

func (uri *AggregateKindUri) ToAggregateRootUri(id int64) *AggregateRootUri {
	return NewAggregateRootUri(uri.namespace, uri.kind, id)
}
