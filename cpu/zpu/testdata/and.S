	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im 0xffff
	nop
	im 0xaaaaaaaa
	and
	breakpoint
