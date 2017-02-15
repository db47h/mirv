// Copyright © 2017 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the “Software”), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package mem provides basic memory components.
//
package mem

import (
	"errors"

	"github.com/db47h/mirv"
)

// Type indicates the type of mapped memory.
//
type Type uint16

// Memory type values.
//
const (
	MemNone Type = iota // non functional memory
	MemRAM              // RAM
	MemIO               // Memory Mapped IO
)

// Interface wraps the methods exported by types that can be used as memory.
//
// Address arguments are always specified relative to the beginning of the
// memory block. For example:
//
//	// a small system with ROM starting at 0x0000, RAM at 0x8000
//	var b mem.Bus
//	rom := mem.New(32768, mirv.LittleEndian)
//	ram := mem.New(32768, mirv.LittleEndian)
//  b.Map(0, rom)
//	b.Map(32768, ram)
//	ram.Write8(4096, 42)        // this should write at physical address 32768+4096
//	p := b.Memory(32768 + 4096) // returns the ram block referenced above
//	if p.Read8(0) != 42 {		// for which index 0 is physical address 32768+4096
//		panic("Wrong memory block")
//	}
//
type Interface interface {
	mirv.ByteOrdered // memory is byte ordered

	Size() mirv.Address // Size in bytes of the memory block
	Type() Type         // Memory type: MemIO or MemRAM

	// Read/Write methods
	Read8(mirv.Address) (uint8, error)
	Write8(mirv.Address, uint8) error
	Read16(mirv.Address) (uint16, error)
	Write16(mirv.Address, uint16) error
	Read32(mirv.Address) (uint32, error)
	Write32(mirv.Address, uint32) error
	Read64(mirv.Address) (uint64, error)
	Write64(mirv.Address, uint64) error
}

// NoMemory is a dummy Memory implementation that has a 0 size and returns a bus
// error for any read or write. It is used by the bus implementation for
// unmapped memory and can aslo be used as a quick scaffolding stub to implement
// types that support only a few addressing modes (like IO devices).
//
type NoMemory struct{}

// Size always returns 0.
//
func (NoMemory) Size() mirv.Address { return 0 }

// Type returns MemNone
//
func (NoMemory) Type() Type { return MemNone }

// ByteOrder returns 0 (unknown byte order)
//
func (NoMemory) ByteOrder() mirv.ByteOrder { return 0 }

var errPage = errors.New("Cross page memory access")

// NewRAM returns a new RAM block of the requested size and byte order.
//
func NewRAM(size mirv.Address, byteOrder mirv.ByteOrder) Interface {
	m := make([]uint8, size)
	if byteOrder == mirv.LittleEndian {
		return (*littleEndian)(&m)
	}
	return (*bigEndian)(&m)
}

//go:generate go run mem_gen.go -o mem_rw.go
