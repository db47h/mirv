// Copyright 2017 Denis Bernard <db047h@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package zpu implements the Zylin ZPU ISA. See http://opencores.org/project,zpu
//
package zpu

import (
	"github.com/db47h/mirv"
	"github.com/db47h/mirv/cpu"
	"github.com/db47h/mirv/mem"
)

type opcode uint8

const (
	opBreakPoint opcode = 0x00
	opIM         opcode = 0x80
	opStoreSP    opcode = 0x40
	opLoadSP     opcode = 0x60
	opAddSP      opcode = 0x10
	opEmulate    opcode = 0x20
	opPopPC      opcode = 0x04
	opLoad       opcode = 0x08
	opStore      opcode = 0x0C
	opPushSP     opcode = 0x02
	opPopSP      opcode = 0x0D
	opAdd        opcode = 0x05
	opAnd        opcode = 0x06
	opOr         opcode = 0x07
	opNot        opcode = 0x09
	opFlip       opcode = 0x0A
	opNop        opcode = 0x0B

	opSwap    opcode = 40
	opSyscall opcode = 60
)

const (
	opIMMask      opcode = 0x80
	opStoreSPMask opcode = 0xE0
	opLoadSPMask  opcode = 0xE0
	opAddSPMask   opcode = 0xF0
	opEmulateMask opcode = 0xE0
)

// State holds the state for a ZPU instance
//
type State struct {
	b      *mem.Bus
	pc     mirv.Address
	sp     mirv.Address
	idim   bool
	halted bool
}

// New instantiates a new ZPU and returns its interface.
//
func New(b *mem.Bus) cpu.Interface {
	z := State{
		b: b,
	}
	z.Reset()
	return &z
}

// Reset resets the ZPU to a known inital state.
//
func (s *State) Reset() {
	s.pc = 0
	_, e := s.b.MappedRange(mem.MemRAM)
	s.sp = e
	s.idim = false
	s.halted = false
}

// SetPC sets the PC to the given address.
//
func (s *State) SetPC(addr mirv.Address) {
	s.pc = addr
}

// PC returns the current program counter.
//
func (s *State) PC() mirv.Address {
	return s.pc
}

// SP returns the current stack pointer.
//
func (s *State) SP() mirv.Address {
	return s.sp
}

func (s *State) tos() uint32 {
	return s.read32(s.sp)
}

func (s *State) push(v uint32) {
	s.sp -= 4
	s.write32(s.sp, v)
}

func (s *State) pop() uint32 {
	v := s.read32(s.sp)
	s.sp += 4
	return v
}

func (s *State) read8(addr mirv.Address) uint8 {
	v, err := s.b.Read8(addr)
	if err != nil {
		panic(err)
	}
	return v
}

func (s *State) write32(addr mirv.Address, v uint32) {
	err := s.b.Write32BE(addr, v)
	if err != nil {
		panic(err)
	}
}

func (s State) read32(addr mirv.Address) uint32 {
	v, err := s.b.Read32BE(addr)
	if err != nil {
		panic(err)
	}
	return v
}

func (s *State) syscall() {

}

// Step steps the simulation forward n cycles. Returns how many cycles where
// performed.
//
func (s *State) Step(n uint64) uint64 {
	var cycles = n

	for ; cycles > 0 && !s.halted; cycles-- {
		var incPC = true

		if !s.idim {
			// TODO: check interupts / exceptions
		}

		insn := opcode(s.read8(s.pc))
		// TODO: check that read8 succeeded

		// Immediate
		if insn&opIMMask == opIM {
			if s.idim {
				s.write32(s.sp, (s.tos()<<7)|uint32(insn&0x7F))
			} else {
				s.push(uint32(int32(insn) << 25 >> 25))
				s.idim = true
			}
			s.pc++
			continue
		}

		// clear idim
		s.idim = false

		switch insn {
		case opBreakPoint:
			return n - cycles
		case opPopPC:
			// Pops address off stack and sets PC
			s.pc = mirv.Address(s.pop())
			incPC = false
		case opLoad:
			// Pops address stored on stack and loads the value of that address onto stack.
			addr := s.tos() & ^uint32(0x03)
			s.write32(s.sp, s.read32(mirv.Address(addr)))
		case opStore:
			// Pops address, then value from stack and stores the value into the memory location of the address.
			addr := s.pop() & ^uint32(0x03)
			s.write32(mirv.Address(addr), s.pop())
		case opPushSP:
			// Pushes stack pointer.
			s.push(uint32(s.sp))
		case opPopSP:
			// Pops value off top of stack and sets SP to that value.
			s.sp = mirv.Address(s.pop())
		case opAdd:
			// Pops two values on stack adds them and pushes the result.
			y := s.pop()
			s.write32(s.sp, s.tos()+y)
		case opAnd:
			// Pops two values off the stack and does a bitwise-and & pushes the result onto the stack
			y := s.pop()
			s.write32(s.sp, s.tos()&y)
		case opOr:
			// Pops two integers, does a bitwise or and pushes result
			y := s.pop()
			s.write32(s.sp, s.tos()|y)
		case opNot:
			// Bitwise inverse of value on stack
			s.write32(s.sp, ^s.tos())
		case opFlip:
			// Reverses the bit order of the value on the stack, i.e. abc->cba, 100->001, 110->011, etc.
			v := s.tos()
			v = (v&0x55555555)<<1 | (v>>1)&0x55555555
			v = (v&0x33333333)<<2 | (v>>2)&0x33333333
			v = (v&0x0F0F0F0F)<<4 | (v>>4)&0x0F0F0F0F
			v = (v << 24) | ((v & 0xFF00) << 8) | ((v >> 8) & 0xFF00) | (v >> 24)
			s.write32(s.sp, v)
		case opNop:

		// implementation of emulated instructions
		case opSwap:
			v := s.read32(s.sp)
			s.write32(s.sp, (v<<16)|(v>>16))
		case opSyscall:
			s.syscall()

		default:
			switch {
			// ops with embedded argument
			case insn&opStoreSPMask == opStoreSP:
				// Pop value off stack and store it in the SP+xxxxx*4 memory location, where xxxxx is a positive integer.
				// The actual memory location is xor'ed with 0x10.
				arg := mirv.Address(insn-opStoreSP) ^ 0x10
				addr := s.sp + (arg * 4)
				s.write32(addr, s.pop())
			case insn&opLoadSPMask == opLoadSP:
				// Push value of memory location SP+xxxxx*4, where xxxxx is a positive integer, onto stack.
				arg := mirv.Address(insn-opLoadSP) ^ 0x10
				addr := s.sp + (arg * 4)
				s.push(s.read32(addr))
			case insn&opAddSPMask == opAddSP:
				addr := s.sp + mirv.Address(insn-opAddSP)*4
				s.write32(s.sp, s.read32(s.sp)+s.read32(addr))
			case insn&opEmulateMask == opEmulate:
				s.push(uint32(s.pc) + 1)
				s.pc = mirv.Address(insn-opEmulate) * 32
				incPC = false
			}
		}

		if incPC {
			s.pc++
		}

	}
	return n - cycles
}
