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

// Address is the guest address type.
type Address uint

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

	// returns the raw memory starting at address addr.
	Get(addr Address) []uint8

	// returns a new memory interface to a sub-page of the given size starting at addr.
	// If Size() is smaller or equal than size, this method may return its receiver.
	Page(addr, size Address) Memory
}
