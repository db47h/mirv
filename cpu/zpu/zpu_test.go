package zpu_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/cpu/zpu"
	"github.com/db47h/mirv/elf"
	"github.com/db47h/mirv/mem"
	"github.com/db47h/mirv/sys"
)

type data struct {
	n   string
	pc  interface{}
	sp  interface{}
	tos uint32
}

const start = 0x20
const top = 1 << 20
const db = 0xDEADBEEF

var td = [...]data{
	{"", 0, top, db},
	{"im0", start + 1, top - 4, 0},
	{"im-1", start + 1, top - 4, ^uint32(0)},
	{"im7abc0123", start + 5, top - 4, 0x7abc0123},
	{"poppc", start + 32, top, db},
	{"pushsp", nil, top - 4, top},
	{"add", nil, top - 4, 0xABCD1233},
	{"load", nil, top - 4, 0xABCD0123},
	{"store", nil, top, 0xABCD0123},
	{"popsp", nil, start + 32, 0xABCD1234},
	{"and", nil, top - 4, 0x0000aaaa},
	{"or", nil, top - 4, 0xffffffff},
	{"not", nil, top - 4, 0xaaaaaaaa},
	{"flip", nil, top - 4, 0xF77DB57B},

	{"loadsp", nil, top - 12, 0xABCD0123},
	{"storesp", nil, top - 4, 0x34567890},
	{"addsp", nil, top - 8, 0x3456788F},
	{"emul0", 0, top - 4, start + 1},
	{"emul31", 31 * 32, top - 4, start + 1},

	{"swap", nil, top - 4, 0xBEEFDEAD},
}

func check(name string, pc interface{}, sp interface{}, tos uint32) error {
	b := sys.NewBus(1<<12, 1<<8)
	b.Map(0, mem.New(1<<20), sys.MemRAM)
	b.Map(1<<20, mem.New(1<<12), sys.MemIO)
	for i := mirv.Address(1 << 20); i < mirv.Address(1<<20+1<<12); i += 4 {
		err := b.Write32BE(i, 0xDEADBEEF)
		if err != nil {
			panic(err)
		}
	}
	z := zpu.New(b)
	z.Reset()

	var err error
	var entry mirv.Address

	if name != "" {
		_, entry, err = elf.Load(b, "testdata/"+name+".elf", false)
		if err != nil {
			return err
		}
	}

	z.SetPC(entry)
	z.Step(1000)

	switch pc := pc.(type) {
	case int:
		if mirv.Address(pc) != z.PC() {
			return fmt.Errorf("%v: expected PC %08X, got %08X", name, mirv.Address(pc), z.PC())
		}
	default:
		if pc != nil {
			return fmt.Errorf("%v: unexpected type for PC: %T", name, pc)
		}
	}

	switch sp := sp.(type) {
	case int:
		if mirv.Address(sp) != z.SP() {
			return fmt.Errorf("%v: expected SP %08X, got %08X", name, mirv.Address(sp), z.SP())
		}
	default:
		if sp != nil {
			return fmt.Errorf("%v: unexpected type for SP: %T", name, sp)
		}
	}

	t, err := b.Read32BE(z.SP())
	if err != nil {
		return fmt.Errorf("error while reading SP @ %08X: %v", z.SP(), err)
	}
	if t != tos {
		return fmt.Errorf("Expected TOS = %08X, got %08X", tos, t)
	}

	return nil
}

func TestISA(t *testing.T) {
	for _, d := range td {
		err := check(d.n, d.pc, d.sp, d.tos)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestNew(t *testing.T) {
	b := sys.NewBus(1<<12, 1<<8)
	b.Map(0, mem.New(1<<16), sys.MemRAM) // 64KiB
	b.Map(0x080A0000, mem.New(1<<12), sys.MemIO)

	arch, entry, err := elf.Load(b, "testdata/hello.elf", false)
	if err != nil {
		t.Fatal(err)
	}
	if arch.Machine != elf.MachineZPU {
		t.Fatalf("Unexpected arch %v", arch)
	}
	z := zpu.New(b)
	z.Reset()
	z.SetPC(entry)

	defer func() {
		if err := recover(); err != nil {
			t.Logf("ZPU panic @ %08X", z.PC())
			t.Fatal(err)
		}
	}()

	// Dummy UART. The newlib implementation of outbyte() does a direct write to
	// IO mem and our bus implementation dos not enable read/write hijacking, so
	// we need to do something bad. TODO: FIX THIS!
	var buf []byte
	var uartDone = make(chan struct{})
	b.Write32BE(0x080A000C, 0x100) // tx ready

	go func(done <-chan struct{}) {
		// we start with 0x080A000C == 0x100 (tx ready)
		// outbyte(byteVal) does a Write32BE(0x080A000C, byteVal)
		t := time.NewTicker(time.Second / (19200 / 8))
		for {
			select {
			case <-t.C:
				c, _ := b.Read16BE(0x080A000C + 2)
				if c&0x100 != 0 {
					// nothing written, skip this tick
					continue
				}
				// read until we get the same value twice
				for cc, _ := b.Read16BE(0x080A000C + 2); cc != c; c = cc {
				}
				// xmit
				buf = append(buf, byte(c))
				// clear byte
				_ = b.Write8(0x080A000C+3, 0)
				// tx ready
				_ = b.Write8(0x080A000C+2, 1)
			case <-done:
				t.Stop()
				return
			}
		}
	}(uartDone)

	z.Step(2000000)

	// wait for the UART to read the last byte
	time.Sleep(2 * time.Second / (19200 / 8))
	close(uartDone)

	if string(buf) != "Hello, World!" {
		t.Fatalf("Expected \"Hello, World!\", got %q", buf)
	}
	t.Logf("ZPU says: %s", buf)
}
