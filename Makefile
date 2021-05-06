GO ?= go

DEPS := .deps/unicorn .deps/capstr .deps/gocui

LIB_SRC := $(wildcard aemulari.v0/*.go)
CMD_COMMON := $(LIB_SRC) $(wildcard cmd/cmdline/*.go) $(wildcard cmd/util/*.go)

AEMULARI_SRC := $(wildcard cmd/aemulari/*.go) $(CMD_COMMON)
AEMULARI_CUI_SRC := $(wildcard cmd/aemulari-cui/*.go) \
					$(wildcard cmd/aemulari-cui/ui/*.go)

BIN  := bin/aemulari bin/aemulari-cui

INSTALL_PATH ?= /usr/local/bin

all: bin/aemulari bin/aemulari-cui

bin:
	@mkdir -p $@

bin/aemulari: $(AEMULARI_SRC) $(DEPS) bin
	$(GO) build -o $@ $<

bin/aemulari-cui: $(AEMULARI_CUI_SRC) $(DEPS) bin
	$(GO) build -o $@ $<

.deps:
	@mkdir -p .deps

.deps/unicorn: .deps
	$(GO) get -u github.com/unicorn-engine/unicorn/bindings/go && touch $@

.deps/capstr: .deps
	$(GO) get -u github.com/lunixbochs/capstr && touch $@

.deps/gocui: .deps
	$(GO) get -u github.com/jroimartin/gocui && touch $@

install:
	cp -p bin/aemulari $(INSTALL_PATH)/
	cp -p bin/aemulari-cui $(INSTALL_PATH)/

test-asm:
	$(MAKE) -C test-asm

clean:
	rm -rf bin
	$(MAKE) -C test-asm clean

realclean: clean
	rm -rf .deps

.PHONY: clean test-asm install
