package query

import "testing"

func TestParse(t *testing.T) {
	cases := []struct{ name, src, want string }{
		{"single", "qty > 3", "(qty > 3)"},
		{"string", "name = \"bolt\"", "(name = \"bolt\")"},
		{"and", "a = 1 and b = 2", "((a = 1) and (b = 2))"},
		{"or", "a = 1 or b = 2", "((a = 1) or (b = 2))"},
		{"precedence", "a = 1 or b = 2 and c = 3", "((a = 1) or ((b = 2) and (c = 3)))"},
		{"parens", "(a = 1 or b = 2) and c = 3", "(((a = 1) or (b = 2)) and (c = 3))"},
		{"escaped quote", "name = \"12\\\" pipe\"", "(name = \"12\\\" pipe\")"},
		{"escaped backslash", "path = \"a\\\\b\"", "(path = \"a\\\\b\")"},
		{"backslash end", "path = \"dir\\\\\" and qty = 1", "((path = \"dir\\\\\") and (qty = 1))"},
		{"decimal", "price <= 9.5", "(price <= 9.5)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := Parse(tc.src)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.src, err)
			}
			if e.String() != tc.want {
				t.Fatalf("Parse(%q)\n got: %s\nwant: %s", tc.src, e, tc.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	for _, src := range []string{"a =", "= 1", "(a = 1", "a = 1 b", `a = "x`} {
		t.Run(src, func(t *testing.T) {
			if _, err := Parse(src); err == nil {
				t.Fatalf("Parse(%q) expected error", src)
			}
		})
	}
}
