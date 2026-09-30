#include <stdio.h>
#include "collections.h"
#include "test.h"

void test_ringbuf(void) {
    ringbuf rb; int v; char name[64];
    rb_init(&rb, 4);
    CHECK("ringbuf/empty pop", !rb_pop(&rb, &v));
    CHECK("ringbuf/push 0", rb_push(&rb, 0));
    CHECK("ringbuf/push 1", rb_push(&rb, 10));
    CHECK("ringbuf/push 2", rb_push(&rb, 20));
    CHECK("ringbuf/push 3", rb_push(&rb, 30));
    CHECK("ringbuf/full rejects", !rb_push(&rb, 99));
    for (int round = 0; round < 12; round++) {
        int got = -1;
        snprintf(name, sizeof name, "ringbuf/wrap round %02d pop", round);
        CHECK(name, rb_pop(&rb, &got) && got == round * 10);
        snprintf(name, sizeof name, "ringbuf/wrap round %02d push", round);
        CHECK(name, rb_push(&rb, (round + 4) * 10));
        snprintf(name, sizeof name, "ringbuf/wrap round %02d peek oldest", round);
        CHECK(name, rb_peek(&rb, 0, &got) && got == (round + 1) * 10);
        snprintf(name, sizeof name, "ringbuf/wrap round %02d peek newest", round);
        CHECK(name, rb_peek(&rb, 3, &got) && got == (round + 4) * 10);
    }
    CHECK("ringbuf/len", rb_len(&rb) == 4);
    CHECK("ringbuf/peek out of range", !rb_peek(&rb, 4, &v));
    rb_free(&rb);
}
