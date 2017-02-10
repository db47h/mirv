
CC = zpu-elf-gcc

.PHONY: all

all: $(patsubst %.S,%.elf,$(wildcard *.S)) hello.elf

hello.elf: hello.c
	$(CC) $< -o $@ -mboard=phi -s

%.elf: %.S
	$(CC) -nostdlib $< -o $@ -s

.PHONY: clean

clean:
	${RM} *.elf
