	.org 32 
	.balign 32, 0
	.globl _start
_start:
	im _store
	load
	breakpoint

	.data
	.balign 32,0
_store:
	.int 0xABCD0123
