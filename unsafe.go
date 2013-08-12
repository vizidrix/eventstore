package eventstore

/*
#include <assert.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

void Copy_memory_loop(void* dest, void* src, int length) {
	unsigned char* cdest = (unsigned char*) dest;
	unsigned char* csrc = (unsigned char*) src;

	int i;
	for (i = 0; i < length; i++) {
		cdest[i] = csrc[i];
	}
}

void write_memory_rep_stosq(void* dest, void* src, int length) {
	//unsigned char* cdest = (unsigned char*) dest;
	//unsigned char* csrc = (unsigned char*) src;

	//memset(cdest, csrc, 8);
  asm("cld\n"
      "rep stosq"
      : : "D" (dest), "c" (1), "a" (src) );
}

*/

//import "C"

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

func ignore_unsafe() {
	log.Printf(fmt.Sprintf("", errors.New("")))
}

/**********************************************************

	UNSAFE METHODS: USE CAUTION WHEN OPERATING

**********************************************************/
func UnsafeCopyBytes(dest []byte, src []byte) {
	//dstPtr := unsafe.Pointer(&dest[0])
	//srcPtr := unsafe.Pointer(&src[0])

	//C.copy_memory_loop(dstPtr, srcPtr, C.int(len(src)))
	//C.Copy_memory_loop(unsafe.Pointer(&dest[0]), unsafe.Pointer(&src[0]), C.int(len(src)))
	//C.write_memory_rep_stosq(unsafe.Pointer(&dest[0]), unsafe.Pointer(&src[0]), C.int(len(src)))
}

func UnsafeCastBytesToHeader(source []byte) []Header {
	length := len(source) >> 3
	if length == 0 {
		return new([0]Header)[:]
	}
	capacity := cap(source) >> 3
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  capacity,
	}
	return *(*[]Header)(unsafe.Pointer(&result))
}

func UnsafeCastHeaderToBytes(source []Header) []byte {
	length := len(source) << 3
	if length == 0 {
		return new([0]byte)[:]
	}
	capacity := cap(source) << 3
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  capacity,
	}
	return *(*[]byte)(unsafe.Pointer(&result))
}

func UnsafeCastBytesToUint64(source []byte) []uint64 {
	length := len(source) >> 3
	if length == 0 {
		return new([0]uint64)[:]
	}
	capacity := cap(source) >> 3
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  capacity,
	}
	return *(*[]uint64)(unsafe.Pointer(&result))
}

func UnsafeCastUint64ToBytes(source []uint64) []byte {
	length := len(source) << 3
	if length == 0 {
		return new([0]byte)[:]
	}
	capacity := cap(source) << 3
	result := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&source[0])),
		Len:  length,
		Cap:  capacity,
	}
	return *(*[]byte)(unsafe.Pointer(&result))
}
