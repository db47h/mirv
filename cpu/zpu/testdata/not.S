	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im 0x55555555
	not
	breakpoint
