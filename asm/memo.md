```
gcc -nostdlib -fno-asynchronous-unwind-tables -I./include -g -fno-stack-protector -m32 -fno-pie -c test2.c
nasm -f elf crt0.asm
ld --entry=start --oformat=binary -Ttext 0x7c00 -m elf_i386 -o test.bin crt0.o test.o
```
