# aemulari

TODO: Overview

# Dependencies

* Libraries:
    * [Unicorn] - CPU Emulator Framework
        * libunicorn.so.1
        * [libunicorn Go bindings]
    * [Capstone] - Disassembly Framework
        * libcapstone.so.1
        * Go bindings: [gapstone]
            * TODO: Migrate to [capstr]
    * [gocui] - Console UI library
    * [go-logging] - Logging library
        * TODO: This is overkill and can be replaced. Logging needs to be routed to UI console anyway.

[Unicorn]: https://github.com/unicorn-engine/unicorn
[libunicorn Go bindings]: https://github.com/unicorn-engine/unicorn/tree/master/bindings/go

[Capstone]: https://github.com/aquynh/capstone
[gapstone]: https://github.com/bnagy/gapstone
[capstr]: https://github.com/lunixbochs/capstr

[gocui]: https://github.com/jroimartin/gocui

[go-logging]: https://github.com/op/go-logging

# Build

TODO

# Installation

TODO

# Examples

TODO: Links to documentation and wiki

# Known Bugs

* aemulari and aemulari-gui hang if the emulated program enters an infinite loop.
  * In this situation, execution is "stuck" in Unicorn's codebase. Perhaps we could better leverage its timeout or implement a "continue" as a series of instruction-limited executions, with opportunities for the user to interrupt and break out.
