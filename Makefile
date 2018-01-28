GO ?= go
DEPS := .deps/unicorn .deps/gapstone .deps/gocui .deps/gologging

all: aemulari aemulari-gui

aemulari: $(DEPS)
	$(GO) build ./cmd/aemulari/aemulari.go

aemulari-gui: $(DEPS)
	$(GO) build ./cmd/aemulari-gui/aemulari-gui.go

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

clean:
	rm -f aemulari aemulari-gui

realclean: clean
	rm -rf .deps

.PHONY: clean
