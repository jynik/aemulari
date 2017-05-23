/*******************************************************************************
 * Test the auxiliary function F, G, H, and I
 ******************************************************************************/

#include "start.h"

aux_tests:
    add         r7, fp, #aux_F
    add         r8, fp, #aux_F_test_00
    bl          aux_test

    add         r7, fp, #aux_G
    add         r8, fp, #aux_G_test_00
    bl          aux_test

    add         r7, fp, #aux_H
    add         r8, fp, #aux_H_test_00
    bl          aux_test

    add         r7, fp, #aux_I
    add         r8, fp, #aux_I_test_00
    bl          aux_test

    bkpt        #0x0000

/*
 * Perform test of aux fn
 *
 * Parameters:
 *      r7      Address aux function to test
 *      r8      Address of first element in array of expected results
 */
aux_test:
    stmfd       sp!, {r0-r8, lr}

    add         r5, fp, #aux_test_00        // Pointer to test case
    add         r6, fp, #aux_test_end       // Pointer to end of test case list

aux_test_loop:
    ldmfd       r5!, {r1-r3}                // Load test vector & advance pointer
    blx         r7                          // Call function under test
    ldr         r4,  [r8]                   // Load expected result
    cmp         r0, r4                      // Test for expected result
    beq         aux_test_loop_end

aux_test_fail:
    bkpt        #0xbad                      // Explode if a test case failed

aux_test_loop_end:
    add         r8, r8, #4                  // Advance result pointer
    cmp         r5, r6                      // Have we run all of our cases?
    blt         aux_test_loop               // Continue looping if not

    ldmfd       sp!, {r0-r8, pc}            // Return


/*
 * Test Vector {
 *      0x00: u32 X
 *      0x04: u32 Y
 *      0x08: u32 Z
 *      0x0C: u32 Result
 * }
 */

aux_test_00:    .word   0x00000000, 0x00000000, 0x00000000
aux_test_01:    .word   0xa55aa55a, 0x5aa55aa5, 0x8118e77e
aux_test_02:    .word   0x00000000, 0xdeadbeef, 0x00000000
aux_test_03:    .word   0x12345678, 0x9abcdef0, 0xbadc0fee
aux_test_04:    .word   0x5f787561, 0x74736574, 0x3230255f
aux_test_05:    .word   0x78302020, 0x78383025, 0x7830202c
aux_test_06:    .word   0x0000000a, 0x00000000, 0x66787561
aux_test_07:    .word   0x2e202020, 0x64726f77, 0x30202020
.set aux_test_end, .
.set aux_test_len, aux_test_02 - aux_test_01

aux_F_test_00:  .word   0x00000000
aux_F_test_01:  .word   0x00004224
aux_F_test_02:  .word   0x00000000
aux_F_test_03:  .word   0xbafc5ff6
aux_F_test_04:  .word   0x7470657e
aux_F_test_05:  .word   0x7830202c
aux_F_test_06:  .word   0x66787561
aux_F_test_07:  .word   0x34202020

aux_G_test_00:  .word   0x00000000
aux_G_test_01:  .word   0xdbbdbddb
aux_G_test_02:  .word   0xdeadbeef
aux_G_test_03:  .word   0x1234d678
aux_G_test_04:  .word   0x56736561
aux_G_test_05:  .word   0x78383021
aux_G_test_06:  .word   0x00000000
aux_G_test_07:  .word   0x64726f77

aux_H_test_00:  .word   0x00000000
aux_H_test_01:  .word   0x7ee71881
aux_H_test_02:  .word   0xdeadbeef
aux_H_test_03:  .word   0x32548766
aux_H_test_04:  .word   0x193b354a
aux_H_test_05:  .word   0x78383029
aux_H_test_06:  .word   0x6678756b
aux_H_test_07:  .word   0x7a726f77

aux_I_test_00:  .word   0xffffffff
aux_I_test_01:  .word   0xa55ae77e
aux_I_test_02:  .word   0x21524110
aux_I_test_03:  .word   0xcd8b2889
aux_I_test_04:  .word   0xab8c9a95
aux_I_test_05:  .word   0x87c7cfd6
aux_I_test_06:  .word   0x99878a9e
aux_I_test_07:  .word   0x8b8d9088

#include "../round.S"
