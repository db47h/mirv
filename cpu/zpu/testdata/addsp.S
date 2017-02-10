	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im -1
	nop
	im 0x34567890
	addsp 4
	breakpoint

