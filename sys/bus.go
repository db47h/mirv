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
// IMPLIED, INCLUDING BUT NOT LIMITEWORK=/tmp/go-build756775462D TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sys

import (
	"fmt"

	"github.com/db47h/mirv"
)

type tag mirv.Address

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
	addr mirv.Address
}

func errBus(op busOp, size uint8, addr mirv.Address) *ErrBus {
	return &ErrBus{op: op, sz: size, addr: addr}
}

func (e *ErrBus) Error() string {
	return fmt.Sprintf("bus error: %v/%d @ address %x", e.op, e.sz, e.addr)
}

// nilMemory dummy mirv.Memory implementation that returns bus error
type nilMemory struct{}

func (nilMemory) Read8(addr mirv.Address) (uint8, error)    { return 0, errBus(opRead, 8, addr) }
func (nilMemory) Read16(addr mirv.Address) (uint16, error)  { return 0, errBus(opRead, 16, addr) }
func (nilMemory) Read32(addr mirv.Address) (uint32, error)  { return 0, errBus(opRead, 32, addr) }
func (nilMemory) Read64(addr mirv.Address) (uint64, error)  { return 0, errBus(opRead, 64, addr) }
func (nilMemory) Write8(addr mirv.Address, v uint8) error   { return errBus(opWrite, 8, addr) }
func (nilMemory) Write16(addr mirv.Address, v uint16) error { return errBus(opWrite, 16, addr) }
func (nilMemory) Write32(addr mirv.Address, v uint32) error { return errBus(opWrite, 32, addr) }
func (nilMemory) Write64(addr mirv.Address, v uint64) error { return errBus(opWrite, 64, addr) }

type cacheEntry struct {
	tag
	m mirv.Memory
}

// Bus is a simplistic system bus.
//
type Bus struct {
	sz   mirv.Address // page size
	bits uint8        // bit count of page part in address
	pom  mirv.Address // page offset mask
	pnm  tag          // page number mask

	pages map[tag]mirv.Memory
	cache []*cacheEntry
}

// NewBus creates a new bus configured with the given parameters.
//
// pageSize must be an exponent of 2 and should match at least the simulated
// CPU's minimum natural page size. A typical page size is 4096 bytes.
//
// Internally, mapped memory pages are kept in a hash map. In order to speed up
// address-to-Memory-interface lookups, recent lookups are kept in a cache (that
// works like a MMU). cacheSize is an exponent of 2 value that determines the
// number of entries in this cache.
//
// The total size of memory that is addressable without costly map lookups is
// cacheSize * pageSize bytes. For best performance, it is advisable to keep
// this value equal to the amount of real memory available to the simulated CPU.
// Each entry in the cache is 24 bytes long on a 64 bits host. Simulations that
// require large amounts of memory might want to use larger page sizes in order
// to compensate for the overhead of a large cache.
//
func NewBus(pageSize mirv.Address, cacheSize uint) *Bus {
	if pageSize == 0 || pageSize&(pageSize-1) != 0 {
		panic("Page size must be an exponent of 2.")
	}
	if cacheSize == 0 || cacheSize&(cacheSize-1) != 0 {
		panic("Cache size must be an exponent of 2.")
	}
	// count bits
	const _m = ^mirv.Address(0)
	var b uint8 = 8 << (_m>>8&1 + _m>>16&1 + _m>>32&1) // (1 << log₂(_m)) * 8 = addr bits
	s, b := b<<1, b-1
	for x := pageSize; s > 0; s >>= 1 {
		if y := x << s; y != 0 {
			b -= s
			x = y
		}
	}
	bus := &Bus{
		sz:    pageSize,
		bits:  b,
		pom:   pageSize - 1,
		pnm:   tag(cacheSize) - 1,
		pages: make(map[tag]mirv.Memory),
		cache: make([]*cacheEntry, cacheSize),
	}
	return bus
}

func (b *Bus) tag(addr mirv.Address) tag {
	return tag(addr >> b.bits)
}

// PageSize returns the configured page size.
//
func (b *Bus) PageSize() mirv.Address {
	return b.sz
}

// Map maps a series of pages starting at addr to the given Memory interfaces.
// Map panics if a page is already mapped or if the address is not page-aligned.
//
// Use the Memory method to check if a given address is mapped.
//
func (b *Bus) Map(addr mirv.Address, mem ...mirv.Memory) {
	for _, m := range mem {
		if addr&b.pom != 0 {
			panic("Address must be page-aligned")
		}
		tag := b.tag(addr)
		if b.pages[tag] != nil {
			panic("Address already mapped")
		}
		b.pages[tag] = m
		n := addr + b.sz
		if n < addr {
			panic("Page mapping past end of addressable memory.")
		}
		addr = n
	}
}

// Unmap unmaps n pages starting at the given address.
//
func (b *Bus) Unmap(addr mirv.Address, n int) {
	for i, a := n, addr; a >= addr && i > 0; i, a = i-1, a+b.sz {
		t := b.tag(a)
		i := t & b.pnm
		if e := b.cache[i]; e.tag == t {
			b.cache[i] = &cacheEntry{^tag(0), nilMemory{}}
		}
		delete(b.pages, t)
	}
}

// Memory returns the Memory interface mapped to address addr, and nil if unmapped.
//
func (b *Bus) Memory(addr mirv.Address) mirv.Memory {
	switch m := b.memory(b.tag(addr)).(type) {
	case nilMemory:
		return nil
	default:
		return m
	}
}

func (b *Bus) memory(tag tag) mirv.Memory {
	i := tag & b.pnm
	if e := b.cache[i]; e != nil && e.tag == tag {
		return e.m
	}
	if m := b.pages[tag]; m != nil {
		b.cache[i] = &cacheEntry{tag, m}
		return m
	}
	return nilMemory{}
}

// Read8 reads unsigned 8 bit value from address.
//
func (b *Bus) Read8(addr mirv.Address) (uint8, error) {
	return b.memory(b.tag(addr)).Read8(addr & b.pom)
}

// Read16 reads unsigned 16 bit value from address.
//
func (b *Bus) Read16(addr mirv.Address) (uint16, error) {
	return b.memory(b.tag(addr)).Read16(addr & b.pom)
}

// Read32 reads unsigned 32 bit value from address.
//
func (b *Bus) Read32(addr mirv.Address) (uint32, error) {
	return b.memory(b.tag(addr)).Read32(addr & b.pom)
}

// Read64 reads unsigned 64 bit value from address.
//
func (b *Bus) Read64(addr mirv.Address) (uint64, error) {
	return b.memory(b.tag(addr)).Read64(addr & b.pom)
}

// Write8 writes unsigned 8 bit value to address.
//
func (b *Bus) Write8(addr mirv.Address, v uint8) error {
	return b.memory(b.tag(addr)).Write8(addr&b.pom, v)
}

// Write16 writes unsigned 16 bit value to address.
//
func (b *Bus) Write16(addr mirv.Address, v uint16) error {
	return b.memory(b.tag(addr)).Write16(addr&b.pom, v)
}

// Write32 writes unsigned 32 bit value to address.
//
func (b *Bus) Write32(addr mirv.Address, v uint32) error {
	return b.memory(b.tag(addr)).Write32(addr&b.pom, v)
}

// Write64 writes unsigned 64 bit value to address.
//
func (b *Bus) Write64(addr mirv.Address, v uint64) error {
	return b.memory(b.tag(addr)).Write64(addr&b.pom, v)
}
