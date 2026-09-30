#include <stdlib.h>
#include <string.h>
#include "collections.h"

static size_t hash(const char *s) {
    size_t h = 1469598103934665603ULL;
    while (*s) { h ^= (unsigned char)*s++; h *= 1099511628211ULL; }
    return h;
}

bool hm_init(hashmap *hm, size_t cap) {
    hm->keys = calloc(cap, sizeof(char *));
    hm->vals = calloc(cap, sizeof(int));
    hm->cap = cap;
    hm->len = 0;
    return hm->keys && hm->vals;
}

void hm_free(hashmap *hm) {
    for (size_t i = 0; i < hm->cap; i++) free(hm->keys[i]);
    free(hm->keys);
    free(hm->vals);
}

bool hm_put(hashmap *hm, const char *key, int val) {
    size_t i = hash(key) % hm->cap;
    for (size_t n = 0; n < hm->cap; n++, i = (i + 1) % hm->cap) {
        if (!hm->keys[i]) {
            hm->keys[i] = strdup(key);
            hm->vals[i] = val;
            hm->len++;
            return true;
        }
        if (strcmp(hm->keys[i], key) == 0) { hm->vals[i] = val; return true; }
    }
    return false;
}

bool hm_get(const hashmap *hm, const char *key, int *out) {
    size_t i = hash(key) % hm->cap;
    for (size_t n = 0; n < hm->cap && hm->keys[i]; n++, i = (i + 1) % hm->cap) {
        if (strcmp(hm->keys[i], key) == 0) { *out = hm->vals[i]; return true; }
    }
    return false;
}
