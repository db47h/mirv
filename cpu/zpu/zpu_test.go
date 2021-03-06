package zpu_test

import (
	"fmt"
	"testing"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/cpu/zpu"
	"github.com/db47h/mirv/elf"
	"github.com/db47h/mirv/mem"
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

type memIO struct {
	mem.Interface
}

func (memIO) Type() mem.Type { return mem.MemIO }

func check(name string, pc interface{}, sp interface{}, tos uint32) error {
	var b mem.Bus
	z := zpu.New(&b)
	b.Map(0, mem.NewRAM(1<<20, z.ByteOrder()))
	b.Map(1<<20, memIO{mem.NewRAM(1<<12, z.ByteOrder())})
	for i := mirv.Address(1 << 20); i < mirv.Address(1<<20+1<<12); i += 4 {
		err := b.Write32(i, 0xDEADBEEF)
		if err != nil {
			panic(err)
		}
	}
	z.Reset()

	var err error
	var entry mirv.Address

	if name != "" {
		_, entry, err = elf.Load(&b, "testdata/"+name+".elf", false)
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

	t, err := b.Read32(z.SP())
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

// Dummy UART. Just intercepts read/writes to MMIO.
// A proper implementation should run in a separate goroutine.
//
type uart struct {
	mem.NoMemory

	txReady byte
	txData  byte

	buf []byte
}

func (u *uart) Size() mirv.Address { return 1 << 12 }
func (*uart) Type() mem.Type       { return mem.MemIO }

// override only 32 bits read/writes. The binary does not make 8 bit accesses,
// so we do not have to care about endianness.
func (u *uart) Read32(addr mirv.Address) (uint32, error) {
	if addr != 0xC {
		return u.NoMemory.Read32(addr)
	}
	return uint32(u.txReady)<<8 | uint32(u.txData), nil
}

func (u *uart) Write32(addr mirv.Address, v uint32) error {
	if addr != 0xC {
		return u.NoMemory.Write32(addr, v)
	}
	u.txData = byte(v)
	u.txReady = byte(v >> 8)
	if u.txReady == 0 {
		u.buf = append(u.buf, u.txData)
		u.txData = 0
		u.txReady = 1
	}
	return nil
}

func TestNew(t *testing.T) {
	uart := uart{txReady: 1, buf: make([]byte, 0, 1024)}
	var b mem.Bus
	z := zpu.New(&b)
	b.Map(0, mem.NewRAM(1<<16, z.ByteOrder())) // 64KiB
	b.Map(0x080A0000, &uart)

	arch, entry, err := elf.Load(&b, "testdata/hello.elf", false)
	if err != nil {
		t.Fatal(err)
	}
	if arch.Machine != elf.MachineZPU {
		t.Fatalf("Unexpected arch %v", arch)
	}
	z.Reset()
	z.SetPC(entry)

	defer func() {
		if err := recover(); err != nil {
			t.Logf("ZPU panic @ %08X", z.PC())
			t.Fatal(err)
		}
	}()

	z.Step(2000000)
	// ts := time.Now()
	// cycles := z.Step(2000000)
	// d := time.Now().Sub(ts)
	// t.Logf("%d cycles / %v -- MIPS: %.3f", cycles, d, (float64(cycles)/1000000.0)/d.Seconds())

	if string(uart.buf) != "Hello, World!" {
		t.Fatalf("Expected \"Hello, World!\", got %q", uart.buf)
	}
	t.Logf("ZPU says: %s", uart.buf)
}
