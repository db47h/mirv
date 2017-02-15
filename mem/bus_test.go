package mem

import (
	"testing"
	"testing/quick"

	"github.com/db47h/mirv"
)

const psz = 1 << 12

func TestBus_Map(t *testing.T) {
	var b Bus
	r := NewRAM(psz*2, mirv.LittleEndian)
	const ba = 4242 << 20
	// 12 bits page size + 8 bits cache size => 20 bits addressable through cache
	b.Map(ba, r)

	if _, m := b.Memory(ba - psz); m.Size() != 0 {
		t.Fatal("Address 0 should not be mapped")
	}
	if _, m := b.Memory(ba + psz*3); m.Size() != 0 {
		t.Fatalf("Address 0x%x should not be mapped", psz*3)
	}
	if len(b.b) != 0 {
		t.Fatalf("Wrong cache size: %d, expected %d", len(b.b), 0)
	}
	b.Map(ba+r.Size(), r)
	if len(b.b) != 1 {
		t.Fatalf("Wrong cache size: %d, expected %d", len(b.b), 1)
	}
}

// test paging
func TestBus_Map_overlap(t *testing.T) {
	var b Bus
	r := NewRAM(psz*2, mirv.LittleEndian)
	b.Map(0, r)
	b.Map(psz*2, r) // map again at a different memory location
	for i := mirv.Address(0); i < psz*2; i += 8 {
		err := b.Write64(i, uint64(i))
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := mirv.Address(0); i < psz*2; i += 8 {
		v, err := b.Read64(i)
		if err != nil {
			t.Fatal(err)
		}
		if v != uint64(i) {
			t.Fatalf("At address 0x%x, expected %d, got %d", i, i, v)
		}
	}
	for i := mirv.Address(0); i < psz*2; i += 8 {
		v, err := b.Read64(i + psz*2)
		if err != nil {
			t.Fatal(err)
		}
		if v != uint64(i) {
			t.Fatalf("At address 0x%x, expected %d, got %d", i+psz*2, i, v)
		}
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
	{func(b *Bus, addr mirv.Address) error { return b.Write64(addr, 0) }, 0, 0, 0, 0}, // do not remove this one, it clears AND checks for cross-page boundary errors
	{func(b *Bus, addr mirv.Address) error { return b.Write8(addr, 42) }, 42, 42 << 8, 42 << 24, 42 << 56},
	{func(b *Bus, addr mirv.Address) error { return b.Write16(addr, 0xbeef) }, 0xbe, 0xbeef, 0xbeef << 16, 0xbeef << 48},
	{func(b *Bus, addr mirv.Address) error { return b.Write32(addr, 0xdeadbeef) }, 0xde, 0xdead, 0xdeadbeef, 0xdeadbeef << 32},
	{func(b *Bus, addr mirv.Address) error { return b.Write64(addr, 0xbadc0feedeadbeef) }, 0xba, 0xbadc, 0xbadc0fee, 0xbadc0feedeadbeef},
}

var tdLE = [...]testData{
	{func(b *Bus, addr mirv.Address) error { return b.Write64(addr, 0) }, 0, 0, 0, 0}, // do not remove this one, it clears AND checks for cross-page boundary errors
	{func(b *Bus, addr mirv.Address) error { return b.Write8(addr, 42) }, 42, 42, 42, 42},
	{func(b *Bus, addr mirv.Address) error { return b.Write16(addr, 0xbeef) }, 0xef, 0xbeef, 0xbeef, 0xbeef},
	{func(b *Bus, addr mirv.Address) error { return b.Write32(addr, 0xdeadbeef) }, 0xef, 0xbeef, 0xdeadbeef, 0xdeadbeef},
	{func(b *Bus, addr mirv.Address) error { return b.Write64(addr, 0xbadc0feedeadbeef) }, 0xef, 0xbeef, 0xdeadbeef, 0xbadc0feedeadbeef},
}

func TestBigEndian(t *testing.T) {
	var b Bus
	r := NewRAM(2*psz, mirv.BigEndian)

	b.Map(psz, r) // map after the first page

	// make sure that we have two pages mapped
	if _, err := b.Read8(0); err == nil {
		_, m := b.Memory(0)
		t.Fatalf("Found address 0 mapped to %v", m)
	}
	if _, err := b.Read8(3 * psz); err == nil {
		t.Fatalf("Found address 0x%X mapped", psz*3)
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
			_ = b.Write64(addr, 0)
			err := d.f(&b, addr)
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
			v16, err := b.Read16(addr)
			if err != nil || v16 != d.r16 {
				t.Logf("@0x%x Read16() failed for sample %d: got %d, %v", addr, i, v16, err)
				return false
			}
			v32, err := b.Read32(addr)
			if err != nil || v32 != d.r32 {
				t.Logf("@0x%x Read32() failed for sample %d: got %d, %v", addr, i, v32, err)
				return false
			}
			v64, err := b.Read64(addr)
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
	var b Bus
	r := NewRAM(2*psz, mirv.LittleEndian)

	b.Map(psz, r) // map after the first page

	// make sure that we have two pages mapped
	if _, err := b.Read8(0); err == nil {
		t.Fatal("Found address 0 mapped")
	}
	if _, err := b.Read8(3 * psz); err == nil {
		t.Fatalf("Found address 0x%X mapped", psz*3)
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
			_ = b.Write64(addr, 0)
			err := d.f(&b, addr)
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
			v16, err := b.Read16(addr)
			if err != nil || v16 != d.r16 {
				t.Logf("@0x%x Read16() failed for sample %d: got %d, %v", addr, i, v16, err)
				return false
			}
			v32, err := b.Read32(addr)
			if err != nil || v32 != d.r32 {
				t.Logf("@0x%x Read32() failed for sample %d: got %d, %v", addr, i, v32, err)
				return false
			}
			v64, err := b.Read64(addr)
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
