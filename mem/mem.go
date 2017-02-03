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

// Package mem implements memory.
//
// The basic usage is to create a block of RAM with NewLittleEndian or
// NewBigEndian then map it into the system address space by calling Map and
// passign it the system bus as argument.
//
// Handling endianness at the memory level is counter intuitive but this is a
// necessary evil: we could just pass pointers (or slices) to the CPU during
// address resolution and handle endianness at the CPU level. This works fine
// with dumb RAM, but memory-mapped IO devices need to know when their mapped
// memory is read/written to. So another solution would be for the CPU to send
// data in the proper order then call Read/Write methods but this is too much
// overhead (Read/Write cannot be inlined).
package mem

import "github.com/db47h/mirv"

type data []uint8

func (p *data) Read8(addr mirv.Address) (uint8, error) {
	if len(*p) == 0 {
		panic("Read8 on empty page")
	}
	return (*p)[addr], nil
}

func (p *data) Write8(addr mirv.Address, v uint8) error {
	if len(*p) == 0 {
		panic("Write8 on empty page")
	}
	(*p)[addr] = v
	return nil
}

// Pager is implemented by large memory structures that need to be split into
// pages in order to be mapped.
//
type Pager interface {
	mirv.Memory

	Pages(pageSize mirv.Address) []mirv.Memory
}

// big-endian memory page
type pageBE struct {
	data
}

func (p *pageBE) Pages(pageSize mirv.Address) []mirv.Memory {
	pages := make([]mirv.Memory, mirv.Address(len(p.data)))
	for i, pg := 0, p.data; pageSize <= mirv.Address(len(pg)); i++ {
		pages[i] = &pageBE{pg[:pageSize]}
		pg = pg[pageSize:]
	}
	return pages
}

func (p *pageBE) Read16(addr mirv.Address) (uint16, error) {
	return mirv.BigEndian.Uint16(p.data[addr:])
}
func (p *pageBE) Read32(addr mirv.Address) (uint32, error) {
	return mirv.BigEndian.Uint32(p.data[addr:])
}
func (p *pageBE) Read64(addr mirv.Address) (uint64, error) {
	return mirv.BigEndian.Uint64(p.data[addr:])
}
func (p *pageBE) Write16(addr mirv.Address, v uint16) error {
	return mirv.BigEndian.PutUint16(p.data[addr:], v)
}
func (p *pageBE) Write32(addr mirv.Address, v uint32) error {
	return mirv.BigEndian.PutUint32(p.data[addr:], v)
}
func (p *pageBE) Write64(addr mirv.Address, v uint64) error {
	return mirv.BigEndian.PutUint64(p.data[addr:], v)
}

// BigEndianRAM returns a block of RAM configured for BigEndian access.
// The given size should be page aligned.
func BigEndianRAM(size mirv.Address) Pager {
	return &pageBE{make(data, size)}
}
