package mirv

import "errors"

// Address is the guest address type.
type Address uint

// Memory wraps the methods that memory-like devices implement. The address
// passed to the methods is relative to the start of the page to which the
// Memory implementation is mapped.
//
// Read/Writes to unmapped memory result in a *ErrBus error.
//
type Memory interface {
	Read8(addr Address) (uint8, error)   // Read unsigned 8 bit value from address.
	Read16(addr Address) (uint16, error) // Read unsigned 16 bit value from address.
	Read32(addr Address) (uint32, error) // Read unsigned 32 bit value from address.
	Read64(addr Address) (uint64, error) // Read unsigned 64 bit value from address.

	Write8(addr Address, v uint8) error   // Write unsigned 8 bit value to address.
	Write16(addr Address, v uint16) error // Write unsigned 16 bit value to address.
	Write32(addr Address, v uint32) error // Write unsigned 32 bit value to address.
	Write64(addr Address, v uint64) error // Write unsigned 64 bit value to address.
}

// Errors
var (
	ErrAlign     = errors.New("Unaligned read/write")
	ErrCrossPage = errors.New("Read/write across page boundary")
)

// A ByteOrder specifies how to convert byte sequences into 16-, 32-, or 64-bit unsigned integers.
// Error values can only be ErrCrossPage.
//
type ByteOrder interface {
	Uint16([]byte) (uint16, error)
	Uint32([]byte) (uint32, error)
	Uint64([]byte) (uint64, error)
	PutUint16([]byte, uint16) error
	PutUint32([]byte, uint32) error
	PutUint64([]byte, uint64) error
}

// Byte order helpers. Both implement the ByteOrder interface.
var (
	LittleEndian = littleEndian{}
	BigEndian    = bigEndian{}
)

type littleEndian struct{}

func (littleEndian) Uint16(b []byte) (uint16, error) {
	if len(b) < 2 {
		return 0, ErrCrossPage
	}
	return uint16(b[0]) | uint16(b[1])<<8, nil
}

func (littleEndian) PutUint16(b []byte, v uint16) error {
	if len(b) < 2 {
		return ErrCrossPage
	}
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	return nil
}

func (littleEndian) Uint32(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, ErrCrossPage
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

func (littleEndian) PutUint32(b []byte, v uint32) error {
	if len(b) < 4 {
		return ErrCrossPage
	}
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return nil
}

func (littleEndian) Uint64(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, ErrCrossPage
	}
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56, nil
}

func (littleEndian) PutUint64(b []byte, v uint64) error {
	if len(b) < 8 {
		return ErrCrossPage
	}
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	return nil
}

type bigEndian struct{}

func (bigEndian) Uint16(b []byte) (uint16, error) {
	if len(b) < 2 {
		return 0, ErrCrossPage
	}
	return uint16(b[1]) | uint16(b[0])<<8, nil
}

func (bigEndian) PutUint16(b []byte, v uint16) error {
	if len(b) < 2 {
		return ErrCrossPage
	}
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return nil
}

func (bigEndian) Uint32(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, ErrCrossPage
	}
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24, nil
}

func (bigEndian) PutUint32(b []byte, v uint32) error {
	if len(b) < 4 {
		return ErrCrossPage
	}
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
	return nil
}

func (bigEndian) Uint64(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, ErrCrossPage
	}
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56, nil
}

func (bigEndian) PutUint64(b []byte, v uint64) error {
	if len(b) < 8 {
		return ErrCrossPage
	}
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
	return nil
}
