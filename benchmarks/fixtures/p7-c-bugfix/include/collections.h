#ifndef COLLECTIONS_H
#define COLLECTIONS_H
#include <stddef.h>
#include <stdbool.h>

/* Fixed-capacity FIFO ring buffer of ints. */
typedef struct {
    int *data;
    size_t cap, head, len;
} ringbuf;

bool rb_init(ringbuf *rb, size_t cap);
void rb_free(ringbuf *rb);
bool rb_push(ringbuf *rb, int v);   /* false when full */
bool rb_pop(ringbuf *rb, int *out); /* false when empty */
bool rb_peek(const ringbuf *rb, size_t i, int *out); /* i-th oldest */
size_t rb_len(const ringbuf *rb);

/* Growable string buffer. */
typedef struct { char *s; size_t len, cap; } strbuf;
void sb_init(strbuf *sb);
void sb_free(strbuf *sb);
void sb_append(strbuf *sb, const char *s);
void sb_appendf_int(strbuf *sb, int v);

/* String -> int hash map (open addressing). */
typedef struct { char **keys; int *vals; size_t cap, len; } hashmap;
bool hm_init(hashmap *hm, size_t cap);
void hm_free(hashmap *hm);
bool hm_put(hashmap *hm, const char *key, int val);
bool hm_get(const hashmap *hm, const char *key, int *out);
#endif
