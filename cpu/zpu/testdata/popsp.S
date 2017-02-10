	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im _newSP
	popsp
	breakpoint

	.balign 32, 0
_newSP:
	.int 0xABCD1234
