package eventstore

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func filestore_ignore() { log.Println(fmt.Sprintf("", 10)) }

//var fs FileStore = osFS{}

type IFileStore interface {
	IsNotExist(err error) bool
	Mkdir(name string, perm os.FileMode) error
	Create(name string) (*os.File, error)
	Open(name string) (*os.File, error)
	Stat(name string) (os.FileInfo, error)
}

type IFile interface {
	io.Closer
	io.Reader
	io.ReaderAt
	//io.ReadCloser
	//io.ReadSeeker
	io.Seeker
	io.Writer
	io.WriterAt
	//io.WriteCloser
	//io.WriteSeeker
	Stat() (os.FileInfo, error)
}

type IDirectory interface {
}

// osFS implements fileSystem using the local disk.
type osFileStore struct{}

func (osFileStore) IsNotExist(err error) bool                 { return os.IsNotExist(err) }
func (osFileStore) Mkdir(name string, perm os.FileMode) error { return os.Mkdir(name, perm) }
func (osFileStore) Create(name string) (*os.File, error)      { return os.Create(name) }
func (osFileStore) Open(name string) (*os.File, error)        { return os.Open(name) }
func (osFileStore) Stat(name string) (os.FileInfo, error)     { return os.Stat(name) }

const (
	Append_WriteOnly          = os.O_APPEND | os.O_WRONLY
	Create_Truncate_WriteOnly = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
)

/*
func Write(file string, flag int, data []byte) {
	if file, err := os.Open(file, flag, 0666); err == nil {
		// File opened successfully
	} else {
		log.Printf("Unable to open file: %s - %s", file, err)
	}

}
*/

func WriteData(dest io.Writer, data []byte) error {
	return binary.Write(dest, binary.BigEndian, data)
}

func getFile(path string) (*os.File, error) {
	var fs IFileStore = osFileStore{}

	if file, err := fs.Open(path); err == nil {
		return file, err
	} else {
		return nil, err
	}
}

func makeFile(path string) (*os.File, error) {
	var fs IFileStore = osFileStore{}

	if _, err := fs.Stat(path); err != nil {
		if fs.IsNotExist(err) {
			log.Println("File doesn't exist")
			if _, err := fs.Create(path); err != nil {
				log.Println("Unable to create file: %s", err)
				return nil, err
			}
		} else {
			log.Println("Unknown error")
			return nil, err
		}
	}
	log.Printf("File exists: %s", path)
	file, err := fs.Open(path)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return nil, err
	}
	return file, nil
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
