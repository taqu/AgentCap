#include "test.h"
int failures, checks;
int main(void) {
    test_ringbuf();
    test_strbuf();
    test_hashmap();
    printf("%d checks, %d failures\n", checks, failures);
    return failures ? 1 : 0;
}
