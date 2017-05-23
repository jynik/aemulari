# Generate a .bin for each .S
SRC := $(wildcard *.S)
HDR := $(wildcard *.h)
BIN := $(SRC:.S=.bin)
OBJ := $(SRC:.S=.o)

PHONY_TARGETS += clean

%.bin: %.o
	$(OBJCOPY) $(OBJCOPYFLAGS) $< $@

%.o: %.S $(HDR) $(EXTRA_SRC) $(EXTRA_HDR)
	$(AS) $(ASFLAGS) -c $< -o $@

ifneq (${KEEP_OBJ},)
.PRECIOUS: %.o
endif

