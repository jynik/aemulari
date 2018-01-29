AS 		:= arm-none-eabi-gcc
ASFLAGS := -x assembler-with-cpp -Wall -Wextra -Werror

OBJCOPY := arm-none-eabi-objcopy
OBJCOPYFLAGS := -O binary

HDR 		:= $(wildcard *.h)

ARM_SRC 	:= $(wildcard *.arm.S)
ARM_OBJ := $(SRC:.S=.o)

THUMB_SRC	:= $(wildcard *.thumb.S)
THUMB_OBJ	:= $(wildcard *.thumb.o)

BIN 		:= $(ARM_SRC:.S=.bin) $(THUMB_SRC:.S=.bin)

ifneq ($(STACK_BASE),)
	ASFLAGS += "-DSTACK_BASE=#$(STACK_BASE)"
endif

all: $(BIN)

%.arm.o: %.arm.S $(HDR)
	$(AS) $(ASFLAGS) -c $< -o $@

%.thumb.o: %.thumb.S $(HDR)
	$(AS) -mthumb -mcpu=Cortex-M4 $(ASFLAGS) -c $< -o $@

%.bin: %.o $(HDR)
	$(OBJCOPY) $(OBJCOPYFLAGS) $< $@

clean:
	rm -f *.bin
