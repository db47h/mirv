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

package mem

import (
	"errors"
	"fmt"
	"io"

	"github.com/db47h/mirv"
)

var (
	errOverlap     = errors.New("memory block overlap")
	errOverflow    = errors.New("memory block overflows address space")
	errNoMemoryMap = errors.New("no memory mapped")
)

//go:generate stringer -type busOp .
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

var nilMemory = &block{
	s: ^mirv.Address(0),
	e: 0,
	m: NoMemory{},
}

type block struct {
	s, e mirv.Address
	m    Interface
}

func (b *block) overlaps(blk *block) bool {
	return blk.s >= b.s && blk.s <= b.e || b.s >= blk.s && b.s <= blk.e
}

func (b *block) contains(addr mirv.Address) bool {
	return b != nil && addr <= b.e && addr >= b.s
}

// Bus is a simplistic memory bus. The current implementation only provides
// guest <-> host memory mapping and helper functions for reading and writing
// data with different byte orders.
//
// Read and writes do not need to be aligned but cannot cross block boundaries.
// For example:
//
//	var b Bus
//	// map 2 x 32KiB RAM blocks at addresses 0 and 0x8000 respectively.
//	b.Map(0x0000, mem.NewRAM(0x8000))
//	b.Map(0x8000, mem.NewRAM(0x8000))
//	u16, err := b.Read16LE(0x0001) // OK: unaligned but all bytes are in first block
//	u16, err = b.Read16LE(0x7FFF) // Will return an error: both bytes are in different blocks
//
// Bus keeps mapped blocks of memory in a slice. When doing guest->host address
// resolution, it does a binary search for the corresponding interface in that
// slice. In order to improve up performance, if also keeps a reference to a
// "preferred" memory block that will always be checked first. This preferred
// memory block is by default the first mapped block, and can also be set by the
// user by calling the Preferred method.
//
type Bus struct {
	b []*block
	p *block // preferred mem block
}

//go:generate go run bus_gen.go -o bus_rw.go

// Map maps a memory block starting at addr to the given Interface. Map
// returns a non nil error if the block is block already mapped, overlaps with
// another mapped block or if addr+m.Size() is greater than the maximum value of
// mirv.Address.
//
func (b *Bus) Map(addr mirv.Address, m Interface) error {
	if m.Size() == 0 {
		return nil
	}
	end := addr + (m.Size() - 1)
	if end < addr {
		return errOverflow
	}
	return b.insert(&block{
		s: addr,
		e: end,
		m: m,
	})
}

func (b *Bus) insertIdx(blk *block) int {
	var l, h = 0, len(b.b)
	for l != h {
		i := l + (h-l)/2 // == (l+h)/2 without overflow
		bi := b.b[i]
		if blk.overlaps(bi) {
			return -1
		}
		if bi.s > blk.s {
			h = i
		} else {
			l = i + 1
		}
	}
	return l
}

func (b *Bus) insert(blk *block) error {
	if len(b.b) == 0 && b.p == nil {
		b.p = blk
		return nil
	}
	i := b.insertIdx(blk)
	if b.p != nil && blk.overlaps(b.p) || i < 0 {
		return errOverlap
	}
	b.b = append(b.b, nil)
	copy(b.b[i+1:], b.b[i:])
	b.b[i] = blk
	return nil
}

// Preferred sets the preferred memory block. When resolving guest to host
// addresses, the memory block containing addr will be checked first.
//
func (b *Bus) Preferred(addr mirv.Address) {
	if b.p.contains(addr) {
		return
	}
	o := b.findIdx(addr)
	// Unmapped, don't touch anything.
	if o < 0 {
		// TODO: panic, error?
		return
	}
	if b.p == nil {
		b.p = b.b[o]
		copy(b.b[o:], b.b[o+1:])
		b.b = b.b[:len(b.b)-1]
		return
	}
	i := b.insertIdx(b.p)
	var t *block
	t, b.p = b.p, b.b[o]
	if i < o {
		// shift right
		copy(b.b[i+1:], b.b[i:o])
	} else if i > o {
		// shift left
		copy(b.b[o:], b.b[o+1:i])
		i--
	}
	b.b[i] = t
}

// Remap maps or remaps the memory block containing the given address. If the given
// address is already mapped, attempt .
//
// Remap panics if the size of the new memory Interface is too large to fit.
//
// This function is meant to help implement the brk/sbrk syscalls and dynamic
// memory bank swapping.
//
func (b *Bus) Remap(addr mirv.Address, m Interface) error {
	if b.p.contains(addr) {
		return nil
	}
	i := b.findIdx(addr)
	if i < 0 {
		return b.Map(addr, m)
	}

	blk := b.b[i]
	end := addr + (m.Size() - 1)
	if s := m.Size(); s > blk.m.Size() && i < len(b.b)-1 {
		if end >= b.b[i+1].s {
			return errOverlap
		}
	}
	blk.m = m
	blk.e = end
	return nil
}

// MappedRange reports the largest addressable range [low, high) for the given
// memory type. i.e. only the memory addresses low and high-1 are guaranteed to
// be mapped, but there may be unmapped addresses in between.
//
// As a result of 2-complement arithmetic, the high address may be 0 if the
// highest memory address is mapped.
//
// The only case where this function returns a non nil error is when there is no
// mapped memory of the requested type.
//
// The purpose of this function is to ease setup of some CPUs that default some
// registers to start or end of memory.
//
func (b *Bus) MappedRange(t Type) (low, high mirv.Address, err error) {
	var ok bool // true if low/high have changed
	low = ^mirv.Address(0)
	if blk := b.p; blk != nil && blk.m.Type() == t {
		low = blk.s
		high = blk.e
		ok = true
	}

	for _, blk := range b.b {
		if blk.m.Type() != t {
			continue
		}
		ok = true
		if blk.s < low {
			low = blk.s
		}
		if blk.e > high {
			high = blk.e
		}
	}

	if ok {
		return low, high + 1, nil
	}
	return 0, 0, errNoMemoryMap
}

// Memory returns the base address and memory Interface mapped to address addr.
// If the address is not mapped, it returns a 0 sized Memory interface:
//
//	base, m := bus.Memory(addr)
//	if m.Size() == 0 {
//		// addr is not mapped
//		// ...
//	}
//	log.Printf("0x%X is in a %d bytes block mapped at 0x%X", addr, m.Size(), base)
//
func (b *Bus) Memory(addr mirv.Address) (mirv.Address, Interface) {
	if len(b.b) == 0 {
		return 0, nilMemory.m
	}
	e := b.memory(addr)
	return e.s, e.m
}

// findIdx returns the index if the block containing addr. If not found, returns
// -1. It does not check b.p.
//
func (b *Bus) findIdx(addr mirv.Address) int {
	for bb, l := b.b, len(b.b); l > 0; l = len(bb) {
		i := l / 2
		blk := bb[i]
		if blk.s > addr {
			bb = bb[:i]
			continue
		}
		if blk.e < addr {
			bb = bb[i+1:]
			continue
		}
		return i
	}
	return -1
}

// find returns the *block containing addr or nilMemory if not found. Does not
// check b.p.
//
func (b *Bus) find(addr mirv.Address) *block {
	for bb, l := b.b, len(b.b); l > 0; l = len(bb) {
		i := l / 2
		blk := bb[i]
		if blk.s > addr {
			bb = bb[:i]
			continue
		}
		if blk.e < addr {
			bb = bb[i+1:]
			continue
		}
		return blk
	}
	return nilMemory
}

// memory finds the *block containing addr. Does check b.p.
//
func (b *Bus) memory(addr mirv.Address) *block {
	if b.p.contains(addr) {
		return b.p
	}
	return b.find(addr)
}

type busWriter struct {
	addr mirv.Address
	b    *Bus
}

func (w *busWriter) Write(p []byte) (n int, err error) {
	var (
		b   = w.b
		blk = b.memory(w.addr)
		m   = blk.m
	)
	for _, c := range p {
		if w.addr > blk.e {
			blk = b.memory(w.addr)
			m = blk.m
		}
		if err := m.Write8(w.addr-blk.s, c); err != nil {
			return n, io.EOF
		}
		w.addr++
		n++
	}
	return n, nil
}

// Writer returns an io.Writer to the mapped memory starting at addr.
//
// Unlike the Bus Read/Write methods, the returned writer can write
// across several memory blocks as long as they are contiguous.
//
func (b *Bus) Writer(addr mirv.Address) io.Writer {
	return &busWriter{addr, b}
}
