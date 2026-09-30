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
