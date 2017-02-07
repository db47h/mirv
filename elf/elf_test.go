package elf_test

import (
	"testing"

	"github.com/db47h/mirv/elf"
	"github.com/db47h/mirv/mem"
	"github.com/db47h/mirv/sys"
)

func TestLoad(t *testing.T) {
	var err error
	b := sys.NewBus(1<<12, 1<<13)
	r := mem.New(1 << 25)
	b.Map(0, r)

	arch, entry, err := elf.Load("/home/denis/devel/hello", b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Machine: %v, entry: 0x%X", arch, entry)
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
			x, err = b.Read32LE(entry)
		} else {
			x, err = b.Read32BE(entry)
		}
		v = uint64(x)
	} else {
		if arch.Data == elf.DataLittle {
			v, err = b.Read64LE(entry)
		} else {
			v, err = b.Read64BE(entry)
		}
	}
	if err != nil {
		t.Fatalf("Failed to read memory @ 0x%X: %v", entry, err)
	}
	t.Logf("Data @ 0x%X: 0x%X", entry, v)
}
