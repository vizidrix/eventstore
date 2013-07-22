package eventstore

import (
	"io"
	"os"
)

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
