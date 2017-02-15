package elf_test

import (
	"testing"

	"github.com/db47h/mirv/elf"
	"github.com/db47h/mirv/mem"
)

func TestLoad(t *testing.T) {
	var err error
	var b mem.Bus
	// Do not pre-allocate memory. The ELF loader will do it for us.

	arch, entry, err := elf.Load(&b, "testdata/hello.riscv", true)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("Machine: %v, entry: 0x%X", arch, entry)
	if arch.Data == elf.DataNone {
		t.Fatalf("Unsupported byte order %v", arch.Data)
	}
	if arch.Class == elf.ClassNone {
		t.Fatalf("Unsupported arch class %v", arch.Class)
	}
	var v uint64
	if arch.Class == elf.Class32 {
		var x uint32
		if arch.Data == elf.DataLittle {
			x, err = b.Read32(entry)
		} else {
			x, err = b.Read32(entry)
		}
		v = uint64(x)
	} else {
		if arch.Data == elf.DataLittle {
			v, err = b.Read64(entry)
		} else {
			v, err = b.Read64(entry)
		}
	}
	if err != nil {
		t.Fatalf("Failed to read memory @ 0x%X: %v", entry, err)
	}
	// t.Logf("Data @ 0x%X: 0x%X", entry, v)
	if v != 0x5197 {
		t.Fatalf("Value @ entry point = 0x%X, != 0x%X", v, 0x5197)
	}
}
