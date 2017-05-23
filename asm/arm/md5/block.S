/*******************************************************************************
 * Implementation of per-block operation
 ******************************************************************************/

#include "md5.h"
#include "round.S"

/* Helper macro for performing round operation */
.macro  call_round_op aux_fn a b c d k s t
    ldr     r0, [r8, #MD5_CTX_\a]
    ldr     r1, [r8, #MD5_CTX_\b]
    ldr     r2, [r8, #MD5_CTX_\c]
    ldr     r3, [r8, #MD5_CTX_\d]
    ldr     r4, [r8, #(\k * 4)]
    ldr     r5, =#\s
    ldr     r6, =#\t
    ldr     r7, [r8, #MD5_CTX_BASE]
    add     r7, #aux_\aux_fn
    bl      round_op
.endm

/* The following macros implement per-round operations.
 * These are broken out to allow them to be pulled into some unit test programs
 */

.macro perform_round_1
    call_round_op F A B C D I1_01 S11 T1_01
    call_round_op F D A B C I1_02 S12 T1_02
    call_round_op F C D A B I1_03 S13 T1_03
    call_round_op F B C D A I1_04 S14 T1_04

    call_round_op F A B C D I1_05 S11 T1_05
    call_round_op F D A B C I1_06 S12 T1_06
    call_round_op F C D A B I1_07 S13 T1_07
    call_round_op F B C D A I1_08 S14 T1_08

    call_round_op F A B C D I1_09 S11 T1_09
    call_round_op F D A B C I1_10 S12 T1_10
    call_round_op F C D A B I1_11 S13 T1_11
    call_round_op F B C D A I1_12 S14 T1_12

    call_round_op F A B C D I1_13 S11 T1_13
    call_round_op F D A B C I1_14 S12 T1_14
    call_round_op F C D A B I1_15 S13 T1_15
    call_round_op F B C D A I1_16 S14 T1_16
.endm

.macro perform_round_2
    call_round_op G A B C D I2_01 S21 T2_01
    call_round_op G D A B C I2_02 S22 T2_02
    call_round_op G C D A B I2_03 S23 T2_03
    call_round_op G B C D A I2_04 S24 T2_04

    call_round_op G A B C D I2_05 S21 T2_05
    call_round_op G D A B C I2_06 S22 T2_06
    call_round_op G C D A B I2_07 S23 T2_07
    call_round_op G B C D A I2_08 S24 T2_08

    call_round_op G A B C D I2_09 S21 T2_09
    call_round_op G D A B C I2_10 S22 T2_10
    call_round_op G C D A B I2_11 S23 T2_11
    call_round_op G B C D A I2_12 S24 T2_12

    call_round_op G A B C D I2_13 S21 T2_13
    call_round_op G D A B C I2_14 S22 T2_14
    call_round_op G C D A B I2_15 S23 T2_15
    call_round_op G B C D A I2_16 S24 T2_16
.endm

.macro perform_round_3
    call_round_op H A B C D I3_01 S31 T3_01
    call_round_op H D A B C I3_02 S32 T3_02
    call_round_op H C D A B I3_03 S33 T3_03
    call_round_op H B C D A I3_04 S34 T3_04

    call_round_op H A B C D I3_05 S31 T3_05
    call_round_op H D A B C I3_06 S32 T3_06
    call_round_op H C D A B I3_07 S33 T3_07
    call_round_op H B C D A I3_08 S34 T3_08

    call_round_op H A B C D I3_09 S31 T3_09
    call_round_op H D A B C I3_10 S32 T3_10
    call_round_op H C D A B I3_11 S33 T3_11
    call_round_op H B C D A I3_12 S34 T3_12

    call_round_op H A B C D I3_13 S31 T3_13
    call_round_op H D A B C I3_14 S32 T3_14
    call_round_op H C D A B I3_15 S33 T3_15
    call_round_op H B C D A I3_16 S34 T3_16
.endm

.macro perform_round_4
    call_round_op I A B C D I4_01 S41 T4_01
    call_round_op I D A B C I4_02 S42 T4_02
    call_round_op I C D A B I4_03 S43 T4_03
    call_round_op I B C D A I4_04 S44 T4_04

    call_round_op I A B C D I4_05 S41 T4_05
    call_round_op I D A B C I4_06 S42 T4_06
    call_round_op I C D A B I4_07 S43 T4_07
    call_round_op I B C D A I4_08 S44 T4_08

    call_round_op I A B C D I4_09 S41 T4_09
    call_round_op I D A B C I4_10 S42 T4_10
    call_round_op I C D A B I4_11 S43 T4_11
    call_round_op I B C D A I4_12 S44 T4_12

    call_round_op I A B C D I4_13 S41 T4_13
    call_round_op I D A B C I4_14 S42 T4_14
    call_round_op I C D A B I4_15 S43 T4_15
    call_round_op I B C D A I4_16 S44 T4_16
.endm

/*******************************************************************************
 * Round constant definitions
 *
 * I<x>_<y>     X[i] index for round <x>, operation <y>
 * S<x><y>      Rotation value (<<< s)
 * T<x>_<y>     T[i] value for round <x>, operation <y>
 ******************************************************************************/

// Round 1
.set I1_01, 0
.set I1_02, 1
.set I1_03, 2
.set I1_04, 3
.set I1_05, 4
.set I1_06, 5
.set I1_07, 6
.set I1_08, 7
.set I1_09, 8
.set I1_10, 9
.set I1_11, 10
.set I1_12, 11
.set I1_13, 12
.set I1_14, 13
.set I1_15, 14
.set I1_16, 15

.set S11,   7
.set S12,   12
.set S13,   17
.set S14,   22

.set T1_01, 0xd76aa478
.set T1_02, 0xe8c7b756
.set T1_03, 0x242070db
.set T1_04, 0xc1bdceee
.set T1_05, 0xf57c0faf
.set T1_06, 0x4787c62a
.set T1_07, 0xa8304613
.set T1_08, 0xfd469501
.set T1_09, 0x698098d8
.set T1_10, 0x8b44f7af
.set T1_11, 0xffff5bb1
.set T1_12, 0x895cd7be
.set T1_13, 0x6b901122
.set T1_14, 0xfd987193
.set T1_15, 0xa679438e
.set T1_16, 0x49b40821

// Round 2
.set I2_01, 1
.set I2_02, 6
.set I2_03, 11
.set I2_04, 0
.set I2_05, 5
.set I2_06, 10
.set I2_07, 15
.set I2_08, 4
.set I2_09, 9
.set I2_10, 14
.set I2_11, 3
.set I2_12, 8
.set I2_13, 13
.set I2_14, 2
.set I2_15, 7
.set I2_16, 12

.set S21,   5
.set S22,   9
.set S23,   14
.set S24,   20

.set T2_01, 0xf61e2562
.set T2_02, 0xc040b340
.set T2_03, 0x265e5a51
.set T2_04, 0xe9b6c7aa
.set T2_05, 0xd62f105d
.set T2_06, 0x02441453
.set T2_07, 0xd8a1e681
.set T2_08, 0xe7d3fbc8
.set T2_09, 0x21e1cde6
.set T2_10, 0xc33707d6
.set T2_11, 0xf4d50d87
.set T2_12, 0x455a14ed
.set T2_13, 0xa9e3e905
.set T2_14, 0xfcefa3f8
.set T2_15, 0x676f02d9
.set T2_16, 0x8d2a4c8a

// Round 3
.set I3_01, 5
.set I3_02, 8
.set I3_03, 11
.set I3_04, 14
.set I3_05, 1
.set I3_06, 4
.set I3_07, 7
.set I3_08, 10
.set I3_09, 13
.set I3_10, 0
.set I3_11, 3
.set I3_12, 6
.set I3_13, 9
.set I3_14, 12
.set I3_15, 15
.set I3_16, 2

.set S31,   4
.set S32,   11
.set S33,   16
.set S34,   23

.set T3_01, 0xfffa3942
.set T3_02, 0x8771f681
.set T3_03, 0x6d9d6122
.set T3_04, 0xfde5380c
.set T3_05, 0xa4beea44
.set T3_06, 0x4bdecfa9
.set T3_07, 0xf6bb4b60
.set T3_08, 0xbebfbc70
.set T3_09, 0x289b7ec6
.set T3_10, 0xeaa127fa
.set T3_11, 0xd4ef3085
.set T3_12, 0x04881d05
.set T3_13, 0xd9d4d039
.set T3_14, 0xe6db99e5
.set T3_15, 0x1fa27cf8
.set T3_16, 0xc4ac5665

// Round 4
.set I4_01, 0
.set I4_02, 7
.set I4_03, 14
.set I4_04, 5
.set I4_05, 12
.set I4_06, 3
.set I4_07, 10
.set I4_08, 1
.set I4_09, 8
.set I4_10, 15
.set I4_11, 6
.set I4_12, 13
.set I4_13, 4
.set I4_14, 11
.set I4_15, 2
.set I4_16, 9

.set S41,   6
.set S42,   10
.set S43,   15
.set S44,   21

.set T4_01, 0xf4292244
.set T4_02, 0x432aff97
.set T4_03, 0xab9423a7
.set T4_04, 0xfc93a039
.set T4_05, 0x655b59c3
.set T4_06, 0x8f0ccc92
.set T4_07, 0xffeff47d
.set T4_08, 0x85845dd1
.set T4_09, 0x6fa87e4f
.set T4_10, 0xfe2ce6e0
.set T4_11, 0xa3014314
.set T4_12, 0x4e0811a1
.set T4_13, 0xf7537e82
.set T4_14, 0xbd3af235
.set T4_15, 0x2ad7d2bb
.set T4_16, 0xeb86d391


/*
 * Process a 64-byte (16-word) block
 *
 * Input Parameters:
 *  r0  -   Pointer to digest context
 *  r1  -   Pointer to the 64-byte input buffer
 *
 * Output Parameters:
 *  r0  -   The digest context's internal state will be updated
 *
 */
process_block:
    stmfd   sp!, {r0-r8, lr}

    // Move inputs to other registers before reusing them as fn args
    mov     r7, r1  // r7 = X[]
    mov     r8, r0  // r8 = md5_ctx

    perform_round_1
    perform_round_2
    perform_round_3
    perform_round_4

    // Store state in context structure
    ldr     r5, [r8, #MD5_CTX_A]
    add     r5, r0
    str     r5, [r8, #MD5_CTX_A]

    ldr     r5, [r8, #MD5_CTX_B]
    add     r5, r1
    str     r5, [r8, #MD5_CTX_B]

    ldr     r5, [r8, #MD5_CTX_C]
    add     r5, r2
    str     r5, [r8, #MD5_CTX_C]

    ldr     r5, [r8, #MD5_CTX_D]
    add     r5, r3
    str     r5, [r8, #MD5_CTX_D]

    ldmfd   sp!, {r0-r8, pc}

.byte 0x00, 0x02, 0x03, \
      0x04, 0x05
