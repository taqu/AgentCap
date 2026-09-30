package lexer

import (
	"strings"
	"testing"
)

func render(toks []Token) []string {
	var out []string
	for _, t := range toks {
		if t.Kind == EOF {
			continue
		}
		out = append(out, t.Kind.String()+":"+t.Text)
	}
	return out
}

func TestTokenize(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []string
	}{
		{"simple", "name", []string{"Ident:name"}},
		{"number", "42", []string{"Number:42"}},
		{"decimal", "3.25", []string{"Number:3.25"}},
		{"eq", "qty = 3", []string{"Ident:qty", "Op:=", "Number:3"}},
		{"ne", "qty != 3", []string{"Ident:qty", "Op:!=", "Number:3"}},
		{"le", "qty <= 3", []string{"Ident:qty", "Op:<=", "Number:3"}},
		{"ge", "qty >= 3", []string{"Ident:qty", "Op:>=", "Number:3"}},
		{"parens", "(a)", []string{"LParen:(", "Ident:a", "RParen:)"}},
		{"plain string", "\"widget\"", []string{"String:widget"}},
		{"empty string", "\"\"", []string{"String:"}},
		{"spaces in string", "\"blue widget\"", []string{"String:blue widget"}},
		{"escaped quote", "\"say \\\"hi\\\"\"", []string{"String:say \"hi\""}},
		{"escaped backslash", "\"C:\\\\temp\"", []string{"String:C:\\temp"}},
		{"escaped newline", "\"a\\nb\"", []string{"String:a\nb"}},
		{"escaped tab", "\"a\\tb\"", []string{"String:a\tb"}},
		{"backslash then quote", "\"x\\\\\"", []string{"String:x\\"}},
		{"two strings", "\"a\\\\\" \"b\"", []string{"String:a\\", "String:b"}},
		{"and expr", "a = 1 and b = 2", []string{"Ident:a", "Op:=", "Number:1", "Ident:and", "Ident:b", "Op:=", "Number:2"}},
		{"underscore ident", "unit_price", []string{"Ident:unit_price"}},
		{"string compare", "name = \"x\"", []string{"Ident:name", "Op:=", "String:x"}},
		{"tight ops", "a<b", []string{"Ident:a", "Op:<", "Ident:b"}},
		{"nested parens", "((a))", []string{"LParen:(", "LParen:(", "Ident:a", "RParen:)", "RParen:)"}},
		{"trailing space", "a   ", []string{"Ident:a"}},
		{"tab separated", "a\t=\t1", []string{"Ident:a", "Op:=", "Number:1"}},
		{"mixed escapes", "\"q\\\"\\\\\\n\"", []string{"String:q\"\\\n"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Tokenize(tc.src)
			if err != nil {
				t.Fatalf("Tokenize(%q) error: %v", tc.src, err)
			}
			if strings.Join(render(got), "|") != strings.Join(tc.want, "|") {
				t.Fatalf("Tokenize(%q)\n got: %q\nwant: %q", tc.src, render(got), tc.want)
			}
		})
	}
}

func TestTokenizeErrors(t *testing.T) {
	for _, src := range []string{`"open`, `"bad \q"`, `a ! b`, `#`, `"trailing\`} {
		t.Run(src, func(t *testing.T) {
			if _, err := Tokenize(src); err == nil {
				t.Fatalf("Tokenize(%q) expected error", src)
			}
		})
	}
}
