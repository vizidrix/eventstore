package eventstore

/*
//#cgo CFLAGS: -I/usr/include/c-client-src
//#cgo LDFLAGS: -lgodb_mt
// From: https://code.google.com/p/go-sqlite/source/browse/go1/sqlite3/sqlite3.go
//#cgo CFLAGS: -Os
//#cgo CFLAGS: -DNDEBUG=1
//#cgo CFLAGS: -DSQLITE_THREADSAFE=2
//#cgo CFLAGS: -DSQLITE_TEMP_STORE=2
//#cgo CFLAGS: -DSQLITE_USE_URI=1
//#cgo CFLAGS: -DSQLITE_ENABLE_FTS3_PARENTHESIS=1
//#cgo CFLAGS: -DSQLITE_ENABLE_FTS4=1
//#cgo CFLAGS: -DSQLITE_ENABLE_RTREE=1
//#cgo CFLAGS: -DSQLITE_ENABLE_STAT3=1
//#cgo CFLAGS: -DSQLITE_SOUNDEX=1
//#cgo CFLAGS: -DSQLITE_OMIT_AUTHORIZATION=1
//#cgo CFLAGS: -DSQLITE_OMIT_AUTOINIT=1
//#cgo CFLAGS: -DSQLITE_OMIT_LOAD_EXTENSION=1
//#cgo CFLAGS: -DSQLITE_OMIT_TRACE=1
//#cgo CFLAGS: -DSQLITE_OMIT_UTF16=1

#include "godb.h"

//# include "godb.cgo"
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	//"reflect"
	"unsafe"
)

func godb_ignore() {
	log.Printf(fmt.Sprintf("", errors.New("")))
	temp := unsafe.Pointer(&struct{}{})
	log.Printf("% v", temp)
	result := C.temp(10)
	log.Printf("Result: %d ", result)

}

func Stuff(value int) (int, error) {
	result := C.temp(C.int(value))
	log.Printf("Result: %d ", result)

	str := C.CString("/go/datafile.godb")
	//str := C.CString("/datafile.godb")
	defer C.free(unsafe.Pointer(str))
	var handle C.int
	handle, err := C.godb_open_file(str, 0666)
	if err != nil {
		log.Printf("Error opening file: %s", err)
		return int(handle), err
	}
	//defer C.free(unsafe.Pointer(handle))
	log.Printf("File handle: %d", handle)

	return int(handle), nil
}
