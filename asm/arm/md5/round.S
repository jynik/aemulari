/*******************************************************************************
 * MD5 round operations and auxiliary functions
 ******************************************************************************/

/*
 * Perform auxiliary function F(X, Y, Z) = XY v not(X)Z
 *
 * Input Parameters:
 *  r1 - X
 *  r2 - Y
 *  r3 - Z
 *
 * Output Parameters:
 *  r0 - Result
 */
aux_F:
    stmfd       sp!, {r1-r3, lr}
    and         r2, r1, r2      // r2 <- XY       = r1 and r2
    bic         r3, r3, r1      // r3 <- not(X)Z  = r3 and NOT(r1)
    orr         r0, r2, r3
    ldmfd       sp!, {r1-r3, pc}

/*
 * Perform auxiliary function G(X, Y, Z) = XZ v Y not(Z)
 *
 * Input Parameters:
 *  r1 - X
 *  r2 - Y
 *  r3 - Z
 *
 * Output Parameters:
 *  r0 - Result
 */
aux_G:
    stmfd       sp!, {r1-r3, lr}
    and         r1, r1, r3  // r1 <- XZ
    bic         r2, r2, r3  // r2 <- Y not(Z)
    orr         r0, r1, r2
    ldmfd       sp!, {r1-r3, pc}

/*
 * Perform auxiliary function H(X, Y, Z) = X xor Y xor Z
 *
 * Input Parameters:
 *  r1 - X
 *  r2 - Y
 *  r3 - Z
 *
 * Output Parameters:
 *  r0 - Result
 */
aux_H:
    stmfd       sp!, {r1-r3, lr}
    eor         r1, r1, r2  // r1 <- X xor Y
    eor         r0, r1, r3  // r0 <- (X xor Y) xor Z
    ldmfd       sp!, {r1-r3, pc}

/*
 * Perform auxiliary function I(X, Y, Z) = Y xor (X v not(Z))
 *
 * Input Parameters:
 *  r1 - X
 *  r2 - Y
 *  r3 - Z
 *
 * Output Parameters:
 *  r0 - Result
 */
aux_I:
    stmfd       sp!, {r1-r3, lr}
    mvn         r3, r3      // r3 <= NOT(Z)     TODO try Thumb-2 orn instr?
    orr         r3, r1, r3  // r2 <- X v NOT(Z)
    eor         r0, r2, r3  // r0 <- Y xor (X v NOT(Z))
    ldmfd       sp!, {r1-r3, pc}

/*
 * Helper implementation of the round operation "[abcd k s i]"
 *
 * Input Parameters:
 *  r0 - a
 *  r1 - b
 *  r2 - c
 *  r3 - d
 *  r4 - X[k]
 *  r5 - s
 *  r6 - t[i]
 *  r7 - Round-specific auxiliary function (F, G, H, or I)
 *
 * Output Parameters:
 *  r0 - Updated `a` value
 */
round_op:
    stmfd   sp!, {r1-r8, lr}
    mov     r8, r0          // Back up `a`
    blx     r7              // r0 = aux_fn(b, c, d)
    add     r0, r0, r8      // r0 += a
    add     r0, r0, r4      // r0 += X[k]
    add     r0, r0, r6      // r0 += t[i]
    mov     r6, #32         // s = 32 - s   ROL(x, n) == ROR(x, 32-n)
    sub     r5, r6, r5      // r0 <<< s
    ror     r0, r0, r5
    add     r0, r0, r1      // r0 += b
    ldmfd   sp!, {r1-r8, pc}
