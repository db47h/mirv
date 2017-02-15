package mem_test

import (
	"testing"

	"github.com/db47h/mirv"
	"github.com/db47h/mirv/mem"
)

const psz = 1 << 12

func BenchmarkBus_Write64(b *testing.B) {
	var bus mem.Bus
	r := mem.NewRAM(psz, mirv.LittleEndian)
	bus.Map(0, r)
	for i := 0; i < b.N; i++ {
		if err := bus.Write64(0, 12345); err != nil {
			b.Fatal(err)
		}
	}
}
