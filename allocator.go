package wordexp

/*
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

// allocator is used to allocate and free native memory and also
// record how many allocations are outstanding.
type allocator interface {
	Alloc(size C.size_t) (unsafe.Pointer, error)
	Free(ptr unsafe.Pointer)
}

var _ allocator = (*oomAllocator)(nil)
var _ allocator = (*nativeAllocator)(nil)

type nativeAllocator struct{}

type oomAllocator struct {
	nativeAllocator
}

func (a *nativeAllocator) Alloc(size C.size_t) (unsafe.Pointer, error) {
	// What does C.malloc() return?
	//
	// From the Go 1.8+ docs:
	//
	// As a special case, C.malloc does not call the C library
	// malloc directly but instead calls a Go helper function that
	// wraps the C library malloc but guarantees never to return
	// nil. If C's malloc indicates out of memory, the helper
	// function crashes the program, like when Go itself runs out
	// of memory.
	//
	// Note: this is not a change in Go 1.8, merely a
	// clarification of the existing behaviour.
	//
	// See: https://golang.org/issue/16309.
	p := C.malloc(size)
	return p, nil
}

func (a *nativeAllocator) Free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func newNativeAllocator() *nativeAllocator {
	return &nativeAllocator{}
}

func (a *oomAllocator) Alloc(size C.size_t) (unsafe.Pointer, error) {
	// We return a as a pointer to make go vet happy. But callers
	// should really be checking for error first and we
	// categorically return an error here.
	return unsafe.Pointer(a), errors.New("out of memory")
}

func newOutOfMemoryAllocator() *oomAllocator {
	return &oomAllocator{}
}
