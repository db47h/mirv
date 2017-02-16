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

// Package cpu is the root package for CPU implementations.
//
package cpu

import "github.com/db47h/mirv"

// Interface wraps the methods for a CPU.
//
type Interface interface {
	mirv.ByteOrdered

	// Initialize the CPU to a known initial state. This function must be
	// explicitly called before calling Step for the first time. Since most CPUs
	// will examine the memory setup during Reset, it should only be called
	// after all memory has been mapped.
	//
	Reset()

	// Step the simulation forward n cycles. Returns the number of cycles elapsed.
	// This function may return early in some cases (like a HALT or breakpoint instruction).
	//
	Step(n uint64) uint64

	// SetPC set the Program Counter register to the given address.
	//
	SetPC(newPC mirv.Address)

	// PC returns the current value of the Program Counter register.
	//
	PC() mirv.Address

	// SP returns the current value of the Stack Pointer register.
	//
	SP() mirv.Address
}
