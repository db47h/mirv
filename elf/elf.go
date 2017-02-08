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
// TODO:
//
//	- rename this package to something more generic and related to image loading.
//	- add a loadRaw function
//	- kernel loading
//
package elf

import (
	self "debug/elf"
	"fmt"
	"io"
	"strconv"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/mem"
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
	ClassNone Class = iota // Unknown class.
	Class32                // 32-bit architecture.
	Class64                // 64-bit architecture.
)

//go:generate stringer -type Data "$GOFILE"

// Data values.
const (
	DataNone   Data = iota // Unknown data format.
	DataLittle             // 2's complement little-endian.
	DataBig                // 2's complement big-endian.
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

type zeroReader struct{}

func (zeroReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// alloc allocates and maps memory @addr with the given size.
// it allocates only the necessary pages.
//
func alloc(b *sys.Bus, addr, size mirv.Address) {
	ps := b.PageSize()
	pm := ps - 1
	// adjust addr & size
	size += addr & pm
	size = (size + pm) & ^pm
	addr &= ^pm
	var start, cur, end mirv.Address = 0, addr, addr + size

	// Try to allocate in large chunks instead of allocating page by page.
	for cur != end {
		// look for first unmapped page
		for start = cur; b.Memory(start).Size() != 0 && start != end; start += ps {
		}
		if start == end {
			break
		}
		// look for next mapped page
		for cur = start + ps; b.Memory(cur).Size() == 0 && cur != end; cur += ps {
		}
		b.Map(start, mem.New(cur-start))
	}
}

// Load loads an ELF file and returns the architecture, start address and error
// if any. If the autoAlloc parameter is true, guest memory will automatically
// be allocated and mapped in the guest's address space.
//
// Note that there is no memory access control mechanism in the current version.
// However this sill be implemented in future versions and auto-allocation will
// also auto-configure memory access control for all loaded segments (even for
// memory pages allocated and mapped manually before calling Load).
//
// The loader is rather primitive and has some limitations:
//
// Panics on files with program segments larger or equal to 0x8000000000000000
// bytes.
//
// Only statically linked executables are supported.
//
func Load(name string, bus *sys.Bus, autoAlloc bool) (arch Arch, entry mirv.Address, err error) {
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

	if f.Type != self.ET_EXEC {
		return arch, entry, fmt.Errorf("unsupported elf file type %v", f.Type)
	}

	for _, p := range f.Progs {
		if p.Type != self.PT_LOAD {
			return arch, entry, fmt.Errorf("unsupported prog type %v", p.Type)
		}
		r := p.Open()
		n := int64(p.Filesz)
		if n < 0 || int64(p.Memsz) < 0 {
			panic("ELF file too large")
		}
		if autoAlloc {
			alloc(bus, mirv.Address(p.Paddr), mirv.Address(p.Memsz))
		}
		w := bus.Writer(mirv.Address(p.Paddr))
		n, err := io.CopyN(w, r, int64(p.Filesz))
		if err != nil {
			return arch, entry, err
		}
		// zero-fill the gap between p.Filesz and p.Memsz
		// This is to conform to the ELF spec. As a side effect, this clears the
		// BSS, but this should not be taken for granted.
		if uint64(n) < p.Memsz {
			_, err = io.CopyN(w, zeroReader{}, int64(p.Memsz)-n)
			if err != nil {
				return arch, entry, err
			}
		}
	}

	return arch, entry, nil
}
