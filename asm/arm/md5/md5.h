#ifndef __MD5_H__
#define __MD5_H__

/*
 *  MD5 Context {
 *      0x00: u32 A
 *      0x04: u32 B
 *      0x08: u32 C
 *      0x0C: u32 D
 *      0x10: u32 count_low
 *      0x14: u32 count_high
 *      0x18: u32 base_addr of MD5 code
 *  }
 */

.set MD5_CTX_A,             0x00
.set MD5_CTX_B,             0x04
.set MD5_CTX_C,             0x08
.set MD5_CTX_D,             0x0C
.set MD5_CTX_COUNT_LOW,     0x10
.set MD5_CTX_COUNT_HIGH,    0x14
.set MD5_CTX_BASE,          0x18
.set MD5_CTX_SIZE,          0x1c

#endif
