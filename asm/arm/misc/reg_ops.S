add r0, r1, r2          // r0 = r1 + r2
ror r1, #5              // Rotate right 5 bits
sub r3, r1, r0          // r3 = r1 - r0
eor r2, r3, r1          // r2 = r3 xor r1
ldr r0, =#0xdeadbeef    // r0 = 0xdeadbeef (Note where 0xdeadbeef is stored!)
b  .                    // while (1) {}
