	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im 0xDEADBEEF
	flip
	breakpoint
