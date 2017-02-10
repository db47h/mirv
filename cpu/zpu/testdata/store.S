	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im 0xABCD0123
	pushsp
	im 4 
	add
	store
	breakpoint
