package eventstore

import (
	"fmt"
)

type AggregateKindURI struct {
	namespace string
	kind      string
}

func NewAggregateKindURI(namespace string, kind string) *AggregateKindURI {
	return &AggregateKindURI{
		namespace: namespace,
		kind:      kind,
	}
}

func (uri *AggregateKindURI) Namespace() string {
	return uri.namespace
}

func (uri *AggregateKindURI) Kind() string {
	return uri.kind
}

func (uri *AggregateKindURI) RelativePath() string {
	return fmt.Sprint("%s/%s",
		uri.Namespace,
		uri.Kind)
}

func (uri *AggregateKindURI) ToAggregateRootURI(id int64) *AggregateRootURI {
	return NewAggregateRootURI(uri.namespace, uri.kind, id)
}
