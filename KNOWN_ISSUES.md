The following are known issues with aemulari.v0 and its associated programs:

* The debugger will hang if the emulated program enters an infinite loop.
    * In this situation, execution is "stuck" in unicorn's codebase. Perhaps we
      could better leverage its timeout or implement "continue" as a series of
      instruction-limited executions, with opportunities for the user to
      interrupt and break out.

* Execution will halt at the last instruction in the "code" memory region. 
    * This is non-ideal, as a user should be able to execute until hitting
      a breakpoint or exception, as well as execute across memory-mapped
      boundaries.
    * This is simply the nature of the unicorn API. When starting emulation,
      a start and end address are provided; the emulated machine halts at
      the latter. Grep the unicorn codebase for `addr_end` for more info.
    * Options for addressing this:
        * Submit a `CS_OPT_NO_END` patch to the unicorn maintainer.
        * Workaround: Pick the top-most unmapped address and supply this
            as `addr_end` when we emulator.
