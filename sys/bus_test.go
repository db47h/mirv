package sys

import (
	"testing"
	"testing/quick"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/mem"
)

const psz = 1 << 12

func TestBus_Map(t *testing.T) {
	b := NewBus(psz, 1<<8)
	r := mem.New(psz * 2)
	const ba = 4242 << 20
	// 12 bits page size + 8 bits cache size => 20 bits addressable through cache
	b.Map(ba, r, MemRAM)
	// check tag for this address
	if tag := b.tag(ba); tag != ba>>12 {
		t.Fatalf("Wrong tag value for base address 0x%d. Got 0x%x, expected 0x%d.", ba, tag, ba>>12)
	}

	if b.Memory(ba-psz).Size() != 0 {
		t.Fatal("Address 0 should not be mapped")
	}
	if b.Memory(ba+psz*3).Size() != 0 {
		t.Fatalf("Address 0x%x should not be mapped", psz*3)
	}
	// check cache
	if len(b.cache) != 1<<8 {
		t.Fatalf("Wrong cache size: %d, expected %d", len(b.cache), 1<<8)
	}
	for i := range b.cache {
		if _, ok := b.cache[i].m.(nilMemory); !ok {
			t.Fatalf("Unexpected cache entry %d: %T", i, b.cache[i].m)
		}
	}

	// page check
	var ta mirv.Address = ba + psz
	v8, err := b.Read8(ta)
	if err != nil || v8 != 0 {
		t.Fatalf("Unexpected Read8(0x%d) result: %d, %v", ta, v8, err)
	}
	if err = b.Write8(ta, 42); err != nil {
		t.Fatal(err)
	}
	// test various methods of getting back our value
	if v8 = r.Get(psz)[0]; v8 != 42 {
		t.Errorf("r.Get returned %d", v8)
	}
	if v8 = r.Page(psz, psz).Get(0)[0]; v8 != 42 {
		t.Errorf("Sub-page Get returned %d", v8)
	}
	if v8 = b.Memory(ta).Get(0)[0]; v8 != 42 {
		t.Errorf("Memory returned %d", v8)
	}

	// check cache
	e := b.cache[b.tag(ta)&0xff]
	if e.tag != b.tag(ta) {
		t.Fatalf("Wrong tag 0x%x, mem type: %T", e.tag, e.m)
	}
}

// test paging
func TestBus_Map_overlap(t *testing.T) {
	b := NewBus(psz, 1<<8)
	r := mem.New(psz * 2)
	b.Map(0, r, MemRAM)
	b.Map(psz*2, r, MemRAM) // map again at a different memory location
	for i := mirv.Address(0); i < psz*2; i += 8 {
		err := b.Write64LE(i, uint64(i))
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := mirv.Address(0); i < psz*2; i += 8 {
		v, err := b.Read64LE(i)
		if err != nil {
			t.Fatal(err)
		}
		if v != uint64(i) {
			t.Fatalf("At address 0x%x, expected %d, got %d", i, i, v)
		}
	}
	for i := mirv.Address(0); i < psz*2; i += 8 {
		v, err := b.Read64LE(i + psz*2)
		if err != nil {
			t.Fatal(err)
		}
		if v != uint64(i) {
			t.Fatalf("At address 0x%x, expected %d, got %d", i+psz*2, i, v)
		}
	}
}

func TestBus_RAM(t *testing.T) {
	b := NewBus(1<<12, 8)
	r0 := mem.New((1 << 12) << 2)
	r1 := mem.New((1 << 12) << 8)
	r2 := mem.New((1 << 12) << 1)
	b.Map(0x00001000, r0, MemIO)
	b.Map(0x40000000, r1, MemRAM)
	b.Map(0x10000000, r0, MemIO)
	b.Map(0x00005000, r2, MemRAM)
	b.Map(0x80000000, r0, MemIO)
	l, h := b.MemRange(MemRAM)
	if l != 0x5000 && h != 0x40000000+0x10000-1 {
		t.Fatalf("Got 0x%X, 0x%X, expected 0x%X, 0x%X", l, h, 0x5000, 0x40000000+0x10000-1)
	}
	l, h = b.MemRange(MemIO)
	if l != 0x1000 && h != 0x80000000+0x4000-1 {
		t.Fatalf("Got 0x%X, 0x%X, expected 0x%X, 0x%X", l, h, 0x1000, 0x80000000+0x4000-1)
	}
}

type testData struct {
	f   func(*Bus, mirv.Address) error
	r8  uint8
	r16 uint16
	r32 uint32
	r64 uint64
}

var tdBE = [...]testData{
	{func(b *Bus, addr mirv.Address) error { return b.Write64BE(addr, 0) }, 0, 0, 0, 0}, // do not remove this one, it clears AND checks for cross-page boundary errors
	{func(b *Bus, addr mirv.Address) error { return b.Write8(addr, 42) }, 42, 42 << 8, 42 << 24, 42 << 56},
	{func(b *Bus, addr mirv.Address) error { return b.Write16BE(addr, 0xbeef) }, 0xbe, 0xbeef, 0xbeef << 16, 0xbeef << 48},
	{func(b *Bus, addr mirv.Address) error { return b.Write32BE(addr, 0xdeadbeef) }, 0xde, 0xdead, 0xdeadbeef, 0xdeadbeef << 32},
	{func(b *Bus, addr mirv.Address) error { return b.Write64BE(addr, 0xbadc0feedeadbeef) }, 0xba, 0xbadc, 0xbadc0fee, 0xbadc0feedeadbeef},
}

var tdLE = [...]testData{
	{func(b *Bus, addr mirv.Address) error { return b.Write64LE(addr, 0) }, 0, 0, 0, 0}, // do not remove this one, it clears AND checks for cross-page boundary errors
	{func(b *Bus, addr mirv.Address) error { return b.Write8(addr, 42) }, 42, 42, 42, 42},
	{func(b *Bus, addr mirv.Address) error { return b.Write16LE(addr, 0xbeef) }, 0xef, 0xbeef, 0xbeef, 0xbeef},
	{func(b *Bus, addr mirv.Address) error { return b.Write32LE(addr, 0xdeadbeef) }, 0xef, 0xbeef, 0xdeadbeef, 0xdeadbeef},
	{func(b *Bus, addr mirv.Address) error { return b.Write64LE(addr, 0xbadc0feedeadbeef) }, 0xef, 0xbeef, 0xdeadbeef, 0xbadc0feedeadbeef},
}

func TestBigEndian(t *testing.T) {
	b := NewBus(psz, 1<<10)
	r := mem.New(2*psz + 2)

	b.Map(psz, r, MemRAM) // map after the first page

	// make sure that we have two pages mapped
	if _, err := b.Read8(0); err == nil {
		t.Fatal("Found address 0 mapped")
	}
	if _, err := b.Read8(3 * psz); err == nil {
		t.Fatalf("Found address %x mapped", psz*3)
	}

	tf := func(addr16 uint16) bool {
		addr := mirv.Address(addr16 >> 2)
		if addr < psz || addr >= psz*3 {
			// should be unmapped.
			if _, err := b.Read8(addr); err == nil {
				t.Logf("Unexpected success reading unmapped address %d", addr)
				return false
			}
			return true
		}
		for i, d := range tdBE {
			_ = b.Write64BE(addr, 0)
			err := d.f(b, addr)
			if err != nil {
				if addr&7 != 0 {
					return true
				}
				t.Logf("@0x%x f() failed for sample %d: %v", addr, i, err)
				return false
			}
			v8, err := b.Read8(addr)
			if err != nil || v8 != d.r8 {
				t.Logf("@0x%x Read8() failed for sample %d: got %d, %v", addr, i, v8, err)
				return false
			}
			v16, err := b.Read16BE(addr)
			if err != nil || v16 != d.r16 {
				t.Logf("@0x%x Read16() failed for sample %d: got %d, %v", addr, i, v16, err)
				return false
			}
			v32, err := b.Read32BE(addr)
			if err != nil || v32 != d.r32 {
				t.Logf("@0x%x Read32() failed for sample %d: got %d, %v", addr, i, v32, err)
				return false
			}
			v64, err := b.Read64BE(addr)
			if err != nil || v64 != d.r64 {
				t.Logf("@0x%x Read64() failed for sample %d: got %d, %v", addr, i, v64, err)
				return false
			}
		}
		return true
	}
	err := quick.Check(tf, &quick.Config{MaxCount: 65536})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLittleEndian(t *testing.T) {
	b := NewBus(psz, 1<<10)
	r := mem.New(2*psz + 2)

	b.Map(psz, r, MemRAM) // map after the first page

	// make sure that we have two pages mapped
	if _, err := b.Read8(0); err == nil {
		t.Fatal("Found address 0 mapped")
	}
	if _, err := b.Read8(3 * psz); err == nil {
		t.Fatalf("Found address %x mapped", psz*3)
	}

	tf := func(addr16 uint16) bool {
		addr := mirv.Address(addr16 >> 2)
		if addr < psz || addr >= psz*3 {
			// should be unmapped.
			if _, err := b.Read8(addr); err == nil {
				t.Logf("Unexpected success reading unmapped address %d", addr)
				return false
			}
			return true
		}
		for i, d := range tdLE {
			_ = b.Write64LE(addr, 0)
			err := d.f(b, addr)
			if err != nil {
				if addr&7 != 0 {
					return true
				}
				t.Logf("@0x%x f() failed for sample %d: %v", addr, i, err)
				return false
			}
			v8, err := b.Read8(addr)
			if err != nil || v8 != d.r8 {
				t.Logf("@0x%x Read8() failed for sample %d: got %d, %v", addr, i, v8, err)
				return false
			}
			v16, err := b.Read16LE(addr)
			if err != nil || v16 != d.r16 {
				t.Logf("@0x%x Read16() failed for sample %d: got %d, %v", addr, i, v16, err)
				return false
			}
			v32, err := b.Read32LE(addr)
			if err != nil || v32 != d.r32 {
				t.Logf("@0x%x Read32() failed for sample %d: got %d, %v", addr, i, v32, err)
				return false
			}
			v64, err := b.Read64LE(addr)
			if err != nil || v64 != d.r64 {
				t.Logf("@0x%x Read64() failed for sample %d: got %d, %v", addr, i, v64, err)
				return false
			}
		}
		return true
	}
	err := quick.Check(tf, &quick.Config{MaxCount: 65536})
	if err != nil {
		t.Fatal(err)
	}
}
