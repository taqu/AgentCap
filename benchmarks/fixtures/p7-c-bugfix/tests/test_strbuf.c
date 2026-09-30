#include <string.h>
#include "collections.h"
#include "test.h"

void test_strbuf(void) {
    strbuf sb;
    sb_init(&sb);
    for (int i = 0; i < 20; i++) {
        char name[48];
        sb_appendf_int(&sb, i);
        sb_append(&sb, ",");
        snprintf(name, sizeof name, "strbuf/append %02d", i);
        CHECK(name, sb.len == strlen(sb.s));
    }
    CHECK("strbuf/content", strncmp(sb.s, "0,1,2,3,", 8) == 0);
    sb_free(&sb);
}
