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
	"fmt"
	"io"

	"github.com/db47h/mirv"
)

//go:generate stringer -type busOp bus.go mem.go
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

var nilMemory = VoidMemory{}

type tag mirv.Address

type cacheEntry struct {
	tag tag
	m   mirv.Memory
}

// Bus is a simplistic memory bus. The current implementation only provides
// guest <-> host memory mapping and helper functions for reading and writing
// data with different byte orders.
//
// Mapped memory is split into pages. Read and writes do not need to be aligned
// but cannot cross page boundaries. i.e. with a page size of 4096 bytes, trying
// to read a uint64 at address 4095 will result in a bus error. This should be
// of no consequence where the simulated CPU does not support unaligned
// read/writes, but extra steps must be taken with others.
//
type Bus struct {
	sz   mirv.Address // page size
	bits uint8        // bit count of page part in address
	pom  mirv.Address // page offset mask
	pnm  tag          // page number mask

	pages map[tag]*typedMem
	cache []cacheEntry
}

// NewBus creates a new bus configured with the given parameters.
//
// pageSize is the page size in bytes. It must be an exponent of 2 and should
// match at least the simulated CPU's minimum natural page size. A typical page
// size is 4096 bytes.
//
// Internally, mapped memory pages are kept in a hash map. In order to speed up
// address-to-Memory-interface lookups, recent lookups are kept in a cache (that
// works like a MMU). cacheSize is an exponent of 2 value that determines the
// number of entries in this cache.
//
// The total size of memory that is addressable without costly map lookups is
// cacheSize * pageSize bytes. For best performance, it is advisable to keep
// this value equal to the amount of real memory available to the simulated CPU.
// Each entry in the cache is 24 bytes long (plus 16 bytes for memory slices) on
// a 64 bits host. That's an overhead of about 1% for a typical page size of
// 4096 bytes. Simulations that require large amounts of memory might want to
// use larger page sizes in order to compensate for the overhead of a large
// cache.
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
		pages: make(map[tag]*typedMem),
		cache: make([]cacheEntry, cacheSize),
	}

	// prefill cache
	ne := cacheEntry{^tag(0), nilMemory}
	for i := range bus.cache {
		bus.cache[i] = ne
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

// Map maps a the memory pages starting at addr to the given Memory interfaces.
// Map panics if a page is already mapped or if the address is not page-aligned.
//
// Use the Memory method to check if a given address is mapped.
//
// The memType parameter does not affect the page mapping in any way. It only
// serves as a differentiator for the MappedRange method.
//
func (b *Bus) Map(addr mirv.Address, m mirv.Memory, memType Type) {
	if addr&b.pom != 0 {
		panic("Address must be page-aligned")
	}

	pages := m.Size() >> b.bits
	for i, pa := mirv.Address(0), addr; i < pages; i++ {
		tag := b.tag(pa)
		if b.pages[tag] != nil {
			panic("Address already mapped")
		}
		b.pages[tag] = &typedMem{m: m.Page(pa-addr, b.sz), t: memType}
		n := pa + b.sz
		if n <= addr {
			if i == pages-1 {
				// we've just mapped the last page, and nothing more to mapm all is good.
				break
			}
			panic("Page mapping past end of addressable memory.")
		}
		pa = n
	}
}

// Unmap unmaps n pages starting at the given address.
//
func (b *Bus) Unmap(addr mirv.Address, n int) {
	for i, a := n, addr; a >= addr && i > 0; i, a = i-1, a+b.sz {
		t := b.tag(a)
		i := t & b.pnm
		if e := b.cache[i]; e.tag == t {
			b.cache[i].tag, b.cache[i].m = ^tag(0), nilMemory
		}
		delete(b.pages, t)
	}
}

// MappedRange reports the largest addressable range [low, high) for the given
// memory type. i.e. only the memory addresses low and high-1 are guaranteed to
// be mapped, but there may be unmapped pages in between.
//
// As a result of 2-complement arithmetic, the high address may be 0 if the
// highest memory address is mapped.
//
// The purpose of this function is to ease setup of some CPUs that default some registers
// to start or end of memory.
//
func (b *Bus) MappedRange(t Type) (low, high mirv.Address) {
	low = ^mirv.Address(0)
	for tag, m := range b.pages {
		if m.t != t {
			continue
		}
		tag := mirv.Address(tag << b.bits)
		if tag < low {
			low = tag
		}
		if end := tag + b.sz - 1; end > high {
			high = end
		}
	}
	return low, high + 1
}

// Memory returns the Memory interface mapped to address addr. If the address is
// not mapped, it returns a 0 sized Memory interface:
//
//	m := bus.Memory(addr)
//	if m.Size() == 0 {
//		// addr is not mapped
//		// ...
//	}
//
//
func (b *Bus) Memory(addr mirv.Address) mirv.Memory {
	tag := b.tag(addr)
	i := tag & b.pnm
	if e := b.cache[i]; e.tag == tag {
		return e.m
	}
	if m := b.pages[tag]; m != nil {
		m := m.m
		b.cache[i].tag, b.cache[i].m = tag, m
		return m
	}
	return nilMemory
}

// Read8 reads the uint8 value from the memory mapped at address addr.
//
func (b *Bus) Read8(addr mirv.Address) (uint8, error) {
	return b.Memory(addr).Read8(addr & b.pom)
}

// Write8 writes the uint8 value v to the memory mapped at address addr.
//
func (b *Bus) Write8(addr mirv.Address, v uint8) error {
	return b.Memory(addr).Write8(addr&b.pom, v)
}

// Read16LE reads the little endian uint16 value from the memory mapped at address addr.
//
func (b *Bus) Read16LE(addr mirv.Address) (uint16, error) {
	return b.Memory(addr).Read16LE(addr & b.pom)
}

// Read32LE reads the little endian uint16 value from the memory mapped at address addr.
//
func (b *Bus) Read32LE(addr mirv.Address) (uint32, error) {
	return b.Memory(addr).Read32LE(addr & b.pom)
}

// Read64LE reads the little endian uint16 valuefrom the memory mapped  at address addr.
//
func (b *Bus) Read64LE(addr mirv.Address) (uint64, error) {
	return b.Memory(addr).Read64LE(addr & b.pom)
}

// Write16LE writes the little endian uint16 value v to the memory mapped at address addr.
//
func (b *Bus) Write16LE(addr mirv.Address, v uint16) error {
	return b.Memory(addr).Write16LE(addr&b.pom, v)
}

// Write32LE writes the little endian uint32 value v to the memory mapped at address addr.
//
func (b *Bus) Write32LE(addr mirv.Address, v uint32) error {
	return b.Memory(addr).Write32LE(addr&b.pom, v)
}

// Write64LE writes the little endian uint64 value v to the memory mapped at address addr.
//
func (b *Bus) Write64LE(addr mirv.Address, v uint64) error {
	return b.Memory(addr).Write64LE(addr&b.pom, v)
}

// Read16BE reads the big endian uint16 value from the memory mapped at address addr.
//
func (b *Bus) Read16BE(addr mirv.Address) (uint16, error) {
	return b.Memory(addr).Read16BE(addr & b.pom)
}

// Read32BE reads the big endian uint16 value from the memory mapped at address addr.
//
func (b *Bus) Read32BE(addr mirv.Address) (uint32, error) {
	return b.Memory(addr).Read32BE(addr & b.pom)
}

// Read64BE reads the big endian uint16 valuefrom the memory mapped  at address addr.
//
func (b *Bus) Read64BE(addr mirv.Address) (uint64, error) {
	return b.Memory(addr).Read64BE(addr & b.pom)
}

// Write16BE writes the big endian uint16 value v to the memory mapped at address addr.
//
func (b *Bus) Write16BE(addr mirv.Address, v uint16) error {
	return b.Memory(addr).Write16BE(addr&b.pom, v)
}

// Write32BE writes the big endian uint32 value v to the memory mapped at address addr.
//
func (b *Bus) Write32BE(addr mirv.Address, v uint32) error {
	return b.Memory(addr).Write32BE(addr&b.pom, v)
}

// Write64BE writes the big endian uint64 value v to the memory mapped at address addr.
//
func (b *Bus) Write64BE(addr mirv.Address, v uint64) error {
	return b.Memory(addr).Write64BE(addr&b.pom, v)
}

type busWriter struct {
	addr mirv.Address
	b    *Bus
}

func (w *busWriter) Write(p []byte) (n int, err error) {
	var (
		b                = w.b
		page mirv.Memory = b.Memory(w.addr)
		tag              = b.tag(w.addr)
	)
	for _, c := range p {
		if t := b.tag(w.addr); t != tag {
			page = w.b.Memory(w.addr)
			if page.Size() == 0 {
				return n, io.EOF
			}
			tag = t
		}
		if err := page.Write8(w.addr&w.b.pom, c); err != nil {
			return n, io.EOF
		}
		w.addr++
		n++
	}
	return n, nil
}

// Writer returns an io.Writer to the mapped memory starting at addr.
//
func (b *Bus) Writer(addr mirv.Address) io.Writer {
	return &busWriter{addr, b}
}