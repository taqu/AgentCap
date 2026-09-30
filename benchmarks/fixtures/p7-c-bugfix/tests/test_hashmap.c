#include <stdio.h>
#include "collections.h"
#include "test.h"

void test_hashmap(void) {
    hashmap hm;
    hm_init(&hm, 64);
    char key[32], name[64];
    for (int i = 0; i < 40; i++) {
        snprintf(key, sizeof key, "key-%d", i);
        hm_put(&hm, key, i * i);
    }
    for (int i = 0; i < 40; i++) {
        int v = -1;
        snprintf(key, sizeof key, "key-%d", i);
        snprintf(name, sizeof name, "hashmap/get %s", key);
        CHECK(name, hm_get(&hm, key, &v) && v == i * i);
    }
    int v;
    CHECK("hashmap/missing", !hm_get(&hm, "absent", &v));
    hm_free(&hm);
}
