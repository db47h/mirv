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
package mem

import "github.com/db47h/mirv"

type memory []uint8

func (m memory) Size() mirv.Address {
	return mirv.Address(len(m))
}

func (m memory) Get(addr mirv.Address) []uint8 {
	return m[addr:]
}

func (m memory) Page(addr, size mirv.Address) mirv.Memory {
	if addr == 0 && size <= mirv.Address(len(m)) {
		return m
	}
	return m[addr : addr+size]

}

// New returns a new memory block of the requested size.
//
func New(size mirv.Address) mirv.Memory {
	return make(memory, size)
}
