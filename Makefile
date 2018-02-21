GO ?= go

DEPS := .deps/unicorn .deps/capstr .deps/gocui
CMD_COMMON := $(wildcard cmd/cmdline/*.go) $(wildcard cmd/util/*.go)
SRC  := $(wildcard aemulari/*.go)
BIN  := bin/aemulari bin/aemulari-gui

all: bin/aemulari bin/aemulari-gui test-asm

bin:
	@mkdir -p $@

bin/aemulari: ./cmd/aemulari/aemulari.go $(CMD_COMMON) $(SRC) $(DEPS) bin
	$(GO) build -o $@ $<

bin/aemulari-gui: ./cmd/aemulari-gui/aemulari-gui.go $(CMD_COMMON) $(SRC) $(DEPS) bin
	$(GO) build -o $@ $<

.deps:
	@mkdir -p .deps

.deps/unicorn: .deps
	$(GO) get -u github.com/unicorn-engine/unicorn/bindings/go && touch $@

.deps/capstr: .deps
	$(GO) get -u github.com/lunixbochs/capstr && touch $@

.deps/gocui: .deps
	$(GO) get -u github.com/jroimartin/gocui && touch $@

test-asm:
	$(MAKE) -C test-asm

clean:
	rm -rf bin
	$(MAKE) -C test-asm clean

realclean: clean
	rm -rf .deps

.PHONY: clean test-asm
