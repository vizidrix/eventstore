package eventstore

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func event_store_ignore() { log.Println(fmt.Sprintf("", 10)) }

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

// http://graphics.stanford.edu/~seander/bithacks.html
//func PowerOf2(value uint64) bool {
//	return value && !(value & (value - 1))
//}

type EventReader interface {
	Get() (*EventSet, error)
	GetSlice(startIndex int, endIndex int) (*EventSet, error)
}

type EventWriter interface {
	Put(newEvents ...Event) (*EventSet, error)
}

type EventStorer interface {
	Kind(kind *AggregateKind) KindPartitioner
}

type KindPartitioner interface {
	Id(id uint64) AggregatePartitioner
}

type AggregatePartitioner interface {
	EventReader
	EventWriter
}

func Connect(connString string) (EventStorer, error) {
	if strings.HasPrefix(connString, "fs://") {
		return NewFileSystemES(connString), nil
	} else if strings.HasPrefix(connString, "mem://") {
		return NewMemoryES(connString), nil
	} else {
		return nil, errors.New("Unable to find delimiter in connection string")
	}
}

/*
type ReadEventStorer interface {
	// Returns an array of all EventStoreEntry's for the aggregate uri
	//LoadAll(id int64, entries chan<- *EventStoreEntry) error
	LoadAll(id int64) ([]*EventStoreEntry, error)
	// Reutrns an array of all EventStoreEntry's for the aggregate uri that were between the start and end index range
	//LoadIndexRange(id int64, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) error
	LoadIndexRange(id int64, startIndex uint64, endIndex uint64) ([]*EventStoreEntry, error)
}
*/
/*
type AsyncReadEventStorer interface {
	// Returns an array of all EventStoreEntry's for the aggregate uri
	LoadAllAsync(uri *AggregateUri, entries chan<- *EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error)
	// Reutrns an array of all EventStoreEntry's for the aggregate uri that were between the start and end index range
	LoadIndexRangeAsync(uri *AggregateUri, entries chan<- *EventStoreEntry, startIndex uint64, endIndex uint64) (completeChan <-chan struct{}, errorChan <-chan error)
}
*/
/*
type WriteEventStorer interface {
	Append(id int64, entry *EventStoreEntry) error
}
*/

/*
type AsyncWriteEventStorer interface {
	AppendAsync(uri *AggregateUri, entries ...*EventStoreEntry) (completeChan <-chan struct{}, errorChan <-chan error)
}
*/

/*
type SyncEventStorer interface {
	SyncReadEventStorer
	SyncWriteEventStorer
}
*/
/*
type AsyncEventStorer interface {
	AsyncReadEventStorer
	AsyncWriteEventStorer
}
*/

/* -- Replaced by AggregateKind and the -> AggregateUri ability */
/*
type Domain struct {
	path      string
	namespace string
}

type Kind struct {
	path     string
	kindName string
}

type Aggregate struct {
	path string
	id   int64
}

type IEvent interface {
	ToBinary() ([]byte, error)
}
*/
/*
func (es EventStorer) Domain(namespace string) (*Domain, error) {
	path := fmt.Sprintf("%s%s/", es.path, namespace)
	makeDirectory(path)

	return &Domain{
		eventstore: es,
		path:       path,
		namespace:  namespace,
	}, nil
}

func (domain *Domain) Kind(kindName string) (*Kind, error) {
	path := fmt.Sprintf("%s%s/", domain.path, kindName)

	makeDirectory(path)
	return &Kind{
		domain:   domain,
		path:     path,
		kindName: kindName,
	}, nil
}

func (kind *Kind) Aggregate(id int64) (*Aggregate, error) {
	path := fmt.Sprintf("%s%d.agg", kind.path, id)

	makeFile(path)
	return &Aggregate{
		kind: kind,
		path: path,
		id:   id,
	}, nil
}

/*
func (aggregate *Aggregate) LoadAll() ([]IEvent, error) {

}
*/

/*
func (aggregate *Aggregate) Append(event IEvent) error {
	data, _ := event.ToBinary()
	log.Printf("Got data: %d", data)

	file, err := os.OpenFile(aggregate.path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error opening file: %s error %s", aggregate.path, err)
		return err
	}

	count, err := file.Write(data)
	if err != nil {
		log.Printf("Error writing to file: %s with %s error %s", aggregate.path, data, err)
		return err
	}

	log.Printf("Should have written %d bytes to file: %s", count, aggregate.path)
	return nil
}
*/
/*
	//file, err := os.Open(aggregate.path)
	file, err := makeFile(aggregate.path)
	//defer file.Close()
	if err != nil {
		log.Printf("Error opening file: %s", aggregate.path)
		return err
	}

	count, err := file.Write(data)
	file.Close()

	if err != nil {
		log.Printf("Error writing to file: %s", err)
		return err
	}
	log.Printf("Wrote %d bytes to file %s", count, aggregate.path)
	//buffer := new(bytes.Buffer)
	//binary.Write(buffer, binary.BigEndian, int32(len(data)))
	//buffer.Write(data)
	//log.Printf("Buffer: %s", buffer)
	//file, _ := getFile(aggregate.path)

	//WriteData(file, data)

	//sz, _ := file.Write(buffer.Bytes())
	//sz, _ := file.Write(data)

	//log.Printf("%d bytes written", sz)

	return nil
*/
