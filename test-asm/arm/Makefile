AS 		:= arm-none-eabi-gcc
ASFLAGS := -x assembler-with-cpp -Wall -Wextra -Werror

OBJCOPY := arm-none-eabi-objcopy
OBJCOPYFLAGS := -O binary

SRC := $(wildcard *.S)
HDR := $(wildcard *.h)
OBJ := $(SRC:.S=.o)
BIN := $(SRC:.S=.bin)

ifneq ($(STACK_BASE),)
	ASFLAGS += "-DSTACK_BASE=#$(STACK_BASE)"
endif

all: $(BIN)

%.o: %.S $(HDR)
	$(AS) $(ASFLAGS) -c $< -o $@

%.bin: %.o $(HDR)
	$(OBJCOPY) $(OBJCOPYFLAGS) $< $@

clean:
	rm -f *.bin

ifneq ($(KEEP_OBJ),)
.PRECIOUS: %.o
endif
