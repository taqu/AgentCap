#ifndef TEST_H
#define TEST_H
#include <stdio.h>
extern int failures, checks;
#define CHECK(name, cond) do { checks++; if (cond) printf("ok   %s\n", name); \
    else { failures++; printf("FAIL %s  (%s:%d: %s)\n", name, __FILE__, __LINE__, #cond); } } while (0)
void test_ringbuf(void);
void test_strbuf(void);
void test_hashmap(void);
#endif
