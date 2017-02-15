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
	"encoding/binary"

	"github.com/db47h/mirv"
)

// big endian  memory interface
type bigEndian []uint8

func (m *bigEndian) Size() mirv.Address { return mirv.Address(len(*m)) }

func (m *bigEndian) Type() Type { return MemRAM }

func (m *bigEndian) ByteOrder() mirv.ByteOrder { return mirv.BigEndian }

// Read8 returns the 8 bits big endian  value at address addr.
//
func (m *bigEndian) Read8(addr mirv.Address) (uint8, error) {
	if len((*m)[addr:]) < 1 {
		return 0, errPage
	}
	return (*m)[addr], nil
}

// Write8 writes the 8 bits big endian  value to address addr.
//
func (m *bigEndian) Write8(addr mirv.Address, v uint8) error {
	if len((*m)[addr:]) < 1 {
		return errPage
	}
	(*m)[addr] = v
	return nil
}

// Read16 returns the 16 bits big endian  value at address addr.
//
func (m *bigEndian) Read16(addr mirv.Address) (uint16, error) {
	if len((*m)[addr:]) < 2 {
		return 0, errPage
	}
	return binary.BigEndian.Uint16((*m)[addr:]), nil
}

// Write16 writes the 16 bits big endian  value to address addr.
//
func (m *bigEndian) Write16(addr mirv.Address, v uint16) error {
	if len((*m)[addr:]) < 2 {
		return errPage
	}
	binary.BigEndian.PutUint16((*m)[addr:], v)
	return nil
}

// Read32 returns the 32 bits big endian  value at address addr.
//
func (m *bigEndian) Read32(addr mirv.Address) (uint32, error) {
	if len((*m)[addr:]) < 4 {
		return 0, errPage
	}
	return binary.BigEndian.Uint32((*m)[addr:]), nil
}

// Write32 writes the 32 bits big endian  value to address addr.
//
func (m *bigEndian) Write32(addr mirv.Address, v uint32) error {
	if len((*m)[addr:]) < 4 {
		return errPage
	}
	binary.BigEndian.PutUint32((*m)[addr:], v)
	return nil
}

// Read64 returns the 64 bits big endian  value at address addr.
//
func (m *bigEndian) Read64(addr mirv.Address) (uint64, error) {
	if len((*m)[addr:]) < 8 {
		return 0, errPage
	}
	return binary.BigEndian.Uint64((*m)[addr:]), nil
}

// Write64 writes the 64 bits big endian  value to address addr.
//
func (m *bigEndian) Write64(addr mirv.Address, v uint64) error {
	if len((*m)[addr:]) < 8 {
		return errPage
	}
	binary.BigEndian.PutUint64((*m)[addr:], v)
	return nil
}

// little endian  memory interface
type littleEndian []uint8

func (m *littleEndian) Size() mirv.Address { return mirv.Address(len(*m)) }

func (m *littleEndian) Type() Type { return MemRAM }

func (m *littleEndian) ByteOrder() mirv.ByteOrder { return mirv.LittleEndian }

// Read8 returns the 8 bits little endian  value at address addr.
//
func (m *littleEndian) Read8(addr mirv.Address) (uint8, error) {
	if len((*m)[addr:]) < 1 {
		return 0, errPage
	}
	return (*m)[addr], nil
}

// Write8 writes the 8 bits little endian  value to address addr.
//
func (m *littleEndian) Write8(addr mirv.Address, v uint8) error {
	if len((*m)[addr:]) < 1 {
		return errPage
	}
	(*m)[addr] = v
	return nil
}

// Read16 returns the 16 bits little endian  value at address addr.
//
func (m *littleEndian) Read16(addr mirv.Address) (uint16, error) {
	if len((*m)[addr:]) < 2 {
		return 0, errPage
	}
	return binary.LittleEndian.Uint16((*m)[addr:]), nil
}

// Write16 writes the 16 bits little endian  value to address addr.
//
func (m *littleEndian) Write16(addr mirv.Address, v uint16) error {
	if len((*m)[addr:]) < 2 {
		return errPage
	}
	binary.LittleEndian.PutUint16((*m)[addr:], v)
	return nil
}

// Read32 returns the 32 bits little endian  value at address addr.
//
func (m *littleEndian) Read32(addr mirv.Address) (uint32, error) {
	if len((*m)[addr:]) < 4 {
		return 0, errPage
	}
	return binary.LittleEndian.Uint32((*m)[addr:]), nil
}

// Write32 writes the 32 bits little endian  value to address addr.
//
func (m *littleEndian) Write32(addr mirv.Address, v uint32) error {
	if len((*m)[addr:]) < 4 {
		return errPage
	}
	binary.LittleEndian.PutUint32((*m)[addr:], v)
	return nil
}

// Read64 returns the 64 bits little endian  value at address addr.
//
func (m *littleEndian) Read64(addr mirv.Address) (uint64, error) {
	if len((*m)[addr:]) < 8 {
		return 0, errPage
	}
	return binary.LittleEndian.Uint64((*m)[addr:]), nil
}

// Write64 writes the 64 bits little endian  value to address addr.
//
func (m *littleEndian) Write64(addr mirv.Address, v uint64) error {
	if len((*m)[addr:]) < 8 {
		return errPage
	}
	binary.LittleEndian.PutUint64((*m)[addr:], v)
	return nil
}

// Read8 always returns 0 and an error of type *ErrBus.
//
func (NoMemory) Read8(addr mirv.Address) (uint8, error) {
	return 0, errBus(opRead, 1, addr)
}

// Write8 always returns an error of type *ErrBus.
//
func (NoMemory) Write8(addr mirv.Address, v uint8) error {
	return errBus(opWrite, 1, addr)
}

// Read16 always returns 0 and an error of type *ErrBus.
//
func (NoMemory) Read16(addr mirv.Address) (uint16, error) {
	return 0, errBus(opRead, 2, addr)
}

// Write16 always returns an error of type *ErrBus.
//
func (NoMemory) Write16(addr mirv.Address, v uint16) error {
	return errBus(opWrite, 2, addr)
}

// Read32 always returns 0 and an error of type *ErrBus.
//
func (NoMemory) Read32(addr mirv.Address) (uint32, error) {
	return 0, errBus(opRead, 4, addr)
}

// Write32 always returns an error of type *ErrBus.
//
func (NoMemory) Write32(addr mirv.Address, v uint32) error {
	return errBus(opWrite, 4, addr)
}

// Read64 always returns 0 and an error of type *ErrBus.
//
func (NoMemory) Read64(addr mirv.Address) (uint64, error) {
	return 0, errBus(opRead, 8, addr)
}

// Write64 always returns an error of type *ErrBus.
//
func (NoMemory) Write64(addr mirv.Address, v uint64) error {
	return errBus(opWrite, 8, addr)
}
