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

package mem

import (
	"encoding/binary"
	"errors"

	"github.com/db47h/mirv"
)

// Type indicates the type of mapped memory.
//
type Type uint16

// Memory type values.
//
const (
	MemRAM Type = iota
	MemIO
)

// wrapper to track memory type
type typedMem struct {
	m mirv.Memory
	t Type
}

// NoMemory is a dummy Memory implementation that has a 0 size and returns a bus
// error for any read or write. It is used by the bus implementation for
// unmapped memory and this can aslo be used as a quick scaffolding stub to
// implement types that support only a few addressing modes (like IO devices).
//
type NoMemory struct{}

func (NoMemory) Size() mirv.Address                            { return 0 }
func (m NoMemory) Page(mirv.Address, mirv.Address) mirv.Memory { return m }
func (NoMemory) Read8(addr mirv.Address) (uint8, error)        { return 0, errBus(opRead, 1, addr) }
func (NoMemory) Write8(addr mirv.Address, v uint8) error       { return errBus(opWrite, 1, addr) }
func (NoMemory) Read16LE(addr mirv.Address) (uint16, error)    { return 0, errBus(opRead, 2, addr) }
func (NoMemory) Write16LE(addr mirv.Address, v uint16) error   { return errBus(opWrite, 2, addr) }
func (NoMemory) Read32LE(addr mirv.Address) (uint32, error)    { return 0, errBus(opRead, 4, addr) }
func (NoMemory) Write32LE(addr mirv.Address, v uint32) error   { return errBus(opWrite, 4, addr) }
func (NoMemory) Read64LE(addr mirv.Address) (uint64, error)    { return 0, errBus(opRead, 8, addr) }
func (NoMemory) Write64LE(addr mirv.Address, v uint64) error   { return errBus(opWrite, 8, addr) }
func (NoMemory) Read16BE(addr mirv.Address) (uint16, error)    { return 0, errBus(opRead, 2, addr) }
func (NoMemory) Write16BE(addr mirv.Address, v uint16) error   { return errBus(opWrite, 2, addr) }
func (NoMemory) Read32BE(addr mirv.Address) (uint32, error)    { return 0, errBus(opRead, 4, addr) }
func (NoMemory) Write32BE(addr mirv.Address, v uint32) error   { return errBus(opWrite, 4, addr) }
func (NoMemory) Read64BE(addr mirv.Address) (uint64, error)    { return 0, errBus(opRead, 8, addr) }
func (NoMemory) Write64BE(addr mirv.Address, v uint64) error   { return errBus(opWrite, 8, addr) }

type memory []uint8

var errPage = errors.New("Cross page memory access")

// New returns a new memory block of the requested size.
//
func New(size mirv.Address) mirv.Memory {
	m := make(memory, size)
	return &m
}

func (m *memory) Size() mirv.Address {
	return mirv.Address(len(*m))
}

func (m *memory) Page(addr, size mirv.Address) mirv.Memory {
	if addr == 0 && size <= mirv.Address(len(*m)) {
		return m
	}
	n := (*m)[addr : addr+size]
	return &n
}

func (m *memory) Read8(addr mirv.Address) (uint8, error) {
	if len(*m)-int(addr) < 1 {
		return 0, errPage
	}
	return (*m)[addr], nil
}
func (m *memory) Write8(addr mirv.Address, v uint8) error {
	if len(*m)-int(addr) < 1 {
		return errPage
	}
	(*m)[addr] = v
	return nil
}

func (m *memory) Read16LE(addr mirv.Address) (uint16, error) {
	if len(*m)-int(addr) < 2 {
		return 0, errPage
	}
	return binary.LittleEndian.Uint16((*m)[addr:]), nil
}
func (m *memory) Write16LE(addr mirv.Address, v uint16) error {
	if len(*m)-int(addr) < 2 {
		return errPage
	}
	binary.LittleEndian.PutUint16((*m)[addr:], v)
	return nil
}
func (m *memory) Read32LE(addr mirv.Address) (uint32, error) {
	if len(*m)-int(addr) < 4 {
		return 0, errPage
	}
	return binary.LittleEndian.Uint32((*m)[addr:]), nil
}
func (m *memory) Write32LE(addr mirv.Address, v uint32) error {
	if len(*m)-int(addr) < 4 {
		return errPage
	}
	binary.LittleEndian.PutUint32((*m)[addr:], v)
	return nil
}
func (m *memory) Read64LE(addr mirv.Address) (uint64, error) {
	if len(*m)-int(addr) < 8 {
		return 0, errPage
	}
	return binary.LittleEndian.Uint64((*m)[addr:]), nil
}
func (m *memory) Write64LE(addr mirv.Address, v uint64) error {
	if len(*m)-int(addr) < 8 {
		return errPage
	}
	binary.LittleEndian.PutUint64((*m)[addr:], v)
	return nil
}

func (m *memory) Read16BE(addr mirv.Address) (uint16, error) {
	if len(*m)-int(addr) < 2 {
		return 0, errPage
	}
	return binary.BigEndian.Uint16((*m)[addr:]), nil
}
func (m *memory) Write16BE(addr mirv.Address, v uint16) error {
	if len(*m)-int(addr) < 2 {
		return errPage
	}
	binary.BigEndian.PutUint16((*m)[addr:], v)
	return nil
}
func (m *memory) Read32BE(addr mirv.Address) (uint32, error) {
	if len(*m)-int(addr) < 4 {
		return 0, errPage
	}
	return binary.BigEndian.Uint32((*m)[addr:]), nil
}
func (m *memory) Write32BE(addr mirv.Address, v uint32) error {
	if len(*m)-int(addr) < 4 {
		return errPage
	}
	binary.BigEndian.PutUint32((*m)[addr:], v)
	return nil
}
func (m *memory) Read64BE(addr mirv.Address) (uint64, error) {
	if len(*m)-int(addr) < 8 {
		return 0, errPage
	}
	return binary.BigEndian.Uint64((*m)[addr:]), nil
}
func (m *memory) Write64BE(addr mirv.Address, v uint64) error {
	if len(*m)-int(addr) < 8 {
		return errPage
	}
	binary.BigEndian.PutUint64((*m)[addr:], v)
	return nil
}
