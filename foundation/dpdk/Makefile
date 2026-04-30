CC = gcc
CFLAGS = -std=c11 -Wall -Wextra -O2
# TODO: EXTEND: Link against actual DPDK libraries when available
# DPDK_CFLAGS = $(shell pkg-config --cflags libdpdk)
# DPDK_LIBS = $(shell pkg-config --libs libdpdk)

SRCS = src/main.c src/forwarding.c
OBJS = $(SRCS:.c=.o)
TARGET = dpdk-handler

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

test:
	@echo "Unit tests for DPDK module (placeholder)"

clean:
	rm -f $(OBJS) $(TARGET)

.PHONY: all test clean
