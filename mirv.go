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

// Package mirv is the root package for MIRV. It also provides common types
// used by all sub-packages.
//
// The typical setup of a simulation is:
//
//	var bus mem.Bus								// create a memory bus
//	cpu := lm32.New(&bus)						// select preferred CPU
//	sram := mem.NewRAM(1<<20, cpu.ByteOrder)	// RAM
//	pic := iodev.Pic(cpu.ByteOrder)				// IO devices
//	bus.Map(0, sram)							// Map RAM
//	bus.Map(0x80000000, pic)					// Map IO
//	cpu.Reset()									// Reset CPU
//	for {
//		cpu.Step(1000000)						// Run it
//	}
//
// Note that the memory Interface if byte order sensitive. While this
// may seem counter intuitive this was a necessary evil in order to
// limit the number of methods in the interface, as well as making
// it easier to implement IO devices.
//
// Memory access is being the hotest code path, it's been carefully designed to
// use static dispatching of methods as much as possible (i.e. almost no Go
// interfaces) and allow aggressive inlining.
//
// Regarding inlining and the memory interface, this is one situation where
// generics could have been useful. To work around this limitation, the RAM and
// Bus memory interfaces are automatically generated in order to do some manual
// inlining without error prone copy/paste. Memory performance can be almost
// doubled by compiling with -gcflags "-l -l -l -l".
//
package mirv

// Address is the guest address type.
//
type Address uint

// ByteOrder or endianness.
//
type ByteOrder uint8

// ByteOrder values.
//
const (
	LittleEndian ByteOrder = iota + 1 // make it match with ELF
	BigEndian
)

// ByteOrdered is the interface implemented by all ByteOrder-aware types.
//
type ByteOrdered interface {
	ByteOrder() ByteOrder
}
