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
