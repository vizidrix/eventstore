package eventstore

import (
	//"errors"
	//"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	//"os"
	"encoding/binary"
	//"hash/crc32"
	"time"
)

func ignore() { log.Println("") }

type EventStore struct {
	path string
}

type Domain struct {
	eventstore *EventStore
	path       string
	namespace  string
}

type Kind struct {
	domain   *Domain
	path     string
	kindName string
}

type Aggregate struct {
	kind *Kind
	path string
	id   int64
}

type IEvent interface {
	ToBinary() ([]byte, error)
}

func Connect(path string) (*EventStore, error) {
	//path := "/eventstore/"
	makeDirectory(path)

	return &EventStore{path}, nil
}

func (es *EventStore) Domain(namespace string) (*Domain, error) {
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
	path := fmt.Sprintf("%s%d.aggregate", kind.path, id)

	makeFile(path)
	return &Aggregate{
		kind: kind,
		path: path,
		id:   id,
	}, nil
}

func (aggregate *Aggregate) Append(event IEvent) error {
	data, _ := event.ToBinary()
	log.Printf("%d data", data)
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, int32(len(data)))
	buffer.Write(data)
	log.Printf("Buffer: %s", buffer)
	file, _ := getFile(aggregate.path)
	defer file.Close()
	//sz, _ := file.Write(buffer.Bytes())
	sz, _ := file.Write(data)

	log.Printf("%d bytes written", sz)

	return nil
}

func getFile(path string) (IFile, error) {
	var fs IFileStore = osFileStore{}

	if file, err := fs.Open(path); err == nil {
		return file, err
	} else {
		return nil, err
	}
}

func makeFile(path string) (IFile, error) {
	var fs IFileStore = osFileStore{}

	if _, err := fs.Stat(path); err != nil {
		if fs.IsNotExist(err) {
			log.Println("File doesn't exist")
			if file, err := fs.Create(path); err != nil {
				log.Println("Unable to create file: %s", err)
				return file, err
			}
		} else {
			log.Println("Unknown error")
			return nil, err
		}
	}
	log.Printf("File exists: %s", path)
	return nil, nil
}

func makeDirectory(path string) (IFile, error) {
	var fs IFileStore = osFileStore{}

	if _, err := fs.Stat(path); err != nil {
		if fs.IsNotExist(err) {
			// file does not exist
			log.Println("Path doesn't exist")
			//var perm uint32 = 0755
			if err := fs.Mkdir(path, 0755); err != nil {
				log.Println("Unable to create dir: %s", err)
				return nil, err
			}

		} else {
			// other error
			log.Println("Unknown error")
			return nil, err
		}
	}
	log.Printf("Directory exists: %s", path)
	return nil, nil
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
