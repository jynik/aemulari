/*
 * Nested loop for loops with simple arithmetic and bitwise expressions
 *
 *  r0 = loop i in 0..127
 *      r1  = loop j in 0..8
 *          r2 = i + j
 *          r3 = i - j
 *          r4 = i * j
 *          r5 = i xor j
 *          r6 = i or  j
 *          r7 = i and j
 */


    // Initialize i
    mov     r0, #0

loop_i:

    // Initialize j
    mov     r1, #0

loop_j:
    // Evaluate expresions
    add     r2, r0, r1
    sub     r3, r0, r1
    mov     r4, r0
    mul     r4, r1
    mov     r5, r0
    eor     r5, r1
    mov     r6, r0
    orr     r6, r1
    mov     r7, r0
    and     r7, r1

    // Inner loop increment and test
    add     r1, #1
    cmp     r1, #7
    blt     loop_j

    // Outer loop increment and test
    add     r0, #1
    cmp     r0, #127
    blt     loop_i

done:
    bkpt    #0x0000
