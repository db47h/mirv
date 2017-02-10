	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im -1
	nop
	im 0xABCD1234
	add
	breakpoint
