// Copyright 2017 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package mirv is the root package for MIRV. It also provides common types
// used by all sub-packages.
//
package mirv

import "fmt"

// Address is the guest address type.
type Address uint

//go:generate stringer -type busOp "$GOFILE"
type busOp int

const (
	opRead busOp = iota
	opWrite
)

// ErrBus wraps a bus error.
//
type ErrBus struct {
	op   busOp
	sz   uint8
	addr Address
}

func errBus(op busOp, size uint8, addr Address) *ErrBus {
	return &ErrBus{op: op, sz: size, addr: addr}
}

func (e *ErrBus) Error() string {
	return fmt.Sprintf("bus error: %v/%d @ address %x", e.op, e.sz, e.addr)
}

// Memory is implemented by types exposing a memory-like interface.
//
// In the Get and Page methods, addresses are always specified relative to the
// beginning of the memory segment. For example:
//
//	// a small system with ROM starting at 0x0000, RAM at 0x8000
//	b := bus.New(4096, 32768)
//	rom := mem.New(32768)
//	ram := mem.New(32768)
//  b.Map(0, rom)
//	b.Map(32768, ram)
//	ram.Get(4096)[0] = 42		// this should write at physical address 32768+4096
//	p := b.Memory(32768 + 4096) // returns the sub-page from ram referenced above
//	if p.Get(0)[0] != 42 {		// for which index 0 is physical address 32768+4096
//		panic("Bad page")
//	}
type Memory interface {
	// Size in bytes of the memory block.
	Size() Address
	Page(Address, Address) Memory

	Read8(Address) (uint8, error)
	Write8(Address, uint8) error

	Read16LE(Address) (uint16, error)
	Write16LE(Address, uint16) error
	Read32LE(Address) (uint32, error)
	Write32LE(Address, uint32) error
	Read64LE(Address) (uint64, error)
	Write64LE(Address, uint64) error

	Read16BE(Address) (uint16, error)
	Write16BE(Address, uint16) error
	Read32BE(Address) (uint32, error)
	Write32BE(Address, uint32) error
	Read64BE(Address) (uint64, error)
	Write64BE(Address, uint64) error
}

// VoidMemory is a dummy Memory implementation that returns a bus error. It is
// used by the bus implementation for unmapped memory and this can aslo be used
// as a quick scaffolding stub to implement memory types that support only a few
// addressing modes.
//
type VoidMemory struct{}

func (VoidMemory) Size() Address                          { return 0 }
func (m VoidMemory) Page(Address, Address) Memory         { return m }
func (VoidMemory) Read8(addr Address) (uint8, error)      { return 0, errBus(opRead, 1, addr) }
func (VoidMemory) Write8(addr Address, v uint8) error     { return errBus(opWrite, 1, addr) }
func (VoidMemory) Read16LE(addr Address) (uint16, error)  { return 0, errBus(opRead, 2, addr) }
func (VoidMemory) Write16LE(addr Address, v uint16) error { return errBus(opWrite, 2, addr) }
func (VoidMemory) Read32LE(addr Address) (uint32, error)  { return 0, errBus(opRead, 4, addr) }
func (VoidMemory) Write32LE(addr Address, v uint32) error { return errBus(opWrite, 4, addr) }
func (VoidMemory) Read64LE(addr Address) (uint64, error)  { return 0, errBus(opRead, 8, addr) }
func (VoidMemory) Write64LE(addr Address, v uint64) error { return errBus(opWrite, 8, addr) }
func (VoidMemory) Read16BE(addr Address) (uint16, error)  { return 0, errBus(opRead, 2, addr) }
func (VoidMemory) Write16BE(addr Address, v uint16) error { return errBus(opWrite, 2, addr) }
func (VoidMemory) Read32BE(addr Address) (uint32, error)  { return 0, errBus(opRead, 4, addr) }
func (VoidMemory) Write32BE(addr Address, v uint32) error { return errBus(opWrite, 4, addr) }
func (VoidMemory) Read64BE(addr Address) (uint64, error)  { return 0, errBus(opRead, 8, addr) }
func (VoidMemory) Write64BE(addr Address, v uint64) error { return errBus(opWrite, 8, addr) }
