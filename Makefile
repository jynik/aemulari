GO ?= go

DEPS := .deps/unicorn .deps/gapstone .deps/gocui .deps/gologging

SRC  := $(wildcard internal/arch/*.go) \
		$(wildcard internal/cmdline/*.go) \
		$(wildcard internal/debugger/*.go) \
		$(wildcard internal/log/*.go) \
		$(wildcard internal/ui/*.go) \


all: aemulari aemulari-gui

aemulari: ./cmd/aemulari/aemulari.go $(SRC) $(DEPS)
	$(GO) build $<

aemulari-gui: ./cmd/aemulari-gui/aemulari-gui.go $(SRC) $(DEPS)
	$(GO) build $<

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
	make -C test-asm

clean:
	rm -f aemulari aemulari-gui
	make -C test-asm clean

realclean: clean
	rm -rf .deps

.PHONY: clean test-asm
