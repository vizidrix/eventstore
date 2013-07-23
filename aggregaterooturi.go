package eventstore

import (
	"fmt"
)

type AggregateRootURI struct {
	namespace string
	kind      string
	id        int64
}

func NewAggregateRootURI(namespace string, kind string, id int64) *AggregateRootURI {
	return &AggregateRootURI{
		namespace: namespace,
		kind:      kind,
		id:        id,
	}
}

func (uri *AggregateRootURI) Namespace() string {
	return uri.namespace
}

func (uri *AggregateRootURI) Kind() string {
	return uri.kind
}

func (uri *AggregateRootURI) Id() int64 {
	return uri.id
}

func (uri *AggregateRootURI) RelativePath() string {
	return fmt.Sprint("%s/%s/%d",
		uri.Namespace,
		uri.Kind,
		uri.Id)
}

func (uri *AggregateRootURI) ToAggregateKindURI() *AggregateKindURI {
	return NewAggregateKindURI(uri.namespace, uri.kind)
}
