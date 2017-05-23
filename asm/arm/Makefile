export CROSSCOMPILE ?= arm-none-eabi-

export AS := $(CROSSCOMPILE)gcc
export ASFLAGS := -x assembler-with-cpp -Wall -Wextra -Werror

export OBJCOPY := $(CROSSCOMPILE)objcopy
export OBJCOPYFLAGS := -O binary

all: md5 misc

md5:
	@$(MAKE) -C md5

misc:
	@$(MAKE) -C misc

clean:
	@$(MAKE) -C md5 clean

.PHONY: md5 misc clean
