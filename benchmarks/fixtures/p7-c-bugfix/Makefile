CC ?= cc
CFLAGS ?= -std=c11 -D_POSIX_C_SOURCE=200809L -Wall -Wextra -O1 -g
SRC = src/ringbuf.c src/strbuf.c src/hashmap.c
TESTS = tests/test_main.c tests/test_ringbuf.c tests/test_strbuf.c tests/test_hashmap.c

all: libcollections.a

libcollections.a: $(SRC:.c=.o)
	ar rcs $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -Iinclude -c $< -o $@

run_tests: $(SRC) $(TESTS)
	$(CC) $(CFLAGS) -Iinclude $(SRC) $(TESTS) -o $@

test: run_tests
	./run_tests

clean:
	rm -f src/*.o libcollections.a run_tests

.PHONY: all test clean
