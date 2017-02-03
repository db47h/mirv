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

// import "encoding/binary"

// // MemoryLE implements a little-endian memory interface.
// type MemoryLE struct {
// 	tlb
// }

// // NewLE returns a newly initialized MemoryLE
// func NewLE(cacheSize Address) *MemoryLE {
// 	if (cacheSize & (cacheSize - 1)) != 0 {
// 		panic("Cache size must be an exponent of 2.")
// 	}
// 	s := cacheSize >> PageBits
// 	return &MemoryLE{
// 		tlb: tlb{cs: s, cm: s - 1},
// 	}
// }

// // MapMemory maps the memory pages starting at the given address. The memory
// // is immediately allocated.
// func (m *MemoryLE) MapMemory(addr Address, size uint) error {
// 	return m.tlb.mapMemory(addr, size)
// }

// // MapIO maps the specified handler to the given address. See the Controller
// // interface for more details.
// func (m *MemoryLE) MapIO(addr Address, handler Handler) error {
// 	return m.tlb.mapIO(addr, handler)
// }

// // Read8 returns the unsigned 8 bit value at the specified address.
// func (m *MemoryLE) Read8(addr Address) (uint8, error) {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 1 {
// 		return 0, ErrBus
// 	}
// 	return mem[0], nil
// }

// // Read16 returns the unsigned 16 bit value at the specified address.
// func (m *MemoryLE) Read16(addr Address) (uint16, error) {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 2 {
// 		switch handler.(type) {
// 		case noMem:
// 			return 0, ErrBus
// 		default:
// 			return 0, ErrPageAlign
// 		}
// 	}
// 	return binary.LittleEndian.Uint16(mem), nil
// }

// // Read32 returns the unsigned 32 bit value at the specified address.
// func (m *MemoryLE) Read32(addr Address) (uint32, error) {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 4 {
// 		switch handler.(type) {
// 		case noMem:
// 			return 0, ErrBus
// 		default:
// 			return 0, ErrPageAlign
// 		}
// 	}
// 	return binary.LittleEndian.Uint32(mem), nil
// }

// // Read64 returns the unsigned 64 bit value at the specified address.
// func (m *MemoryLE) Read64(addr Address) (uint64, error) {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 8 {
// 		switch handler.(type) {
// 		case noMem:
// 			return 0, ErrBus
// 		default:
// 			return 0, ErrPageAlign
// 		}
// 	}
// 	return binary.LittleEndian.Uint64(mem), nil
// }

// // Write8 writes the unsigned 8 bits value at the specified address.
// func (m *MemoryLE) Write8(addr Address, v uint8) error {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 1 {
// 		return ErrBus
// 	}
// 	mem[0] = v
// 	return nil
// }

// // Write16 writes the unsigned 16 bits value at the specified address.
// func (m *MemoryLE) Write16(addr Address, v uint16) error {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 2 {
// 		switch handler.(type) {
// 		case noMem:
// 			return ErrBus
// 		default:
// 			return ErrPageAlign
// 		}
// 	}
// 	binary.LittleEndian.PutUint16(mem, v)
// 	return nil
// }

// // Write32 writes the unsigned 32 bits value at the specified address.
// func (m *MemoryLE) Write32(addr Address, v uint32) error {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 4 {
// 		switch handler.(type) {
// 		case noMem:
// 			return ErrBus
// 		default:
// 			return ErrPageAlign
// 		}
// 	}
// 	binary.LittleEndian.PutUint32(mem, v)
// 	return nil
// }

// // Write64 writes the unsigned 64 bits value at the specified address.
// func (m *MemoryLE) Write64(addr Address, v uint64) error {
// 	handler := m.handler(addr)
// 	mem := handler.MemAt(addr & _BM)
// 	if len(mem) < 8 {
// 		switch handler.(type) {
// 		case noMem:
// 			return ErrBus
// 		default:
// 			return ErrPageAlign
// 		}
// 	}
// 	binary.LittleEndian.PutUint64(mem, v)
// 	return nil
// }
