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

// Package elf provides utility functions to load ELF files.
//
package elf

import (
	self "debug/elf"
	"fmt"
	"io"
	"strconv"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/sys"
)

// Class corresponds to ELF Header.Ident[EI_CLASS] and Header.Class.
type Class byte

// Data corresponds to ELF Header.Ident[EI_DATA] and Header.Data.
type Data byte

// Machine corresponds to ELF Header.Machine.
type Machine uint16

//go:generate stringer -type Class "$GOFILE"

// Class values.
const (
	ClassNone Class = iota /* Unknown class. */
	Class32                /* 32-bit architecture. */
	Class64                /* 64-bit architecture. */
)

//go:generate stringer -type Data "$GOFILE"

// Data values.
const (
	DataNone   Data = iota /* Unknown data format. */
	DataLittle             /* 2's complement little-endian. */
	DataBig                /* 2's complement big-endian. */
)

// Supported Machine IDs
const (
	MachineZPU   Machine = 106
	MachineLM32  Machine = 138
	MachineRISCV Machine = 243
)

// machine names
var machineNames = []struct {
	mach Machine
	name string
}{
	{MachineZPU, "zpu"},
	{MachineLM32, "lm32"},
	{MachineRISCV, "riscv"},
}

func (m Machine) String() string {
	for _, n := range machineNames {
		if n.mach == m {
			return n.name
		}
	}
	return "unknown-" + strconv.Itoa(int(m))
}

// Arch wraps an architecture description.
//
type Arch struct {
	Machine Machine
	Class   Class
	Data    Data
}

type zeroReader uint64

func (z *zeroReader) Read(p []byte) (n int, err error) {
	if *z == 0 {
		return 0, io.EOF
	}
	var i int
	for i = 0; i < len(p) && *z != 0; i++ {
		p[i] = 0
		*z--
	}
	return i, nil
}

// Load loads an ELF file and returns the architecture, start address and error if any.
//
func Load(name string, bus *sys.Bus) (arch Arch, entry mirv.Address, err error) {
	f, err := self.Open(name)
	if err != nil {
		return Arch{}, 0, err
	}
	defer f.Close()

	entry = mirv.Address(f.Entry)
	arch = Arch{
		Machine: Machine(f.Machine),
		Class:   Class(f.Class),
		Data:    Data(f.Data),
	}

	for _, p := range f.Progs {
		if p.Type != self.PT_LOAD {
			return arch, entry, fmt.Errorf("unsupported prog type %v", p.Type)
		}
		r := p.Open()
		w := bus.WriteSeeker()
		_, err := w.Seek(int64(p.Paddr), io.SeekStart) // TODO: int64(p.Paddr) may not work
		if err != nil {
			return arch, entry, err
		}
		n, err := io.Copy(w, r)
		if err != nil {
			return arch, entry, err
		}
		if n < 0 {
			panic("Negative write count")
		}
		// zero-fill
		if uint64(n) < p.Memsz {
			z := zeroReader(p.Memsz - uint64(n))
			_, err = io.Copy(w, &z)
			if err != nil {
				return arch, entry, err
			}
		}
	}

	return arch, entry, nil
}
