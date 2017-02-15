package mem_test

import (
	"fmt"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/mem"
)

// dummy MMIO. We just reuse RAM and change its type.
type ioMem struct {
	mem.Interface
}

func (m *ioMem) Type() mem.Type { return mem.MemIO }

// exercises memory mapping and reported mapped ranges.
//
func ExampleBus_MappedRange() {
	var pageSize mirv.Address = 0x1000 // 4096
	var b mem.Bus
	r256 := mem.NewRAM(pageSize*256, mirv.LittleEndian)
	r2 := mem.NewRAM(pageSize*2, mirv.LittleEndian)
	b.Map(0x40000000, r256)
	b.Map(0x00005000, r2)
	rIO := &ioMem{mem.NewRAM(pageSize*4, mirv.LittleEndian)}
	b.Map(0x10000000, rIO)
	b.Map(0x00001000, rIO)
	b.Map(0x80000000, rIO)
	l, h, _ := b.MappedRange(mem.MemRAM)
	fmt.Printf("RAM: 0x%X - 0x%X\n", l, h)
	l, h, _ = b.MappedRange(mem.MemIO)
	fmt.Printf("IO : 0x%X - 0x%X\n", l, h)

	// now map the last 2 pages
	addr := 0 - pageSize*2
	b.Map(addr, r2)
	// here MemRange will return a high value of 0
	// because of 2 complement arithmetic.
	l, h, _ = b.MappedRange(mem.MemRAM)
	fmt.Printf("RAM: 0x%X - 0x%X\n", l, h)

	// Output:
	// RAM: 0x5000 - 0x40100000
	// IO : 0x1000 - 0x80004000
	// RAM: 0x5000 - 0x0
}
