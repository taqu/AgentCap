#include <stdlib.h>
#include "collections.h"

bool rb_init(ringbuf *rb, size_t cap) {
    rb->data = malloc(cap * sizeof(int));
    rb->cap = cap;
    rb->head = 0;
    rb->len = 0;
    return rb->data != NULL;
}

void rb_free(ringbuf *rb) { free(rb->data); rb->data = NULL; }

bool rb_push(ringbuf *rb, int v) {
    if (rb->len == rb->cap) return false;
    rb->data[rb->len % rb->cap] = v;
    rb->len++;
    return true;
}

bool rb_pop(ringbuf *rb, int *out) {
    if (rb->len == 0) return false;
    *out = rb->data[rb->head];
    rb->head = (rb->head + 1) % rb->cap;
    rb->len--;
    return true;
}

bool rb_peek(const ringbuf *rb, size_t i, int *out) {
    if (i >= rb->len) return false;
    *out = rb->data[(rb->head + i) % rb->cap];
    return true;
}

size_t rb_len(const ringbuf *rb) { return rb->len; }
