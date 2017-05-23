/*******************************************************************************
 * Test the round tranformation functions FF, GG, HH, II
 * (These are defined in the reference C code)
 ******************************************************************************/

#include "start.h"

/*
 * Macro for invoking round operation with the specified auxiliary function
 */
.macro perform_round_op aux_fn
    stmfd       sp!, {r7, lr}
    add         r7, fp, #\aux_fn
    bl          round_op
    ldmfd       sp!, {r7, pc}
.endm

round_xfrm_tests:
    add         r8, fp, #FF_test_01
    add         r9, fp, #FF_test
    bl          xfrm_test

    add         r8, fp, #GG_test_01
    add         r9, fp, #GG_test
    bl          xfrm_test

    add         r8, fp, #HH_test_01
    add         r9, fp, #HH_test
    bl          xfrm_test

    add         r8, fp, #II_test_01
    add         r9, fp, #II_test
    bl          xfrm_test

    bkpt        #0x0000

/*
 * Perform test of FF, GG, HH, or II transform functions
 *
 * Input Parameters:
 *      r8      Address of first element in array of expected results
 *      r9      Address of xfrm function to test
 *
 * Side Effects:
 *      r8 and r9 are modfied
 */
xfrm_test:
    stmfd       sp!, {r0-r10, lr}

    add         r7,  fp, #xfrm_test_01
    add         r10, fp, #xfrm_test_end

xfrm_test_loop:
    ldmfd       r7!, {r0-r6}
    blx         r9
    ldr         r1, [r8]
    cmp         r0, r1
    beq         xfrm_test_loop_end

xfrm_test_fail:
    bkpt        #0xbad

xfrm_test_loop_end:
    add         r8, r8, #4
    cmp         r7, r10
    blt         xfrm_test_loop

    ldmfd       sp!, {r0-r10, pc}

xfrm_test_01:   .word   0x00000000, 0x00000000, 0x00000000, 0x00000000, 0x00,  0, 0x00000000
xfrm_test_02:   .word   0x00000001, 0x00000100, 0x00010000, 0x01000000, 0x01,  1, 0xa55a5aa5
xfrm_test_03:   .word   0x11111111, 0x22222222, 0x44444444, 0x88888888, 0xcc, 27, 0xaaaaaaaa
xfrm_test_04:   .word   0xffffffff, 0xffffffff, 0xffffffff, 0xffffffff, 0xff, 15, 0xffffffff
xfrm_test_05:   .word   0xdeadbeef, 0xbadc0ffe, 0x12345687, 0x0abcdef0, 0xb3,  7, 0x2f932e10
xfrm_test_06:   .word   0x120f23ed, 0xfe013d13, 0x139c831d, 0x3291d203, 0x38,  6, 0x1973c1a0
.set xfrm_test_end,  .
.set xfrm_test_len,  xfrm_test_06 - xfrm_test_01

FF_test_01:     .word   0x00000000
FF_test_02:     .word   0x4cb4b64f
FF_test_03:     .word   0x9c44444a
FF_test_04:     .word   0x007dffff
FF_test_05:     .word   0xf5be2c0e
FF_test_06:     .word   0x82eb72a2

GG_test_01:     .word   0x00000000
GG_test_02:     .word   0x4ab6b64f
GG_test_03:     .word   0x7a222228
GG_test_04:     .word   0x007dffff
GG_test_05:     .word   0x295a6492
GG_test_06:     .word   0xa1fef62a

HH_test_01:     .word   0x00000000
HH_test_02:     .word   0x4cb6b84f
HH_test_03:     .word   0xcf77777d
HH_test_04:     .word   0x007dffff
HH_test_05:     .word   0x0596add6
HH_test_06:     .word   0xa1d5b195

II_test_01:     .word   0xffffffff
II_test_02:     .word   0x48b2b64d
II_test_03:     .word   0xf999999f
II_test_04:     .word   0x007e7fff
II_test_05:     .word   0xd10fa57c
II_test_06:     .word   0x1baa66d9

#include "../round.S"

/* Round operation tests */
FF_test:
    perform_round_op aux_F

GG_test:
    perform_round_op aux_G

HH_test:
    perform_round_op aux_H

II_test:
    perform_round_op aux_I

