# Position-Independent ARM MD5 Implementation

This directory contains a position-independent ARM MD5 implementation.
It is designed to be easy to read when reviewed alongside RFC 1321, rather
than for performance. 

This code was developed in order to test and debug the aemulari GUI and batch
execution programs (`aemulari-gui` and `aemulari`, respectively). It is 
intended for use in slightly more complicated documentation examples.

## Usage

TODO: Describe how to write a test harness that calls functions in `md5.bin`.

## Source Files

|------------------|------------------------------------------------|
| Filename         | Description                                    |
|------------------|------------------------------------------------|
| `md5.S`          | Top-level "API" for the implementation         |
| `md5.h`          | Constant definitions and macros                |
| `block.S`        | Per-block operations                           |
| `round.S         | Per-round operations and auxiliary functions   |
| `test/start.h`   | Test setup boilerplate                         |
| `test/aux.S`     | Unit test for per-round "auxiliary" functions  |
| `test/round.S    | Unit test for per-round operations             |
| `test/block.S    | Unit test for per-block operations             |
| `test/md5.S      | Top-level integration test                     |
