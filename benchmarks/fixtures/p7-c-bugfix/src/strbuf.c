#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "collections.h"

void sb_init(strbuf *sb) { sb->s = calloc(1, 16); sb->len = 0; sb->cap = 16; }
void sb_free(strbuf *sb) { free(sb->s); sb->s = NULL; }

void sb_append(strbuf *sb, const char *s) {
    size_t n = strlen(s);
    if (sb->len + n + 1 > sb->cap) {
        while (sb->len + n + 1 > sb->cap) sb->cap *= 2;
        sb->s = realloc(sb->s, sb->cap);
    }
    memcpy(sb->s + sb->len, s, n + 1);
    sb->len += n;
}

void sb_appendf_int(strbuf *sb, int v) {
    char tmp[32];
    snprintf(tmp, sizeof tmp, "%d", v);
    sb_append(sb, tmp);
}
