package sys

import (
	"testing"

	"github.com/db47h/mirv/mem"
)

const psz = 1 << 12

func TestBus_MapPage(t *testing.T) {
	b := NewBus(psz, 1<<8)
	r := mem.BigEndianRAM(psz * 2)
	const ba = 4242 << 20
	// 12 bits page size + 8 bits cache size => 20 bits addressable through cache
	b.Map(ba, r.Pages(b.PageSize())...)
	// check tag for this address
	if tag := b.tag(ba); tag != ba>>12 {
		t.Fatalf("Wrong tag value for base address 0x%d. Got 0x%x, expected 0x%d.", ba, tag, ba>>12)
	}

	if b.Memory(ba-psz) != nil {
		t.Fatal("Address 0 should not be mapped")
	}
	if b.Memory(ba+psz*3) != nil {
		t.Fatalf("Address 0x%x should not be mapped", psz*3)
	}
	// check cache
	if len(b.cache) != 1<<8 {
		t.Fatalf("Wrong cache size: %d, expected %d", len(b.cache), 1<<8)
	}
	for i := range b.cache {
		if b.cache[i] != nil {
			t.Fatalf("Unexpected cache entry %d: %T", i, b.cache[i].m)
		}
	}
	// cause a miss
	err := b.Write8(ba, 42)
	if err != nil {
		t.Fatal(err)
	}
	// check cache
	e := b.cache[b.tag(ba)&0xff]
	if e.tag != b.tag(ba) {
		t.Fatalf("Wrong tag 0x%x, mem type: %T", e.tag, e.m)
	}
	if me := b.Memory(ba); me != e.m {
		t.Fatalf("Memory returned unexpected value %T != %T", me, e.m)
	}
}
