package mem_test

import (
	"testing"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/mem"
)

type testData struct {
	f   func(*mem.Bus, mirv.Address) error
	r8  uint8
	r16 uint16
	r32 uint32
	r64 uint64
}

const psz = 1 << 12

// var tdBE = [...]testData{
// 	{func(b *mem.Bus, addr mirv.Address) error { return b.Write64BE(addr, 0) }, 0, 0, 0, 0}, // do not remove this one, it clears AND checks for cross-page boundary errors
// 	{func(b *mem.Bus, addr mirv.Address) error { return b.Write8(addr, 42) }, 42, 42 << 8, 42 << 24, 42 << 56},
// 	{func(b *mem.Bus, addr mirv.Address) error { return b.Write16BE(addr, 0xbeef) }, 0xbe, 0xbeef, 0xbeef << 16, 0xbeef << 48},
// 	{func(b *mem.Bus, addr mirv.Address) error { return b.Write32BE(addr, 0xdeadbeef) }, 0xde, 0xdead, 0xdeadbeef, 0xdeadbeef << 32},
// 	{func(b *mem.Bus, addr mirv.Address) error { return b.Write64BE(addr, 0xbadc0feedeadbeef) }, 0xba, 0xbadc, 0xbadc0fee, 0xbadc0feedeadbeef},
// }

// func TestBigEndianRAM(t *testing.T) {
// 	b := mem.NewBus(psz, 1<<10)
// 	r := mem.BigEndianRAM(2*psz + 2)

// 	b.Map(psz, r.Pages(psz)...) // map after the first page

// 	// make sure that we have two pages mapped
// 	if b.Memory(0) != nil {
// 		t.Fatal("Found address 0 mapped")
// 	}
// 	if b.Memory(3) != nil {
// 		t.Fatalf("Found address %x mapped", psz*3)
// 	}

// 	tf := func(addr16 uint16) bool {
// 		addr := mirv.Address(addr16 >> 2)
// 		if addr < psz || addr >= psz*3 {
// 			// should be unmapped.
// 			if _, err := b.Read8(addr); err == nil {
// 				t.Logf("Unexpected success reading unmapped address %d", addr)
// 				return false
// 			}
// 			return true
// 		}
// 		for i, d := range tdBE {
// 			_ = b.Write64(addr, 0)
// 			err := d.f(b, addr)
// 			if err != nil {
// 				if err == mirv.ErrCrossPage {
// 					// is that so?
// 					if addr&7 != 0 {
// 						return true
// 					}
// 				}
// 				t.Logf("@0x%x f() failed for sample %d: %v", addr, i, err)
// 				return false
// 			}
// 			v8, err := b.Read8(addr)
// 			if err != nil || v8 != d.r8 {
// 				t.Logf("@0x%x Read8() failed for sample %d: got %d, %v", addr, i, v8, err)
// 				return false
// 			}
// 			v16, err := b.Read16(addr)
// 			if err != nil || v16 != d.r16 {
// 				t.Logf("@0x%x Read16() failed for sample %d: got %d, %v", addr, i, v16, err)
// 				return false
// 			}
// 			v32, err := b.Read32(addr)
// 			if err != nil || v32 != d.r32 {
// 				t.Logf("@0x%x Read32() failed for sample %d: got %d, %v", addr, i, v32, err)
// 				return false
// 			}
// 			v64, err := b.Read64(addr)
// 			if err != nil || v64 != d.r64 {
// 				t.Logf("@0x%x Read64() failed for sample %d: got %d, %v", addr, i, v64, err)
// 				return false
// 			}
// 		}
// 		return true
// 	}
// 	err := quick.Check(tf, &quick.Config{MaxCount: 65536})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func BenchmarkBus_Write64(b *testing.B) {
	r := mem.New(psz)
	bus := mem.NewBus(psz, 1<<8)
	bus.Map(0, r, mem.MemRAM)
	for i := 0; i < b.N; i++ {
		if err := bus.Write64LE(0, 12345); err != nil {
			b.Fatal(err)
		}
	}
}
