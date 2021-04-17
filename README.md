# aemulari

**aemulari** is a little Go-based console UI front-end for [Unicorn],
using [Capstone], [capstr], and [gocui]. I wrote it to help me think through
32-bit ARM bare-metal code that I was otherwise reverse engineering statically.

For whatever reason, QEMU and GEF just wasn't doing it for me.

At the moment, you should consider this an open source project that's 
pretty much entirely inactive, as I haven't really touched it much since 2018,
where I presented it at [BSidesROC]. 

I originally intended to keep working on it to release something more
featureful, but alas, how time flies. I never got around to it.

Before you fork it, reach out. I may be interested in picking this back
up as I use it from time to time.

# Dependencies

* [Unicorn] - CPU Emulator Framework
    * libunicorn.so.1
    * [libunicorn Go bindings]
* [Capstone] - Disassembly Framework
    * libcapstone.so.1
    * Go bindings: [capstr]
* [gocui] - Console UI library

The provided *Makefile* will fetch and build Go dependencies. 

However, you will first need to build and install *libunicorn* and
*libcapstone* on your own.

For both of those, the build process is pretty simple:

~~~
git clone <repo>
mkdir <repo>/build && cd <repo>/build
cmake -DCMAKE_INSTALL_PATH=/usr/local
make
sudo make install
~~~


[Unicorn]: https://github.com/unicorn-engine/unicorn
[libunicorn Go bindings]: https://github.com/unicorn-engine/unicorn/tree/master/bindings/go

[Capstone]: https://github.com/aquynh/capstone
[capstr]: https://github.com/lunixbochs/capstr

[gocui]: https://github.com/jroimartin/gocui

[BSidesROC]: https://www.youtube.com/watch?v=CzHaK7cqak4

# Build

To the great displeasure of Gophers everywhere, there's a Make-based build.
(I learned Go in 2017 to write this... I was already behind the times WRT
 gomods and a proper build setup.)

Run `make` and `aemulari-cui` and `aemulari` will be build and written to a 
`bin/` directory.  The former is the UI, while the latter was intended for batch
execution.

# Known Issues

A couple of known issues are presented [here](KNOWN_ISSUES.md). 
