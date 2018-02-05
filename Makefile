GO ?= go

DEPS := .deps/unicorn .deps/gapstone .deps/gocui .deps/gologging
SRC  := $(wildcard aemulari/*.go)
BIN  := bin/aemulari bin/aemulari-gui

all: bin/aemulari bin/aemulari-gui test-asm

bin:
	@mkdir -p $@

bin/aemulari: ./cmd/aemulari/aemulari.go $(SRC) $(DEPS) bin
	$(GO) build -o $@ $<

bin/aemulari-gui: ./cmd/aemulari-gui/aemulari-gui.go $(SRC) $(DEPS) bin
	$(GO) build -o $@ $<

.deps:
	@mkdir -p .deps

.deps/unicorn: .deps
	$(GO) get -u github.com/unicorn-engine/unicorn/bindings/go
	@touch $@

.deps/gapstone: .deps
	$(GO) get -u github.com/bnagy/gapstone
	@touch $@

.deps/gocui: .deps
	$(GO) get -u github.com/jroimartin/gocui
	@touch $@

.deps/gologging: .deps
	$(GO) get -u github.com/op/go-logging
	@touch $@

test-asm:
	$(MAKE) -C test-asm

clean:
	rm -rf bin
	$(MAKE) -C test-asm clean

realclean: clean
	rm -rf .deps

.PHONY: clean test-asm
