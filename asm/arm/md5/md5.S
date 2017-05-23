/*******************************************************************************
 * 32-bit ARM (position-independent) implementation of MD5
 ******************************************************************************/

#include "md5.h"

/********************************************************************************
 * Top-level MD5 functions
 *******************************************************************************/

/*
 * Initialize an MD5 context structure
 *
 * Input Parameters:
 *      r0  - Pointer to MD5 context structure
 *
 * Side Effects:
 *      The MD5 context structure is ready for use.
 */
md5_init:
    stmfd   sp!, {r1, lr}

    // PC is always +8 from where we are. Also consider previous instruction.
    sub     r1, pc, #12
    str     r1, [r0, #MD5_CTX_BASE]

    ldr     r1, =A_INIT
    str     r1, [r0, #MD5_CTX_A]
    ldr     r1, =B_INIT
    str     r1, [r0, #MD5_CTX_B]
    ldr     r1, =C_INIT
    str     r1, [r0, #MD5_CTX_C]
    ldr     r1, =D_INIT
    str     r1, [r0, #MD5_CTX_D]
    mov     r1, #0x0
    str     r1, [r0, #MD5_CTX_COUNT_LOW]
    str     r1, [r0, #MD5_CTX_COUNT_HIGH]
    ldmfd   sp!, {r1, pc}

.set A_INIT, 0x67452301
.set B_INIT, 0xefcdab89
.set C_INIT, 0x98badcfe
.set D_INIT, 0x10325476

#include "block.S"
