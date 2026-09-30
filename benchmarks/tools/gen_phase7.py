#!/usr/bin/env python3
"""Deterministic generator for Phase 7 coding-agent benchmark fixtures.

Every fixture is derived from a known-correct base tree; each task injects one
defect. The base tree therefore doubles as the reference solution:
`gen_phase7.py --check` verifies that every base passes its verifier and every
task fixture fails it before any agent runs (harness validation, Phase 7A).

Usage:
    python3 benchmarks/tools/gen_phase7.py           # (re)write fixtures
    python3 benchmarks/tools/gen_phase7.py --check   # also validate oracles
"""
import json
import os
import random
import shutil
import subprocess
import sys
import tempfile
import textwrap

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
FIXTURES = os.path.join(ROOT, "fixtures")


def d(s):
    return textwrap.dedent(s).lstrip("\n")


# --------------------------------------------------------------- Go base tree

GO = {}
GO["go.mod"] = "module inventory\n\ngo 1.22\n"

GO["lexer/lexer.go"] = d('''
    // Package lexer tokenizes inventory query expressions such as
    //   name = "blue \\"widget\\"" and qty > 10
    package lexer

    import (
    	"fmt"
    	"strings"
    	"unicode"
    )

    type Kind int

    const (
    	EOF Kind = iota
    	Ident
    	Number
    	String
    	Op
    	LParen
    	RParen
    )

    func (k Kind) String() string {
    	return [...]string{"EOF", "Ident", "Number", "String", "Op", "LParen", "RParen"}[k]
    }

    type Token struct {
    	Kind Kind
    	Text string
    	Pos  int
    }

    // Tokenize splits src into tokens. String literals are double-quoted and
    // support the escapes \\\\ \\" \\n and \\t.
    func Tokenize(src string) ([]Token, error) {
    	var out []Token
    	i := 0
    	for i < len(src) {
    		c := rune(src[i])
    		switch {
    		case unicode.IsSpace(c):
    			i++
    		case c == '(':
    			out = append(out, Token{LParen, "(", i})
    			i++
    		case c == ')':
    			out = append(out, Token{RParen, ")", i})
    			i++
    		case c == '"':
    			text, n, err := readString(src[i:])
    			if err != nil {
    				return nil, fmt.Errorf("pos %d: %w", i, err)
    			}
    			out = append(out, Token{String, text, i})
    			i += n
    		case strings.ContainsRune("=!<>", c):
    			j := i + 1
    			if j < len(src) && src[j] == '=' {
    				j++
    			}
    			op := src[i:j]
    			if op == "!" {
    				return nil, fmt.Errorf("pos %d: unexpected '!'", i)
    			}
    			out = append(out, Token{Op, op, i})
    			i = j
    		case unicode.IsDigit(c):
    			j := i
    			for j < len(src) && (unicode.IsDigit(rune(src[j])) || src[j] == '.') {
    				j++
    			}
    			out = append(out, Token{Number, src[i:j], i})
    			i = j
    		case unicode.IsLetter(c) || c == '_':
    			j := i
    			for j < len(src) && (unicode.IsLetter(rune(src[j])) || unicode.IsDigit(rune(src[j])) || src[j] == '_') {
    				j++
    			}
    			out = append(out, Token{Ident, src[i:j], i})
    			i = j
    		default:
    			return nil, fmt.Errorf("pos %d: unexpected %q", i, c)
    		}
    	}
    	out = append(out, Token{EOF, "", len(src)})
    	return out, nil
    }

    // readString reads a quoted literal starting at s[0]=='"' and returns the
    // unescaped text and the number of bytes consumed.
    func readString(s string) (string, int, error) {
    	var b strings.Builder
    	for i := 1; i < len(s); i++ {
    		switch s[i] {
    		case '"':
    			return b.String(), i + 1, nil
    		case '\\\\':
    			if i+1 >= len(s) {
    				return "", 0, fmt.Errorf("unterminated escape")
    			}
    			i++
    			switch s[i] {
    			case 'n':
    				b.WriteByte('\\n')
    			case 't':
    				b.WriteByte('\\t')
    			case '"', '\\\\':
    				b.WriteByte(s[i])
    			default:
    				return "", 0, fmt.Errorf("unknown escape \\\\%c", s[i])
    			}
    		default:
    			b.WriteByte(s[i])
    		}
    	}
    	return "", 0, fmt.Errorf("unterminated string")
    }
''')

LEX_CASES = [
    ("simple", 'name', ["Ident:name"]),
    ("number", '42', ["Number:42"]),
    ("decimal", '3.25', ["Number:3.25"]),
    ("eq", 'qty = 3', ["Ident:qty", "Op:=", "Number:3"]),
    ("ne", 'qty != 3', ["Ident:qty", "Op:!=", "Number:3"]),
    ("le", 'qty <= 3', ["Ident:qty", "Op:<=", "Number:3"]),
    ("ge", 'qty >= 3', ["Ident:qty", "Op:>=", "Number:3"]),
    ("parens", '(a)', ["LParen:(", "Ident:a", "RParen:)"]),
    ("plain string", '"widget"', ["String:widget"]),
    ("empty string", '""', ["String:"]),
    ("spaces in string", '"blue widget"', ["String:blue widget"]),
    ("escaped quote", r'"say \"hi\""', ['String:say "hi"']),
    ("escaped backslash", r'"C:\\temp"', [r"String:C:\temp"]),
    ("escaped newline", r'"a\nb"', ["String:a\nb"]),
    ("escaped tab", r'"a\tb"', ["String:a\tb"]),
    ("backslash then quote", r'"x\\"', [r"String:x\\"[:-1]]),
    ("two strings", r'"a\\" "b"', [r"String:a\\"[:-1], "String:b"]),
    ("and expr", 'a = 1 and b = 2', ["Ident:a", "Op:=", "Number:1", "Ident:and", "Ident:b", "Op:=", "Number:2"]),
    ("underscore ident", 'unit_price', ["Ident:unit_price"]),
    ("string compare", 'name = "x"', ["Ident:name", "Op:=", "String:x"]),
    ("tight ops", 'a<b', ["Ident:a", "Op:<", "Ident:b"]),
    ("nested parens", '((a))', ["LParen:(", "LParen:(", "Ident:a", "RParen:)", "RParen:)"]),
    ("trailing space", 'a   ', ["Ident:a"]),
    ("tab separated", 'a\t=\t1', ["Ident:a", "Op:=", "Number:1"]),
    ("mixed escapes", r'"q\"\\\n"', ['String:q"\\\n']),
]


def go_quote(s):
    out = '"'
    for ch in s:
        if ch == '\\':
            out += '\\\\'
        elif ch == '"':
            out += '\\"'
        elif ch == '\n':
            out += '\\n'
        elif ch == '\t':
            out += '\\t'
        else:
            out += ch
    return out + '"'


def lexer_test():
    rows = []
    for name, src, toks in LEX_CASES:
        rows.append("\t\t{%s, %s, []string{%s}}," % (go_quote(name), go_quote(src), ", ".join(go_quote(t) for t in toks)))
    return d('''
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
    ''') + "\n".join(rows) + "\n" + d('''
    	}
    	for _, tc := range cases {
    		t.Run(tc.name, func(t *testing.T) {
    			got, err := Tokenize(tc.src)
    			if err != nil {
    				t.Fatalf("Tokenize(%q) error: %v", tc.src, err)
    			}
    			if strings.Join(render(got), "|") != strings.Join(tc.want, "|") {
    				t.Fatalf("Tokenize(%q)\\n got: %q\\nwant: %q", tc.src, render(got), tc.want)
    			}
    		})
    	}
    }

    func TestTokenizeErrors(t *testing.T) {
    	for _, src := range []string{`"open`, `"bad \\q"`, `a ! b`, `#`, `"trailing\\`} {
    		t.Run(src, func(t *testing.T) {
    			if _, err := Tokenize(src); err == nil {
    				t.Fatalf("Tokenize(%q) expected error", src)
    			}
    		})
    	}
    }
    ''')


GO["lexer/lexer_test.go"] = lexer_test()

GO["query/ast.go"] = d('''
    package query

    import "fmt"

    // Expr is a boolean expression over item fields.
    type Expr interface{ String() string }

    type Compare struct {
    	Field string
    	Op    string
    	Value Value
    }

    type Value struct {
    	IsString bool
    	Str      string
    	Num      float64
    }

    type And struct{ Left, Right Expr }
    type Or struct{ Left, Right Expr }

    func (c Compare) String() string {
    	if c.Value.IsString {
    		return fmt.Sprintf("(%s %s %q)", c.Field, c.Op, c.Value.Str)
    	}
    	return fmt.Sprintf("(%s %s %g)", c.Field, c.Op, c.Value.Num)
    }
    func (a And) String() string { return "(" + a.Left.String() + " and " + a.Right.String() + ")" }
    func (o Or) String() string  { return "(" + o.Left.String() + " or " + o.Right.String() + ")" }
''')

GO["query/parser.go"] = d('''
    // Package query parses inventory filter expressions.
    package query

    import (
    	"fmt"
    	"strconv"

    	"inventory/lexer"
    )

    type parser struct {
    	toks []lexer.Token
    	pos  int
    }

    // Parse parses src into an expression tree. "and" binds tighter than "or".
    func Parse(src string) (Expr, error) {
    	toks, err := lexer.Tokenize(src)
    	if err != nil {
    		return nil, err
    	}
    	p := &parser{toks: toks}
    	e, err := p.or()
    	if err != nil {
    		return nil, err
    	}
    	if p.peek().Kind != lexer.EOF {
    		return nil, fmt.Errorf("unexpected %q at %d", p.peek().Text, p.peek().Pos)
    	}
    	return e, nil
    }

    func (p *parser) peek() lexer.Token { return p.toks[p.pos] }
    func (p *parser) next() lexer.Token  { t := p.toks[p.pos]; p.pos++; return t }

    func (p *parser) or() (Expr, error) {
    	left, err := p.and()
    	if err != nil {
    		return nil, err
    	}
    	for p.peek().Kind == lexer.Ident && p.peek().Text == "or" {
    		p.next()
    		right, err := p.and()
    		if err != nil {
    			return nil, err
    		}
    		left = Or{left, right}
    	}
    	return left, nil
    }

    func (p *parser) and() (Expr, error) {
    	left, err := p.primary()
    	if err != nil {
    		return nil, err
    	}
    	for p.peek().Kind == lexer.Ident && p.peek().Text == "and" {
    		p.next()
    		right, err := p.primary()
    		if err != nil {
    			return nil, err
    		}
    		left = And{left, right}
    	}
    	return left, nil
    }

    func (p *parser) primary() (Expr, error) {
    	if p.peek().Kind == lexer.LParen {
    		p.next()
    		e, err := p.or()
    		if err != nil {
    			return nil, err
    		}
    		if p.next().Kind != lexer.RParen {
    			return nil, fmt.Errorf("expected )")
    		}
    		return e, nil
    	}
    	field := p.next()
    	if field.Kind != lexer.Ident {
    		return nil, fmt.Errorf("expected field name at %d", field.Pos)
    	}
    	op := p.next()
    	if op.Kind != lexer.Op {
    		return nil, fmt.Errorf("expected operator at %d", op.Pos)
    	}
    	val := p.next()
    	switch val.Kind {
    	case lexer.String:
    		return Compare{field.Text, op.Text, Value{IsString: true, Str: val.Text}}, nil
    	case lexer.Number:
    		n, err := strconv.ParseFloat(val.Text, 64)
    		if err != nil {
    			return nil, err
    		}
    		return Compare{field.Text, op.Text, Value{Num: n}}, nil
    	}
    	return nil, fmt.Errorf("expected value at %d", val.Pos)
    }
''')

QUERY_CASES = [
    ("single", 'qty > 3', '(qty > 3)'),
    ("string", 'name = "bolt"', '(name = "bolt")'),
    ("and", 'a = 1 and b = 2', '((a = 1) and (b = 2))'),
    ("or", 'a = 1 or b = 2', '((a = 1) or (b = 2))'),
    ("precedence", 'a = 1 or b = 2 and c = 3', '((a = 1) or ((b = 2) and (c = 3)))'),
    ("parens", '(a = 1 or b = 2) and c = 3', '(((a = 1) or (b = 2)) and (c = 3))'),
    ("escaped quote", r'name = "12\" pipe"', r'(name = "12\" pipe")'),
    ("escaped backslash", r'path = "a\\b"', r'(path = "a\\b")'),
    ("backslash end", r'path = "dir\\" and qty = 1', r'((path = "dir\\") and (qty = 1))'),
    ("decimal", 'price <= 9.5', '(price <= 9.5)'),
]


def query_test():
    rows = "\n".join("\t\t{%s, %s, %s}," % (go_quote(n), go_quote(s), go_quote(w)) for n, s, w in QUERY_CASES)
    return d('''
    package query

    import "testing"

    func TestParse(t *testing.T) {
    	cases := []struct{ name, src, want string }{
    ''') + rows + "\n" + d('''
    	}
    	for _, tc := range cases {
    		t.Run(tc.name, func(t *testing.T) {
    			e, err := Parse(tc.src)
    			if err != nil {
    				t.Fatalf("Parse(%q): %v", tc.src, err)
    			}
    			if e.String() != tc.want {
    				t.Fatalf("Parse(%q)\\n got: %s\\nwant: %s", tc.src, e, tc.want)
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
    ''')


GO["query/parser_test.go"] = query_test()

GO["query/eval.go"] = d('''
    package query

    import "inventory/store"

    // Match reports whether item satisfies e.
    func Match(e Expr, it store.Item) bool {
    	switch x := e.(type) {
    	case And:
    		return Match(x.Left, it) && Match(x.Right, it)
    	case Or:
    		return Match(x.Left, it) || Match(x.Right, it)
    	case Compare:
    		return compare(x, it)
    	}
    	return false
    }

    func compare(c Compare, it store.Item) bool {
    	if c.Value.IsString {
    		var got string
    		switch c.Field {
    		case "sku":
    			got = it.SKU
    		case "name":
    			got = it.Name
    		default:
    			return false
    		}
    		switch c.Op {
    		case "=":
    			return got == c.Value.Str
    		case "!=":
    			return got != c.Value.Str
    		}
    		return false
    	}
    	var got float64
    	switch c.Field {
    	case "qty":
    		got = float64(it.Qty)
    	case "price":
    		got = it.Price
    	default:
    		return false
    	}
    	switch c.Op {
    	case "=":
    		return got == c.Value.Num
    	case "!=":
    		return got != c.Value.Num
    	case "<":
    		return got < c.Value.Num
    	case "<=":
    		return got <= c.Value.Num
    	case ">":
    		return got > c.Value.Num
    	case ">=":
    		return got >= c.Value.Num
    	}
    	return false
    }
''')

GO["query/eval_test.go"] = d('''
    package query

    import (
    	"testing"

    	"inventory/store"
    )

    func TestMatch(t *testing.T) {
    	it := store.Item{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25}
    	cases := map[string]bool{
    		`qty > 10`:                   true,
    		`qty < 10`:                   false,
    		`name = "bolt"`:              true,
    		`name != "bolt"`:             false,
    		`qty > 10 and price < 1`:     true,
    		`qty > 100 or sku = "B-1"`:   true,
    		`(qty > 100 or qty < 1) and price < 1`: false,
    	}
    	for src, want := range cases {
    		t.Run(src, func(t *testing.T) {
    			e, err := Parse(src)
    			if err != nil {
    				t.Fatal(err)
    			}
    			if got := Match(e, it); got != want {
    				t.Fatalf("Match(%s) = %v, want %v", src, got, want)
    			}
    		})
    	}
    }
''')

GO["store/item.go"] = d('''
    package store

    // Item is one stocked inventory line.
    type Item struct {
    	SKU   string
    	Name  string
    	Qty   int
    	Price float64
    }
''')

GO["store/store.go"] = d('''
    // Package store is an in-memory inventory store.
    package store

    import (
    	"context"
    	"errors"
    	"sort"
    	"sync"
    )

    var ErrNotFound = errors.New("item not found")

    type Store struct {
    	mu    sync.RWMutex
    	items map[string]Item
    }

    func New() *Store { return &Store{items: map[string]Item{}} }

    // Get returns the item stored under sku.
    func (s *Store) Get(ctx context.Context, sku string) (Item, error) {
    	if err := ctx.Err(); err != nil {
    		return Item{}, err
    	}
    	s.mu.RLock()
    	defer s.mu.RUnlock()
    	it, ok := s.items[sku]
    	if !ok {
    		return Item{}, ErrNotFound
    	}
    	return it, nil
    }

    func (s *Store) Put(it Item) {
    	s.mu.Lock()
    	defer s.mu.Unlock()
    	s.items[it.SKU] = it
    }

    func (s *Store) Delete(sku string) bool {
    	s.mu.Lock()
    	defer s.mu.Unlock()
    	_, ok := s.items[sku]
    	delete(s.items, sku)
    	return ok
    }

    // List returns all items ordered by SKU.
    func (s *Store) List() []Item {
    	s.mu.RLock()
    	defer s.mu.RUnlock()
    	out := make([]Item, 0, len(s.items))
    	for _, it := range s.items {
    		out = append(out, it)
    	}
    	sort.Slice(out, func(i, j int) bool { return out[i].SKU < out[j].SKU })
    	return out
    }
''')

GO["store/store_test.go"] = d('''
    package store

    import (
    	"context"
    	"errors"
    	"testing"
    )

    func TestGetPut(t *testing.T) {
    	s := New()
    	s.Put(Item{SKU: "A", Name: "anchor", Qty: 1})
    	it, err := s.Get(context.Background(), "A")
    	if err != nil || it.Name != "anchor" {
    		t.Fatalf("Get = %v, %v", it, err)
    	}
    	if _, err := s.Get(context.Background(), "missing"); !errors.Is(err, ErrNotFound) {
    		t.Fatalf("missing err = %v", err)
    	}
    }

    func TestGetCanceled(t *testing.T) {
    	ctx, cancel := context.WithCancel(context.Background())
    	cancel()
    	if _, err := New().Get(ctx, "A"); !errors.Is(err, context.Canceled) {
    		t.Fatalf("err = %v", err)
    	}
    }

    func TestListSorted(t *testing.T) {
    	s := New()
    	for _, k := range []string{"c", "a", "b"} {
    		s.Put(Item{SKU: k})
    	}
    	got := s.List()
    	if got[0].SKU != "a" || got[2].SKU != "c" {
    		t.Fatalf("List = %v", got)
    	}
    	if !s.Delete("a") || s.Delete("a") {
    		t.Fatal("Delete semantics")
    	}
    }
''')

GO["util/mathx.go"] = d('''
    // Package util holds small shared helpers.
    package util

    // Clamp limits v to [lo, hi].
    func Clamp(v, lo, hi int) int {
    	if v < lo {
    		return lo
    	}
    	if v > hi {
    		return hi
    	}
    	return v
    }

    // Abs returns |v|.
    func Abs(v int) int {
    	if v < 0 {
    		return -v
    	}
    	return v
    }
''')

GO["util/strings.go"] = d('''
    package util

    import "strings"

    // PadRight pads s with spaces to width w.
    func PadRight(s string, w int) string {
    	if len(s) >= w {
    		return s
    	}
    	return s + strings.Repeat(" ", w-len(s))
    }

    // Truncate shortens s to at most n bytes, adding "..." when cut.
    func Truncate(s string, n int) string {
    	if len(s) <= n {
    		return s
    	}
    	if n <= 3 {
    		return s[:n]
    	}
    	return s[:n-3] + "..."
    }
''')

GO["util/util_test.go"] = d('''
    package util

    import "testing"

    func TestClamp(t *testing.T) {
    	for _, c := range [][4]int{{5, 0, 10, 5}, {-1, 0, 10, 0}, {11, 0, 10, 10}} {
    		if got := Clamp(c[0], c[1], c[2]); got != c[3] {
    			t.Errorf("Clamp(%d,%d,%d) = %d", c[0], c[1], c[2], got)
    		}
    	}
    }

    func TestAbs(t *testing.T) {
    	if Abs(-3) != 3 || Abs(4) != 4 {
    		t.Fatal("Abs")
    	}
    }

    func TestStrings(t *testing.T) {
    	if PadRight("ab", 4) != "ab  " || Truncate("abcdefgh", 6) != "abc..." {
    		t.Fatal("strings helpers")
    	}
    }
''')

GO["api/errors.go"] = d('''
    package api

    import "net/http"

    // Code is a stable machine-readable API error code.
    type Code string

    const (
    	CodeNotFound   Code = "not_found"
    	CodeBadQuery   Code = "bad_query"
    	CodeConflict   Code = "conflict"
    	CodeQuota      Code = "quota_exceeded"
    	CodeInternal   Code = "internal"
    	CodeBadRequest Code = "bad_request"
    )

    // HTTPStatus maps an API error code to its HTTP status.
    func HTTPStatus(c Code) int {
    	switch c {
    	case CodeNotFound:
    		return http.StatusNotFound
    	case CodeBadQuery, CodeBadRequest:
    		return http.StatusBadRequest
    	case CodeConflict:
    		return http.StatusConflict
    	case CodeQuota:
    		return http.StatusTooManyRequests
    	}
    	return http.StatusInternalServerError
    }
''')

GO["api/handlers.go"] = d('''
    // Package api exposes the inventory store over HTTP.
    package api

    import (
    	"encoding/json"
    	"errors"
    	"net/http"
    	"strconv"

    	"inventory/query"
    	"inventory/store"
    	"inventory/util"
    )

    type Server struct {
    	Store *store.Store
    	Quota func(r *http.Request) bool
    }

    func (s *Server) Routes() *http.ServeMux {
    	mux := http.NewServeMux()
    	mux.HandleFunc("GET /items/{sku}", s.getItem)
    	mux.HandleFunc("POST /items", s.createItem)
    	mux.HandleFunc("GET /search", s.search)
    	mux.HandleFunc("POST /items/{sku}/adjust", s.adjust)
    	return mux
    }

    func writeError(w http.ResponseWriter, c Code, msg string) {
    	w.Header().Set("Content-Type", "application/json")
    	w.WriteHeader(HTTPStatus(c))
    	json.NewEncoder(w).Encode(map[string]string{"code": string(c), "error": msg})
    }

    func (s *Server) getItem(w http.ResponseWriter, r *http.Request) {
    	if s.Quota != nil && !s.Quota(r) {
    		writeError(w, CodeQuota, "quota exceeded")
    		return
    	}
    	it, err := s.Store.Get(r.Context(), r.PathValue("sku"))
    	if errors.Is(err, store.ErrNotFound) {
    		writeError(w, CodeNotFound, "no such item")
    		return
    	}
    	if err != nil {
    		writeError(w, CodeInternal, err.Error())
    		return
    	}
    	json.NewEncoder(w).Encode(it)
    }

    func (s *Server) createItem(w http.ResponseWriter, r *http.Request) {
    	var it store.Item
    	if err := json.NewDecoder(r.Body).Decode(&it); err != nil || it.SKU == "" {
    		writeError(w, CodeBadRequest, "invalid item")
    		return
    	}
    	if _, err := s.Store.Get(r.Context(), it.SKU); err == nil {
    		writeError(w, CodeConflict, "sku exists")
    		return
    	}
    	s.Store.Put(it)
    	w.WriteHeader(http.StatusCreated)
    }

    func (s *Server) search(w http.ResponseWriter, r *http.Request) {
    	e, err := query.Parse(r.URL.Query().Get("q"))
    	if err != nil {
    		writeError(w, CodeBadQuery, err.Error())
    		return
    	}
    	var out []store.Item
    	for _, it := range s.Store.List() {
    		if query.Match(e, it) {
    			out = append(out, it)
    		}
    	}
    	json.NewEncoder(w).Encode(out)
    }

    func (s *Server) adjust(w http.ResponseWriter, r *http.Request) {
    	it, err := s.Store.Get(r.Context(), r.PathValue("sku"))
    	if err != nil {
    		writeError(w, CodeNotFound, "no such item")
    		return
    	}
    	delta, err := strconv.Atoi(r.URL.Query().Get("by"))
    	if err != nil {
    		writeError(w, CodeBadRequest, "by must be an integer")
    		return
    	}
    	it.Qty = util.Clamp(it.Qty+delta, 0, 1_000_000)
    	s.Store.Put(it)
    	json.NewEncoder(w).Encode(it)
    }
''')

GO["api/api_test.go"] = d('''
    package api

    import (
    	"net/http"
    	"net/http/httptest"
    	"strings"
    	"testing"

    	"inventory/store"
    )

    func server() (*Server, http.Handler) {
    	st := store.New()
    	st.Put(store.Item{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25})
    	st.Put(store.Item{SKU: "N-1", Name: `12" nail`, Qty: 3, Price: 0.05})
    	s := &Server{Store: st}
    	return s, s.Routes()
    }

    func do(h http.Handler, method, url, body string) *httptest.ResponseRecorder {
    	rec := httptest.NewRecorder()
    	h.ServeHTTP(rec, httptest.NewRequest(method, url, strings.NewReader(body)))
    	return rec
    }

    func TestGetItem(t *testing.T) {
    	_, h := server()
    	if rec := do(h, "GET", "/items/B-1", ""); rec.Code != 200 || !strings.Contains(rec.Body.String(), "bolt") {
    		t.Fatalf("get: %d %s", rec.Code, rec.Body)
    	}
    	if rec := do(h, "GET", "/items/X", ""); rec.Code != 404 {
    		t.Fatalf("missing: %d", rec.Code)
    	}
    }

    func TestCreateConflict(t *testing.T) {
    	_, h := server()
    	if rec := do(h, "POST", "/items", `{"SKU":"B-1"}`); rec.Code != http.StatusConflict {
    		t.Fatalf("conflict: got %d want 409", rec.Code)
    	}
    	if rec := do(h, "POST", "/items", `{"SKU":"C-1","Name":"cog"}`); rec.Code != http.StatusCreated {
    		t.Fatalf("create: %d", rec.Code)
    	}
    	if rec := do(h, "POST", "/items", `{`); rec.Code != http.StatusBadRequest {
    		t.Fatalf("bad body: %d", rec.Code)
    	}
    }

    func TestSearch(t *testing.T) {
    	_, h := server()
    	rec := do(h, "GET", `/search?q=name+%3D+%2212%5C%22+nail%22`, "")
    	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "N-1") {
    		t.Fatalf("search escaped: %d %s", rec.Code, rec.Body)
    	}
    	if rec := do(h, "GET", "/search?q=qty+%3E", ""); rec.Code != 400 {
    		t.Fatalf("bad query: %d", rec.Code)
    	}
    }

    func TestAdjustClamps(t *testing.T) {
    	_, h := server()
    	rec := do(h, "POST", "/items/N-1/adjust?by=-10", "")
    	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"Qty":0`) {
    		t.Fatalf("adjust: %d %s", rec.Code, rec.Body)
    	}
    }
''')

GO["report/report.go"] = d('''
    // Package report renders inventory listings.
    package report

    import (
    	"context"
    	"fmt"
    	"strings"

    	"inventory/store"
    	"inventory/util"
    )

    // Render formats items. Supported formats: "table".
    func Render(items []store.Item, format string) (string, error) {
    	switch format {
    	case "table", "":
    		return table(items), nil
    	}
    	return "", fmt.Errorf("unknown format %q", format)
    }

    func table(items []store.Item) string {
    	var b strings.Builder
    	b.WriteString(util.PadRight("SKU", 8) + util.PadRight("NAME", 20) + util.PadRight("QTY", 6) + "PRICE\\n")
    	for _, it := range items {
    		fmt.Fprintf(&b, "%s%s%s%.2f\\n", util.PadRight(it.SKU, 8), util.PadRight(util.Truncate(it.Name, 19), 20), util.PadRight(fmt.Sprint(it.Qty), 6), it.Price)
    	}
    	return b.String()
    }

    // Lookup renders the named SKUs, skipping unknown ones.
    func Lookup(ctx context.Context, s *store.Store, skus []string) string {
    	var items []store.Item
    	for _, sku := range skus {
    		if it, err := s.Get(ctx, sku); err == nil {
    			items = append(items, it)
    		}
    	}
    	return table(items)
    }

    // LowStock lists items whose quantity is at most threshold.
    func LowStock(ctx context.Context, s *store.Store, threshold int) []string {
    	var out []string
    	for _, it := range s.List() {
    		cur, err := s.Get(ctx, it.SKU)
    		if err == nil && cur.Qty <= threshold {
    			out = append(out, fmt.Sprintf("%s (%d left, off by %d)", cur.SKU, cur.Qty, util.Abs(threshold-cur.Qty)))
    		}
    	}
    	return out
    }
''')

GO["report/report_test.go"] = d('''
    package report

    import (
    	"context"
    	"strings"
    	"testing"

    	"inventory/store"
    )

    func TestTable(t *testing.T) {
    	out, err := Render([]store.Item{{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25}}, "table")
    	if err != nil || !strings.Contains(out, "B-1") || !strings.Contains(out, "0.25") {
    		t.Fatalf("table: %q %v", out, err)
    	}
    	if _, err := Render(nil, "xml"); err == nil {
    		t.Fatal("expected unknown format error")
    	}
    }

    func TestLowStock(t *testing.T) {
    	s := store.New()
    	s.Put(store.Item{SKU: "A", Qty: 1})
    	s.Put(store.Item{SKU: "B", Qty: 50})
    	got := LowStock(context.Background(), s, 5)
    	if len(got) != 1 || !strings.HasPrefix(got[0], "A") {
    		t.Fatalf("LowStock = %v", got)
    	}
    	if !strings.Contains(Lookup(context.Background(), s, []string{"B", "zz"}), "50") {
    		t.Fatal("Lookup")
    	}
    }
''')

GO["cmd/inventory/main.go"] = d('''
    package main

    import (
    	"context"
    	"flag"
    	"fmt"
    	"log"
    	"net/http"
    	"os"

    	"inventory/api"
    	"inventory/report"
    	"inventory/store"
    )

    func main() {
    	addr := flag.String("addr", ":8080", "listen address")
    	show := flag.String("show", "", "print one SKU and exit")
    	flag.Parse()
    	st := store.New()
    	st.Put(store.Item{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25})
    	if *show != "" {
    		it, err := st.Get(context.Background(), *show)
    		if err != nil {
    			fmt.Fprintln(os.Stderr, err)
    			os.Exit(1)
    		}
    		out, _ := report.Render([]store.Item{it}, "table")
    		fmt.Print(out)
    		return
    	}
    	log.Fatal(http.ListenAndServe(*addr, (&api.Server{Store: st}).Routes()))
    }
''')

GO["README.md"] = d('''
    # inventory

    Small inventory service: a query language (`lexer`, `query`), an in-memory
    `store`, an HTTP `api`, text `report`s and shared `util` helpers.

    Run all tests with `go test ./...`.
''')


# ------------------------------------------------------------- large variant

def large_tree(base):
    """Adds ~300 generated files: a realistic-but-noisy error catalogue and
    many service modules, so locating a defect requires repository search."""
    t = dict(base)
    rng = random.Random(7)
    words = ["billing", "catalog", "shipping", "ledger", "audit", "pricing", "vendor", "warehouse",
             "returns", "tax", "promo", "loyalty", "invoice", "session", "export", "ingest"]
    codes = []
    for w in words:
        for k in ["missing", "invalid", "locked", "expired", "limit", "timeout", "denied", "stale"]:
            codes.append((w, k))
    # errcatalog: per-domain tables mapping domain codes to API codes.
    for w in words:
        rows = []
        for (dw, k) in codes:
            if dw != w:
                continue
            api = {"missing": "api.CodeNotFound", "invalid": "api.CodeBadRequest", "locked": "api.CodeConflict",
                   "expired": "api.CodeBadRequest", "limit": "api.CodeQuota", "timeout": "api.CodeInternal",
                   "denied": "api.CodeBadRequest", "stale": "api.CodeConflict"}[k]
            rows.append('\t"%s.%s": %s,' % (w, k, api))
        t["errcatalog/%s.go" % w] = d('''
            package errcatalog

            import "inventory/api"

            func init() {
            	register(map[string]api.Code{
            ''') + "\n".join(rows) + "\n" + d('''
            	})
            }
            ''')
    t["errcatalog/catalog.go"] = d('''
        // Package errcatalog maps internal domain error keys ("domain.kind")
        // to public API error codes.
        package errcatalog

        import "inventory/api"

        var table = map[string]api.Code{}

        func register(m map[string]api.Code) {
        	for k, v := range m {
        		table[k] = v
        	}
        }

        // Lookup returns the API code for key, or CodeInternal when unknown.
        func Lookup(key string) api.Code {
        	if c, ok := table[key]; ok {
        		return c
        	}
        	return api.CodeInternal
        }

        // Status returns the HTTP status for a domain error key.
        func Status(key string) int { return api.HTTPStatus(Lookup(key)) }
    ''')
    t["errcatalog/catalog_test.go"] = d('''
        package errcatalog

        import "testing"

        func TestKnownKeys(t *testing.T) {
        	for _, k := range []string{"billing.missing", "tax.invalid", "audit.locked"} {
        		if Lookup(k) == "internal" {
        			t.Errorf("%s unmapped", k)
        		}
        	}
        	if Status("nope.nothing") != 500 {
        		t.Error("unknown keys must be 500")
        	}
        }
    ''')
    # service modules: many files mentioning limits/quota so naive grep is noisy.
    for i, w in enumerate(words):
        for j in range(18):
            name = "%s_%02d" % (w, j)
            fns = []
            for f in range(6):
                kind = rng.choice(["limit", "quota", "timeout", "retry", "window", "burst"])
                fns.append(d('''
                    // %(Fn)s evaluates the %(kind)s policy for %(w)s step %(f)d.
                    // It returns the remaining %(kind)s budget and whether the
                    // caller should back off (see errcatalog "%(w)s.limit").
                    func %(Fn)s(used, max int) (int, bool) {
                    	remaining := max - used
                    	if remaining < %(thr)d {
                    		return 0, true
                    	}
                    	return remaining, false
                    }
                    ''' % {"Fn": "%s%s%d" % (kind.capitalize(), w.capitalize(), f), "kind": kind, "w": w, "f": f, "thr": rng.randint(0, 9)}))
            header = "// Package %s implements %s service policies.\npackage %s\n\n" % (w, w, w)
            t["services/%s/%s.go" % (w, name)] = header + "\n".join(fn.replace("func ", "func P%02d" % j, 1) for fn in fns)
    return t


# ----------------------------------------------------------------- C project

C = {}
C["Makefile"] = d('''
    CC ?= cc
    CFLAGS ?= -std=c11 -D_POSIX_C_SOURCE=200809L -Wall -Wextra -O1 -g
    SRC = src/ringbuf.c src/strbuf.c src/hashmap.c
    TESTS = tests/test_main.c tests/test_ringbuf.c tests/test_strbuf.c tests/test_hashmap.c

    all: libcollections.a

    libcollections.a: $(SRC:.c=.o)
    	ar rcs $@ $^

    %.o: %.c
    	$(CC) $(CFLAGS) -Iinclude -c $< -o $@

    run_tests: $(SRC) $(TESTS)
    	$(CC) $(CFLAGS) -Iinclude $(SRC) $(TESTS) -o $@

    test: run_tests
    	./run_tests

    clean:
    	rm -f src/*.o libcollections.a run_tests

    .PHONY: all test clean
''').replace("    \t", "\t")
C["include/collections.h"] = d('''
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
''')
C["src/ringbuf.c"] = d('''
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
        rb->data[(rb->head + rb->len) % rb->cap] = v;
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
''')
C["src/strbuf.c"] = d('''
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
''')
C["src/hashmap.c"] = d('''
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
''')
C["tests/test.h"] = d('''
    #ifndef TEST_H
    #define TEST_H
    #include <stdio.h>
    extern int failures, checks;
    #define CHECK(name, cond) do { checks++; if (cond) printf("ok   %s\\n", name); \\
        else { failures++; printf("FAIL %s  (%s:%d: %s)\\n", name, __FILE__, __LINE__, #cond); } } while (0)
    void test_ringbuf(void);
    void test_strbuf(void);
    void test_hashmap(void);
    #endif
''')
C["tests/test_main.c"] = d('''
    #include "test.h"
    int failures, checks;
    int main(void) {
        test_ringbuf();
        test_strbuf();
        test_hashmap();
        printf("%d checks, %d failures\\n", checks, failures);
        return failures ? 1 : 0;
    }
''')


def c_ringbuf_tests():
    body = ['#include <stdio.h>', '#include "collections.h"', '#include "test.h"', '', 'void test_ringbuf(void) {',
            '    ringbuf rb; int v; char name[64];', '    rb_init(&rb, 4);']
    body.append('    CHECK("ringbuf/empty pop", !rb_pop(&rb, &v));')
    for i in range(4):
        body.append('    CHECK("ringbuf/push %d", rb_push(&rb, %d));' % (i, i * 10))
    body.append('    CHECK("ringbuf/full rejects", !rb_push(&rb, 99));')
    # many wraparound rounds
    body.append('    for (int round = 0; round < 12; round++) {')
    body.append('        int got = -1;')
    body.append('        snprintf(name, sizeof name, "ringbuf/wrap round %02d pop", round);')
    body.append('        CHECK(name, rb_pop(&rb, &got) && got == round * 10);')
    body.append('        snprintf(name, sizeof name, "ringbuf/wrap round %02d push", round);')
    body.append('        CHECK(name, rb_push(&rb, (round + 4) * 10));')
    body.append('        snprintf(name, sizeof name, "ringbuf/wrap round %02d peek oldest", round);')
    body.append('        CHECK(name, rb_peek(&rb, 0, &got) && got == (round + 1) * 10);')
    body.append('        snprintf(name, sizeof name, "ringbuf/wrap round %02d peek newest", round);')
    body.append('        CHECK(name, rb_peek(&rb, 3, &got) && got == (round + 4) * 10);')
    body.append('    }')
    body.append('    CHECK("ringbuf/len", rb_len(&rb) == 4);')
    body.append('    CHECK("ringbuf/peek out of range", !rb_peek(&rb, 4, &v));')
    body.append('    rb_free(&rb);')
    body.append('}')
    return "\n".join(body) + "\n"


C["tests/test_ringbuf.c"] = c_ringbuf_tests()
C["tests/test_strbuf.c"] = d('''
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
''')
C["tests/test_hashmap.c"] = d('''
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
''')
C["README.md"] = "# collections\n\nSmall C collections library. Build with `make`, test with `make test`.\n"

# --------------------------------------------------------------- C++ project

CPP = {}
CPP["Makefile"] = d('''
    CXX ?= g++
    CXXFLAGS ?= -std=c++17 -Wall -Wextra -O1 -Iinclude
    SRC = src/shapes.cpp src/scene.cpp src/render.cpp src/stats.cpp
    OBJ = $(SRC:.cpp=.o)

    all: scene_demo

    %.o: %.cpp include/geometry.hpp include/scene.hpp
    	$(CXX) $(CXXFLAGS) -c $< -o $@

    scene_demo: $(OBJ) src/main.cpp
    	$(CXX) $(CXXFLAGS) $(OBJ) src/main.cpp -o $@

    run_tests: $(OBJ) tests/tests.cpp
    	$(CXX) $(CXXFLAGS) $(OBJ) tests/tests.cpp -o $@

    test: run_tests
    	./run_tests

    clean:
    	rm -f src/*.o scene_demo run_tests

    .PHONY: all test clean
''').replace("    \t", "\t")
CPP["include/geometry.hpp"] = d('''
    #pragma once
    #include <memory>
    #include <string>

    namespace geo {

    struct Point {
        double x = 0, y = 0;
    };

    // Shape is the polymorphic base of every drawable primitive. All queries
    // are const: shapes are immutable once placed in a Scene.
    class Shape {
    public:
        explicit Shape(std::string name) : name_(std::move(name)) {}
        virtual ~Shape() = default;
        virtual double area() const = 0;
        virtual double perimeter() const = 0;
        virtual Point centroid() const = 0;
        const std::string& name() const { return name_; }

    private:
        std::string name_;
    };

    class Circle : public Shape {
    public:
        Circle(std::string name, Point c, double r);
        double area() const override;
        double perimeter() const override;
        Point centroid() const override;

    private:
        Point c_;
        double r_;
    };

    class Rect : public Shape {
    public:
        Rect(std::string name, Point min, Point max);
        double area() const override;
        double perimeter() const override;
        Point centroid() const override;

    private:
        Point min_, max_;
    };

    class Triangle : public Shape {
    public:
        Triangle(std::string name, Point a, Point b, Point c);
        double area() const override;
        double perimeter() const override;
        Point centroid() const override;

    private:
        Point a_, b_, c_;
    };

    using ShapePtr = std::unique_ptr<Shape>;

    }  // namespace geo
''')
CPP["include/scene.hpp"] = d('''
    #pragma once
    #include <functional>
    #include <map>
    #include <string>
    #include <vector>
    #include "geometry.hpp"

    namespace geo {

    class Scene {
    public:
        void add(ShapePtr s);
        std::size_t size() const { return shapes_.size(); }
        double total_area() const;
        // Returns shapes sorted by descending area.
        std::vector<const Shape*> by_area() const;
        const Shape* find(const std::string& name) const;
        void for_each(const std::function<void(const Shape&)>& fn) const;

    private:
        std::vector<ShapePtr> shapes_;
        std::map<std::string, const Shape*> index_;
    };

    std::string render(const Scene& scene);

    struct Summary {
        double min_area, max_area, mean_area;
    };
    Summary summarize(const Scene& scene);

    }  // namespace geo
''')
CPP["src/shapes.cpp"] = d('''
    #include <cmath>
    #include "geometry.hpp"

    namespace geo {
    namespace {
    double dist(Point a, Point b) { return std::hypot(a.x - b.x, a.y - b.y); }
    }

    Circle::Circle(std::string name, Point c, double r) : Shape(std::move(name)), c_(c), r_(r) {}
    double Circle::area() const { return M_PI * r_ * r_; }
    double Circle::perimeter() const { return 2 * M_PI * r_; }
    Point Circle::centroid() const { return c_; }

    Rect::Rect(std::string name, Point min, Point max) : Shape(std::move(name)), min_(min), max_(max) {}
    double Rect::area() const { return (max_.x - min_.x) * (max_.y - min_.y); }
    double Rect::perimeter() const { return 2 * ((max_.x - min_.x) + (max_.y - min_.y)); }
    Point Rect::centroid() const { return {(min_.x + max_.x) / 2, (min_.y + max_.y) / 2}; }

    Triangle::Triangle(std::string name, Point a, Point b, Point c) : Shape(std::move(name)), a_(a), b_(b), c_(c) {}
    double Triangle::area() const {
        return std::fabs((b_.x - a_.x) * (c_.y - a_.y) - (c_.x - a_.x) * (b_.y - a_.y)) / 2;
    }
    double Triangle::perimeter() const { return dist(a_, b_) + dist(b_, c_) + dist(c_, a_); }
    Point Triangle::centroid() const { return {(a_.x + b_.x + c_.x) / 3, (a_.y + b_.y + c_.y) / 3}; }

    }  // namespace geo
''')
CPP["src/scene.cpp"] = d('''
    #include <algorithm>
    #include "scene.hpp"

    namespace geo {

    void Scene::add(ShapePtr s) {
        index_[s->name()] = s.get();
        shapes_.push_back(std::move(s));
    }

    double Scene::total_area() const {
        double sum = 0;
        for (const auto& s : shapes_) sum += s->area();
        return sum;
    }

    std::vector<const Shape*> Scene::by_area() const {
        std::vector<const Shape*> out;
        for (const auto& s : shapes_) out.push_back(s.get());
        std::sort(out.begin(), out.end(), [](const Shape* a, const Shape* b) { return a->area() > b->area(); });
        return out;
    }

    const Shape* Scene::find(const std::string& name) const {
        auto it = index_.find(name);
        return it == index_.end() ? nullptr : it->second;
    }

    void Scene::for_each(const std::function<void(const Shape&)>& fn) const {
        for (const auto& s : shapes_) fn(*s);
    }

    }  // namespace geo
''')
CPP["src/render.cpp"] = d('''
    #include <iomanip>
    #include <sstream>
    #include "scene.hpp"

    namespace geo {

    std::string render(const Scene& scene) {
        std::ostringstream os;
        os << std::fixed << std::setprecision(2);
        for (const Shape* s : scene.by_area()) {
            Point c = s->centroid();
            os << s->name() << " area=" << s->area() << " perim=" << s->perimeter()
               << " at (" << c.x << "," << c.y << ")\\n";
        }
        return os.str();
    }

    }  // namespace geo
''')
CPP["src/stats.cpp"] = d('''
    #include <algorithm>
    #include <limits>
    #include "scene.hpp"

    namespace geo {

    Summary summarize(const Scene& scene) {
        Summary s{std::numeric_limits<double>::max(), 0, 0};
        scene.for_each([&](const Shape& shape) {
            s.min_area = std::min(s.min_area, shape.area());
            s.max_area = std::max(s.max_area, shape.area());
        });
        s.mean_area = scene.size() ? scene.total_area() / scene.size() : 0;
        return s;
    }

    }  // namespace geo
''')
CPP["src/main.cpp"] = d('''
    #include <iostream>
    #include "scene.hpp"

    int main() {
        geo::Scene scene;
        scene.add(std::make_unique<geo::Circle>("sun", geo::Point{0, 0}, 2));
        scene.add(std::make_unique<geo::Rect>("field", geo::Point{0, 0}, geo::Point{4, 3}));
        scene.add(std::make_unique<geo::Triangle>("roof", geo::Point{0, 0}, geo::Point{4, 0}, geo::Point{2, 3}));
        std::cout << geo::render(scene);
    }
''')
CPP["tests/tests.cpp"] = d('''
    #include <cmath>
    #include <cstdio>
    #include <memory>
    #include "scene.hpp"

    static int failures = 0;
    #define CHECK(name, cond) do { if (cond) std::printf("ok   %s\\n", name); else { failures++; std::printf("FAIL %s\\n", name); } } while (0)

    int main() {
        using namespace geo;
        Scene scene;
        scene.add(std::make_unique<Circle>("c", Point{0, 0}, 1));
        scene.add(std::make_unique<Rect>("r", Point{0, 0}, Point{2, 3}));
        scene.add(std::make_unique<Triangle>("t", Point{0, 0}, Point{4, 0}, Point{0, 3}));
        CHECK("circle area", std::fabs(scene.find("c")->area() - M_PI) < 1e-9);
        CHECK("rect area", scene.find("r")->area() == 6);
        CHECK("triangle area", scene.find("t")->area() == 6);
        CHECK("triangle perimeter", scene.find("t")->perimeter() == 12);
        CHECK("total area", std::fabs(scene.total_area() - (12 + M_PI)) < 1e-9);
        auto order = scene.by_area();
        CHECK("by_area size", order.size() == 3);
        CHECK("by_area last is circle", order[2]->name() == "c");
        Summary s = summarize(scene);
        CHECK("summary max", s.max_area == 6);
        CHECK("summary min", std::fabs(s.min_area - M_PI) < 1e-9);
        CHECK("render mentions roof-less scene", render(scene).find("t area=6.00") != std::string::npos);
        std::printf("%d failures\\n", failures);
        return failures ? 1 : 0;
    }
''')
CPP["README.md"] = "# scene\n\nC++17 2D scene library. Build with `make`, test with `make test`.\n"


# ------------------------------------------------------------ task defects

def replace_once(tree, path, old, new):
    assert tree[path].count(old) == 1, (path, old)
    tree[path] = tree[path].replace(old, new)


def go_test_repair():
    t = dict(GO)
    # Escaped backslash writes the backslash but also swallows the next char.
    replace_once(t, "lexer/lexer.go", "\t\t\tcase '\"', '\\\\':\n\t\t\t\tb.WriteByte(s[i])\n",
                 "\t\t\tcase '\"':\n\t\t\t\tb.WriteByte(s[i])\n\t\t\tcase '\\\\':\n\t\t\t\tb.WriteByte(s[i])\n\t\t\t\ti++\n")
    return t


def go_compile_repair():
    t = dict(GO)
    # API change: Get gains a context parameter; implementation and its tests
    # are updated, callers are not.
    for path in ["api/handlers.go", "report/report.go", "cmd/inventory/main.go"]:
        t[path] = t[path].replace("Store.Get(r.Context(), ", "Store.Get(").replace("s.Get(ctx, ", "s.Get(").replace("st.Get(context.Background(), ", "st.Get(")
    t["cmd/inventory/main.go"] = t["cmd/inventory/main.go"].replace('\t"context"\n', "")
    return t


def go_feature():
    return dict(GO)  # feature is added by the agent; hidden test checks it


def go_refactor():
    return dict(GO)


def go_review_states():
    """Base (committed) = GO. Setup overlays an uncommitted refactor of the
    error mapping and handler helpers that introduces one regression."""
    after = {}
    after["api/errors.go"] = d('''
        package api

        import "net/http"

        // Code is a stable machine-readable API error code.
        type Code string

        const (
        	CodeNotFound   Code = "not_found"
        	CodeBadQuery   Code = "bad_query"
        	CodeConflict   Code = "conflict"
        	CodeQuota      Code = "quota_exceeded"
        	CodeInternal   Code = "internal"
        	CodeBadRequest Code = "bad_request"
        )

        // statusByCode is the single source of truth for code -> HTTP status.
        var statusByCode = map[Code]int{
        	CodeNotFound:   http.StatusNotFound,
        	CodeBadQuery:   http.StatusBadRequest,
        	CodeBadRequest: http.StatusBadRequest,
        	CodeConflict:   http.StatusBadRequest,
        	CodeQuota:      http.StatusTooManyRequests,
        	CodeInternal:   http.StatusInternalServerError,
        }

        // HTTPStatus maps an API error code to its HTTP status.
        func HTTPStatus(c Code) int {
        	if s, ok := statusByCode[c]; ok {
        		return s
        	}
        	return http.StatusInternalServerError
        }
    ''')
    h = GO["api/handlers.go"]
    h = h.replace("func writeError(w http.ResponseWriter, c Code, msg string) {", "// respondError writes a JSON error body with the mapped status.\nfunc respondError(w http.ResponseWriter, c Code, msg string) {")
    h = h.replace("\t\twriteError(", "\t\trespondError(")
    after["api/handlers.go"] = h
    r = GO["report/report.go"]
    r = r.replace("// Render formats items. Supported formats: \"table\".", "// Render formats items for display. Supported formats: \"table\" (default).")
    after["report/report.go"] = r
    u = GO["util/strings.go"].replace("// Truncate shortens s to at most n bytes, adding \"...\" when cut.", "// Truncate shortens s to at most n bytes. When s is cut, the result ends\n// with \"...\" and is exactly n bytes long.")
    after["util/strings.go"] = u
    return after


def go_search_large():
    t = large_tree(GO)
    # Defect hidden among 128 near-identical catalogue rows.
    replace_once(t, "errcatalog/warehouse.go", '"warehouse.limit": api.CodeQuota', '"warehouse.limit": api.CodeInternal')
    return t


def c_bugfix():
    t = dict(C)
    # Wraparound bug: push ignores head once the buffer has wrapped.
    replace_once(t, "src/ringbuf.c", "rb->data[(rb->head + rb->len) % rb->cap] = v;", "rb->data[rb->len % rb->cap] = v;")
    return t


def cpp_compile_repair():
    t = dict(CPP)
    # Derived classes predate the const-correct Shape interface.
    h = t["include/geometry.hpp"]
    for cls in ["Circle", "Rect", "Triangle"]:
        pass
    h = h.replace("double area() const override;", "double area() override;")
    h = h.replace("double perimeter() const override;", "double perimeter() override;")
    t["include/geometry.hpp"] = h
    s = t["src/shapes.cpp"]
    for m in ["area", "perimeter"]:
        for cls in ["Circle", "Rect", "Triangle"]:
            s = s.replace("double %s::%s() const" % (cls, m), "double %s::%s()" % (cls, m))
    t["src/shapes.cpp"] = s
    return t


FIXTURE_TREES = {
    "p7-go-test-repair": go_test_repair,
    "p7-go-compile-repair": go_compile_repair,
    "p7-go-feature": go_feature,
    "p7-go-refactor": go_refactor,
    "p7-go-search-large": go_search_large,
    "p7-c-bugfix": c_bugfix,
    "p7-cpp-compile-repair": cpp_compile_repair,
}


def write_tree(path, tree):
    if os.path.exists(path):
        shutil.rmtree(path)
    for rel, content in sorted(tree.items()):
        full = os.path.join(path, rel)
        os.makedirs(os.path.dirname(full), exist_ok=True)
        with open(full, "w", newline="\n") as f:
            f.write(content)


def generate():
    for name, fn in FIXTURE_TREES.items():
        write_tree(os.path.join(FIXTURES, name), fn())
    review = dict(GO)
    for rel, content in go_review_states().items():
        review[".states/after/" + rel] = content
    write_tree(os.path.join(FIXTURES, "p7-go-review"), review)


# -------------------------------------------------------------- oracle check

TESTS_UNCHANGED = ["git", "diff", "--quiet", "HEAD", "--", "*_test.go"]

HIDDEN_CSV_TEST = """package report

import (
	"testing"

	"inventory/store"
)

func TestHiddenCSV(t *testing.T) {
	out, err := Render([]store.Item{{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25}, {SKU: "N-2", Name: "nail, 12\\" long", Qty: 3, Price: 1}}, "csv")
	if err != nil {
		t.Fatal(err)
	}
	want := "sku,name,qty,price\\nB-1,bolt,12,0.25\\nN-2,\\"nail, 12\\"\\" long\\",3,1.00\\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
	if tbl, _ := Render(nil, "table"); tbl == "" {
		t.Fatal("table format regressed")
	}
}
"""

HIDDEN_CATALOG_TEST = """package errcatalog

import "testing"

func TestHiddenCatalog(t *testing.T) {
	for _, w := range []string{"billing", "warehouse", "tax", "vendor"} {
		if got := Status(w + ".limit"); got != 429 {
			t.Errorf("%s.limit = %d, want 429", w, got)
		}
		if got := Status(w + ".timeout"); got != 500 {
			t.Errorf("%s.timeout = %d, want 500", w, got)
		}
		if got := Status(w + ".locked"); got != 409 {
			t.Errorf("%s.locked = %d, want 409", w, got)
		}
	}
}
"""

HIDDEN_MATHX_TEST = """package mathx

import "testing"

func TestHiddenMathx(t *testing.T) {
	if Clamp(-5, 0, 9) != 0 || Clamp(50, 0, 9) != 9 || Clamp(4, 0, 9) != 4 || Abs(-7) != 7 || Abs(3) != 3 {
		t.Fatal("mathx behavior changed")
	}
}
"""


def hidden(pkg_dir, name, source, run):
    """Verifier-only test injected after the agent finishes (§19)."""
    script = "cat > %s/%s <<'ACAP_HIDDEN'\n%sACAP_HIDDEN\n%s; rc=$?; rm -f %s/%s; exit $rc" % (pkg_dir, name, source, run, pkg_dir, name)
    return {"run": {"argv": ["sh", "-c", script]}}


def run(*argv):
    return {"run": {"argv": list(argv)}}


REVIEW_FILES = sorted(go_review_states().keys())

WORKLOADS = [
    dict(name="phase7/go-test-repair", fixture="p7-go-test-repair", category="test-repair", language="go", scale="medium",
         task="The Go test suite (`go test ./...`) is failing. Find the root cause in the implementation and fix it.\nDo not modify any test files.",
         verify=[run("go", "test", "./..."), run("go", "vet", "./..."), run(*TESTS_UNCHANGED)], ref="go"),
    dict(name="phase7/go-compile-repair", fixture="p7-go-compile-repair", category="compile-repair", language="go", scale="medium",
         task="The module no longer builds: `store.Store.Get` recently gained a `context.Context` first parameter, but its callers were not updated.\nUpdate the code so that `go build ./...` and `go test ./...` succeed. Keep the new Get signature and do not modify test files.",
         verify=[run("go", "build", "./..."), run("go", "test", "./..."), run(*TESTS_UNCHANGED),
                 run("grep", "-q", "func (s \\*Store) Get(ctx context.Context, sku string)", "store/store.go")], ref="go"),
    dict(name="phase7/go-feature-csv", fixture="p7-go-feature", category="feature", language="go", scale="medium",
         task="Add a \"csv\" output format to `report.Render` in package report.\nOutput a header line `sku,name,qty,price`, then one line per item in the given order; price uses exactly two decimals.\nFields containing a comma, double quote or newline are quoted per RFC 4180 (quotes doubled). Every line ends with \\n.\nKeep the existing formats working and make sure `go test ./...` passes.",
         verify=[run("go", "test", "./..."), hidden("report", "zz_hidden_test.go", HIDDEN_CSV_TEST, "go test ./report/")], ref="go-feature"),
    dict(name="phase7/go-refactor-mathx", fixture="p7-go-refactor", category="refactor", language="go", scale="medium",
         task="Refactor: move `Clamp` and `Abs` out of package `util` into a new package `inventory/mathx` (directory `mathx/`).\nRemove them from util, update every caller, and move their tests. Behavior must not change; `go build ./...` and `go test ./...` must pass.",
         verify=[run("go", "build", "./..."), run("go", "test", "./..."),
                 run("sh", "-c", "! grep -rnE 'util\\.(Clamp|Abs)\\b' --include=*.go . && ! grep -nE '^func (Clamp|Abs)\\(' util/*.go"),
                 hidden("mathx", "zz_hidden_test.go", HIDDEN_MATHX_TEST, "go test ./mathx/")], ref="go-refactor"),
    dict(name="phase7/go-review-diff", fixture="p7-go-review", category="git-review", language="go", scale="medium",
         setup=[{"copy": {"from": ".states/after/" + f, "to": f}} for f in REVIEW_FILES],
         task="The working tree contains an uncommitted refactor (inspect it with git). Review the change: it introduces a behavioral regression.\nFix the regression while keeping the refactor's structure (the `statusByCode` table and the `respondError` helper).\nDo not commit, and do not modify test files.",
         verify=[run("go", "test", "./..."), run(*TESTS_UNCHANGED), run("grep", "-q", "statusByCode", "api/errors.go"),
                 run("grep", "-q", "respondError", "api/handlers.go")], ref="go-review"),
    dict(name="phase7/go-search-large", fixture="p7-go-search-large", category="search", language="go", scale="large",
         task="Clients report that warehouse capacity-limit errors are returned as HTTP 500 instead of 429 Too Many Requests.\nFind the cause in this repository and fix it without changing any unrelated error mappings. `go test ./...` must pass.",
         verify=[run("go", "build", "./..."), run("go", "test", "./..."),
                 hidden("errcatalog", "zz_hidden_test.go", HIDDEN_CATALOG_TEST, "go test ./errcatalog/")], ref="go-large"),
    dict(name="phase7/c-bugfix-ringbuf", fixture="p7-c-bugfix", category="bug-fix", language="c", scale="small",
         task="`make test` fails in this C library. Find and fix the bug in the library sources under src/.\nDo not modify anything under tests/.",
         verify=[run("sh", "-c", "make -s clean >/dev/null 2>&1; make -s test"), run("git", "diff", "--quiet", "HEAD", "--", "tests")], ref="c"),
    dict(name="phase7/cpp-compile-repair", fixture="p7-cpp-compile-repair", category="compile-repair", language="cpp", scale="small",
         task="The C++ project fails to build with `make`. Fix it so that `make` and `make test` succeed.\nThe `Shape` interface in include/geometry.hpp is intentionally const-correct: do not remove `const` from Shape's methods. Do not modify tests/.",
         verify=[run("sh", "-c", "make -s clean >/dev/null 2>&1; make -s && make -s test"), run("git", "diff", "--quiet", "HEAD", "--", "tests"),
                 run("grep", "-q", "virtual double area() const = 0;", "include/geometry.hpp")], ref="cpp"),
]


def yaml_str(s):
    return json.dumps(s)  # JSON strings are valid YAML scalars


def write_workloads():
    out = os.path.join(ROOT, "workloads", "phase7")
    os.makedirs(out, exist_ok=True)
    for wl in WORKLOADS:
        lines = ["# Generated by benchmarks/tools/gen_phase7.py; do not edit.", "version: 1",
                 "name: " + wl["name"], "fixture: ../../fixtures/" + wl["fixture"],
                 "category: " + wl["category"], "language: " + wl["language"], "scale: " + wl["scale"],
                 "cache_policy: warm", "git:", "  init: true", "timeout: 15m", "task: " + yaml_str(wl["task"])]
        if wl.get("setup"):
            lines.append("setup:")
            for st in wl["setup"]:
                lines.append("  - copy: {from: %s, to: %s}" % (yaml_str(st["copy"]["from"]), yaml_str(st["copy"]["to"])))
        lines.append("verify:")
        for v in wl["verify"]:
            lines.append("  - run:")
            lines.append("      argv: " + json.dumps(v["run"]["argv"]))
        with open(os.path.join(out, wl["name"].split("/")[1] + ".yaml"), "w", newline="\n") as f:
            f.write("\n".join(lines) + "\n")


def check():
    """Reference solutions must pass; task fixtures must fail (Phase 7A)."""
    ok = True
    for wl in WORKLOADS:
        fixture = os.path.join(FIXTURES, wl["fixture"])
        for label in ("fixture", "reference"):
            with tempfile.TemporaryDirectory() as tmp:
                if label == "fixture":
                    shutil.copytree(fixture, tmp, dirs_exist_ok=True)
                    shutil.rmtree(os.path.join(tmp, ".states"), ignore_errors=True)
                else:
                    write_tree(tmp, REFERENCES[wl["ref"]]())
                subprocess.run("git init -q && git add -A && git -c user.email=b@x -c user.name=b commit -qm base", shell=True, cwd=tmp, check=True)
                if label == "fixture":
                    for st in wl.get("setup", []):
                        shutil.copy(os.path.join(fixture, st["copy"]["from"]), os.path.join(tmp, st["copy"]["to"]))
                results = [subprocess.run(v["run"]["argv"], cwd=tmp, capture_output=True).returncode == 0 for v in wl["verify"]]
                passed = all(results)
                want = label == "reference"
                ok &= passed == want
                print("%s %-28s %-9s verify=%s %s" % ("ok " if passed == want else "BAD", wl["name"], label, "pass" if passed else "fail", results))
    return ok


def ref_feature():
    t = dict(GO)
    t["report/report.go"] = t["report/report.go"].replace('''	case "table", "":
		return table(items), nil
	}''', '''	case "table", "":
		return table(items), nil
	case "csv":
		var b strings.Builder
		b.WriteString("sku,name,qty,price\\n")
		for _, it := range items {
			fmt.Fprintf(&b, "%s,%s,%d,%.2f\\n", csvField(it.SKU), csvField(it.Name), it.Qty, it.Price)
		}
		return b.String(), nil
	}''') + '''
func csvField(s string) string {
	if strings.ContainsAny(s, ",\\"\\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}
'''
    return t


def ref_refactor():
    t = dict(GO)
    t["mathx/mathx.go"] = GO["util/mathx.go"].replace("// Package util holds small shared helpers.\npackage util", "// Package mathx holds integer math helpers.\npackage mathx")
    del t["util/mathx.go"]
    t["util/util_test.go"] = t["util/util_test.go"]
    for p in ["api/handlers.go", "report/report.go"]:
        t[p] = t[p].replace('"inventory/util"', '"inventory/mathx"\n\t"inventory/util"').replace("util.Clamp", "mathx.Clamp").replace("util.Abs", "mathx.Abs")
    t["api/handlers.go"] = t["api/handlers.go"].replace('\t"inventory/mathx"\n\t"inventory/util"', '\t"inventory/mathx"')
    ut = t["util/util_test.go"]
    ut = ut[:ut.index("func TestClamp")] + ut[ut.index("func TestStrings"):]
    t["util/util_test.go"] = ut
    t["mathx/mathx_test.go"] = d('''
        package mathx

        import "testing"

        func TestClamp(t *testing.T) {
        	if Clamp(-1, 0, 3) != 0 || Clamp(5, 0, 3) != 3 || Abs(-2) != 2 {
        		t.Fatal("mathx")
        	}
        }
    ''')
    return t


def ref_review():
    t = dict(GO)
    for rel, content in go_review_states().items():
        t[rel] = content.replace("CodeConflict:   http.StatusBadRequest", "CodeConflict:   http.StatusConflict")
    return t


REFERENCES = {
    "go": lambda: dict(GO),
    "go-large": lambda: large_tree(GO),
    "c": lambda: dict(C),
    "cpp": lambda: dict(CPP),
    "go-feature": ref_feature,
    "go-refactor": ref_refactor,
    "go-review": ref_review,
}

if __name__ == "__main__":
    generate()
    write_workloads()
    if "--check" in sys.argv:
        sys.exit(0 if check() else 1)
