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
// The Memory interface is an unusually big interface. A previous version had a
// single Get() method that returned a uint8 slice and all read/writes went
// through mem.Bus. This was not practical for memory mapped IO: we need a way
// for IO devices to catch writes.
//
// Possible solutions:
//
//	- Current version with 16 methods. Ugly, but acceptable performance.
//	- a single Get(addr) []uint8. In order for MMIO devices to catch writes.
//	  Need to implement a clock mechanism (after every CPU cycle, call a
//	  bus.Tick() method that propagate the clock tick to all devices).
//	  Not yet tested.
//	- Define the Uint16LittleEndian, Uint32BigEndian types, etc. with
//    read/write methods:
//
//	// types.go
//	type Data interface {
//		Write([]uint8) error
//	}
//
//	type Uint16LE uint16
//
//	func (v Uint16LE) Write(dst []uint8) {
//		binary.LittleEndian.PutUint16(dst, v)
//	}
//
//	// constructor for Uint16LE -- also needs a method like memory.Read(Address) []uint8
//	func Read16LE(src []uint8) Uint16LE {
//		return Uint16LE(binary.Uint16(src))
//	}
//
//	// mem.go
//	func (m memory) Write(dst Address, v Data) {
//		v.Write(m[address:])
//	}
//
// The Memory interface would be down to 4 methods, but tests showed very bad
// performance.
//
// Any suggestions welcome.
//
package mirv

// Address is the guest address type.
type Address uint

// Memory is implemented by types exposing a memory-like interface.
//
// Address arguments are always specified relative to the beginning of the
// memory segment. For example:
//
//	// a small system with ROM starting at 0x0000, RAM at 0x8000
//	b := bus.New(4096, 32768)
//	rom := mem.New(32768)
//	ram := mem.New(32768)
//  b.Map(0, rom)
//	b.Map(32768, ram)
//	ram.Write8(4096, 42)        // this should write at physical address 32768+4096
//	p := b.Memory(32768 + 4096) // returns the ram page referenced above
//	if p.Read8(0) != 42 {		// for which index 0 is physical address 32768+4096
//		panic("Bad page")
//	}
//
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
