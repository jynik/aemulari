#ifndef START_H__
#define START_H__

#ifndef STACK_BASE
#   define STACK_BASE #0x10000
#endif

// Set up our stack and write our base address to fp
_test_start:
    mov     fp, pc      // Per ARM DDI 00029E: This actually gives PC + 8
    sub     fp, fp, #8
    mov     sp, STACK_BASE

#endif
