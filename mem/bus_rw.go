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
	"github.com/db47h/mirv"
)

// Read8 returns the 8 bits value at address addr.
//
func (b *Bus) Read8(addr mirv.Address) (uint8, error) {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Read8(addr - blk.s)
}

// Write8 writes the 8 bits value to address addr.
//
func (b *Bus) Write8(addr mirv.Address, v uint8) error {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Write8(addr - blk.s, v)
}

// Read16 returns the 16 bits value at address addr.
//
func (b *Bus) Read16(addr mirv.Address) (uint16, error) {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Read16(addr - blk.s)
}

// Write16 writes the 16 bits value to address addr.
//
func (b *Bus) Write16(addr mirv.Address, v uint16) error {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Write16(addr - blk.s, v)
}

// Read32 returns the 32 bits value at address addr.
//
func (b *Bus) Read32(addr mirv.Address) (uint32, error) {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Read32(addr - blk.s)
}

// Write32 writes the 32 bits value to address addr.
//
func (b *Bus) Write32(addr mirv.Address, v uint32) error {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Write32(addr - blk.s, v)
}

// Read64 returns the 64 bits value at address addr.
//
func (b *Bus) Read64(addr mirv.Address) (uint64, error) {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Read64(addr - blk.s)
}

// Write64 writes the 64 bits value to address addr.
//
func (b *Bus) Write64(addr mirv.Address, v uint64) error {
	blk := b.p
	if !blk.contains(addr) {
		blk = b.find(addr)
	}
	return blk.m.Write64(addr - blk.s, v)
}

