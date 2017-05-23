PHONY_TARGETS := test

include ../../include/bin_per_source.mk

all: test

test: $(BIN)
	$(MAKE) -C ./test

clean:
	rm -f $(BIN) $(OBJ)
	$(MAKE) -C ./test clean

.PHONY: $(PHONY_TARGETS)
